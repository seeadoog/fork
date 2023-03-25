package fork

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ProcFunc func(child *Cmd) error

const (
	ForkFlag       = "forked-child-process"
	forkFlag       = "--" + ForkFlag
	childSocketDir = "sockets"
)

func init() {
	flag.Bool(ForkFlag, false, "mark the process is forked process")
	rand.Seed(time.Now().Unix())
	//err := os.Mkdir(childSocketDir, 777)
	//if err != nil {
	//	panic("make dir error:" + err.Error())
	//}
}

func IsChildren() bool {
	for _, arg := range os.Args {
		if arg == forkFlag {
			return true
		}
	}
	return false
}

type Forker struct {
	pf               ProcFunc
	childrenCmd      map[int64]*Cmd
	cid              int64
	lock             sync.Mutex
	n                int
	once             sync.Once
	done             chan struct{}
	masterListenAddr string
	serverPipe       *ServerPipe
	tf               TransportFactory
	marshal          Marshaller
	healthySeconds   int
	onChildFinished  func(cmd *Cmd, err error) (createNew bool)
	wg               sync.WaitGroup
}

func NewForker(n int) *Forker {
	f := &Forker{
		childrenCmd: map[int64]*Cmd{},
		cid:         0,
		lock:        sync.Mutex{},
		n:           n,
		once:        sync.Once{},
		done:        make(chan struct{}),
	}
	return f
}

// SetOnChildFinish 设置子进程结束的执行器，返回值为true 时会重新创建子进程
func (f *Forker) SetOnChildFinish(hd func(cmd *Cmd, err error) bool) {
	f.onChildFinished = hd
}

func (f *Forker) SetTransportFactory(tf TransportFactory) *Forker {
	f.tf = tf
	return f
}

func (f *Forker) SetTransportMarshaller(marshaller Marshaller) *Forker {
	f.marshal = marshaller
	return f
}

// SetHealthySeconds 设置子进程的健康启动时间，如果第一次启动时子进程的的生存时间小于改时间，则认为进程启动失败，主进程也会退出
func (f *Forker) SetHealthySeconds(s int) {
	f.healthySeconds = s
}

const (
	envMasterPipe = "__master_pipe_addr__"
	envChildId    = "__child_cmd_id__"
)

func (f *Forker) masterListen() error {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:")
	if err != nil {
		return fmt.Errorf("master resolve tcp  addr error:%w", err)
	}

	ls, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return fmt.Errorf("master listen tcp addr error:%w", err)
	}

	f.masterListenAddr = ls.Addr().String()
	st, err := f.tf.NewServerTransport(ls)
	if err != nil {
		return fmt.Errorf("master creat transport error:%w", err)
	}
	f.serverPipe = &ServerPipe{
		serverTransport: st,
		marshaller:      f.marshal,
	}
	return nil

}

