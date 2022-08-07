package fork

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestFork(t *testing.T) {
	f := NewForker(1)
	err := f.ForkProcess(func(f *MasterTool) error {
		addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8080")
		if err != nil {
			return fmt.Errorf("resolve addr error:%w", err)
		}
		ls, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			return fmt.Errorf("listen error:%w", err)
		}
		file, err := ls.File()
		if err != nil {
			return fmt.Errorf("convert to file error:%w", err)
		}
		f.SetForkPreRun(func(child *exec.Cmd) error {
			child.ExtraFiles = append(child.ExtraFiles, file)
			return nil
		})

		go func() {
			for {
				time.Sleep(5 * time.Second)
				f.RangeChildren(func(cmd *Cmd) bool {
					res, err := CallR[string](cmd.ClientPipe(), "/metrics", "haha")
					if err != nil {
						fmt.Println("error:", err)
						return false
					}
					fmt.Println("metrics:", res)
					return true
				})
			}
		}()
		return nil
	}, func(ct *ChildrenTool) error {
		fmt.Println("fork child")
		lsfile := os.NewFile(4, "LISTEN")
		ls, err := net.FileListener(lsfile)
		RegisterHandler(ct.ServePipe(), "/metrics", func(in string) (out string, err error) {
			return "hello", nil
		})
		if err != nil {
			return err
		}
		for {
			conn, err := ls.Accept()
			if err != nil {
				return err
			}
			go func() {
				io.Copy(conn, conn)
			}()
		}
	})
	if err != nil {
		panic(err)
	}
}

func TestRandfile(t *testing.T) {

	fmt.Println(randListenFile())
	fmt.Println(randListenFile())
}