func (f *Forker) createChildCmd() (*Cmd, error) {
	cmd := exec.Command(os.Args[0], append(os.Args[1:], forkFlag)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cid := atomic.AddInt64(&f.cid, 1)
	c := &Cmd{
		cmd: cmd,
		id:  cid,
		f:   f,
	}
	err := c.childListenPort()
	if err != nil {
		return nil, err
	}

	cmd.ExtraFiles = append(cmd.ExtraFiles, c.listenFile)

	c.setEnv(envMasterPipe, f.masterListenAddr)
	c.setEnv(envChildId, strconv.Itoa(int(cid)))

	if f.pf != nil {
		err = f.pf(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil

}

func (f *Forker) RangeChild(rf func(cmd *Cmd) bool) {
	if IsChildren() {
		panic("range child cannot called by child process")
	}
	f.lock.Lock()
	for _, cmd := range f.childrenCmd {
		if !rf(cmd) {
			break
		}
	}
	f.lock.Unlock()
}

// SetPreForkChild 设置子进程启动前执行的函数
func (f *Forker) SetPreForkChild(fun ProcFunc) {
	f.pf = fun
}

func (f *Forker) newChildClientPipe() (*ClientPipe, error) {
	url := os.Getenv(envMasterPipe)
	if url == "" {
		return nil, fmt.Errorf("cannot get master server addr")
	}
	ct, err := f.tf.NewClientTransport(url)
	if err != nil {
		return nil, fmt.Errorf("create master client pipe error:%w", err)
	}

	c := &ClientPipe{
		clientTransport: ct,
		marshaller:      f.marshal,
	}
	return c, nil
}

func (f *Forker) createChildHandlerTooL() (*ChildrenTool, error) {
	file := os.NewFile(3, "child listen file")
	ls, err := net.FileListener(file)
	if err != nil {
		return nil, fmt.Errorf("child create file listen error:%w", err)
	}
	//创建子进程和主进程通信的桥梁
	st, err := f.tf.NewServerTransport(ls)
	if err != nil {
		return nil, fmt.Errorf("create child server transport error:%w", err)
	}

	cp, err := f.newChildClientPipe()
	if err != nil {
		return nil, err
	}

	childId, err := strconv.Atoi(os.Getenv(envChildId))
	if err != nil {
		childId = -1
	}
	return &ChildrenTool{
		pipe: &ServerPipe{
			serverTransport: st,
			marshaller:      f.marshal,
		},
		cliPipe: cp,
		childId: int64(childId),
	}, nil
}

func (f *Forker) init() {
	if f.tf == nil {
		f.tf = defaultTransFactory
	}
	if f.marshal == nil {
		f.marshal = DefaultMarshal
	}
	if f.healthySeconds <= 0 {
		f.healthySeconds = 30
	}

}

// ForkProcess 启动主进程和子进程，如果需要从主进程传递文件到子进程，可以在doMasterPre 中设置子进程的文件，子进程获取文件需要从4开始.
func (f *Forker) ForkProcess(doMasterPre func(f *MasterTool) error, doChild func(c *ChildrenTool) error) error {
	f.init()
	if IsChildren() {
		t, err := f.createChildHandlerTooL()
		if err != nil {
			return err
		}
		return doChild(t)
	} else {
		err := f.masterListen()
		if err != nil {
			return err
		}
		err = doMasterPre(&MasterTool{
			f: f,
		})
		if err != nil {
			return err
		}
		err = f.fork(f.n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Forker) IsMaster() bool {
	return !IsChildren()
}

func (f *Forker) handleSignals() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	s := <-sc
	fmt.Println("receive sig :", s, "send signal to all children")
	f.close()
	f.lock.Lock()
	for _, cmd := range f.childrenCmd {
		cmd.cmd.Process.Signal(s)
	}
	f.lock.Unlock()
}

func (f *Forker) fork(n int) error {
	res := make(chan *proc, n)
	//startProcessTime := time.Now()
	wg := &f.wg

	for i := 0; i < n; i++ {
		cmd, err := f.createChildCmd()
		if err != nil {
			return fmt.Errorf("create children error:%w", err)
		}
		err = cmd.cmd.Start()
		if err != nil {
			return fmt.Errorf("children cmd start error:%w", err)
		}
		wg.Add(1)

		f.lock.Lock()
		f.childrenCmd[cmd.id] = cmd
		f.lock.Unlock()

		go func() {
			err := cmd.cmd.Wait()
			wg.Done()
			//子进程30s 内推出了，说明子进程启动异常，需要立即结束进程
			//if time.Since(startProcessTime) < time.Duration(f.healthySeconds)*time.Second {
			//	f.close()
			//	return
			//}

			res <- &proc{
				cmd: cmd,
				err: err,
			}
		}()
	}

	go f.handleSignals()

	go func() {
		for {
			select {
			case p := <-res:
				hd := f.onChildFinished
				if hd == nil {
					hd = func(cmd *Cmd, err error) (createNew bool) {
						return false
					}
				}
				createNew := hd(p.cmd, p.err)
				if !createNew {
					continue
				}
				// 等待1s 后再启动子进程
				time.Sleep(1 * time.Second)
				cmd, err := f.createChildCmd()
				if err != nil {
					return
				}
				err = cmd.cmd.Start()
				if err != nil {
					return
				}
				wg.Add(1)

				go func() {
					err := cmd.cmd.Wait()
					res <- &proc{
						cmd: cmd,
						err: err,
					}
					wg.Done()
				}()
				f.lock.Lock()
				f.childrenCmd[cmd.id] = cmd
				delete(f.childrenCmd, p.cmd.id)
				f.lock.Unlock()
			case <-f.done:
				return

			}
		}
	}()
	//wg.Wait()
	return nil
}

// Wait 会等待所有子进程执行完毕,只有主进程可以执行wait
func (f *Forker) Wait() {
	if IsChildren() {
		panic("Wait cannot called by child process")
	}
	f.wg.Wait()
}

func (f *Forker) close() {
	f.once.Do(func() {
		close(f.done)
	})
}

type Cmd struct {
	cmd        *exec.Cmd
	id         int64
	listenFile *os.File
	clientPip  *ClientPipe
	f          *Forker
}

func (c *Cmd) ID() int64 {
	return c.id
}

func (c *Cmd) Cmd() *exec.Cmd {
	return c.cmd
}

func (c *Cmd) ClientPipe() *ClientPipe {
	return c.clientPip
}

func (c *Cmd) Proc() *exec.Cmd {
	return c.cmd
}

func (c *Cmd) ChildrenId() int64 {
	return c.id
}

func (c *Cmd) childListenPort() error {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:")
	if err != nil {
		return fmt.Errorf("child resolve tcp  addr error:%w", err)
	}

	ls, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return fmt.Errorf("child listen tcp addr error:%w", err)
	}
	c.listenFile, err = ls.File()
	if err != nil {
		return fmt.Errorf("child tcp file error:%w", err)
	}
	ct, err := c.f.tf.NewClientTransport(ls.Addr().String())
	if err != nil {
		return fmt.Errorf("child creat client transport error:%w", err)
	}
	c.clientPip = &ClientPipe{
		clientTransport: ct,
		marshaller:      c.f.marshal,
	}
	return nil
}

//func (c *Cmd) childListenFile() error {
//	filename := childSocketDir + fmt.Sprintf("/child_%v.sock", c.id)
//	addr, err := net.ResolveUnixAddr("unix", filename)
//	if err != nil {
//		return fmt.Errorf("resolve child listen file error:%w", err)
//	}
//	ls, err := net.ListenUnix("unix", addr)
//	if err != nil {
//		return fmt.Errorf("listen child listen file error:%w", err)
//	}
//	c.listenFile, err = ls.File()
//	if err != nil {
//		return fmt.Errorf("child listen to file error:%w", err)
//	}
//	return nil
//}

func (c *Cmd) setEnv(k, v string) {
	c.cmd.Env = append(c.cmd.Env, k+"="+v)
}

type proc struct {
	cmd *Cmd
	err error
}

func randListenFile() string {
	src := make([]byte, 16)
	_, _ = rand.Read(src)
	return hex.EncodeToString(src)
}
