package transports

import (
	"fmt"
	"github.com/seeadoog/fork"
	"github.com/valyala/fasthttp"
	"net"
	"sync"
	"unsafe"
)

type fastHttpServerTransport struct {
	ls   net.Listener
	hds  map[string]func(in []byte) ([]byte, error)
	lock sync.RWMutex
}

func NewFastHttpTransport(ls net.Listener) fork.ServerTransport {
	h := &fastHttpServerTransport{
		ls:   ls,
		hds:  map[string]func(in []byte) ([]byte, error){},
		lock: sync.RWMutex{},
	}
	h.init()
	return h
}

func (f *fastHttpServerTransport) init() {
	go func() {
		err := fasthttp.Serve(f.ls, func(ctx *fasthttp.RequestCtx) {
			f.lock.RLock()
			handler := f.hds[byte2string(ctx.Request.URI().Path())]
			f.lock.RUnlock()

			writeError := func(code int, err string) {
				ctx.Response.SetStatusCode(code)
				ctx.WriteString(err)
			}
			if handler == nil {
				writeError(404, "not found !!!")
				return
			}
			output, err := handler(ctx.Request.Body())
			if err != nil {
				writeError(500, err.Error())
				return
			}
			ctx.Write(output)
			return
		})
		if err != nil {
			panic(err)
		}
	}()
	return
}

func (f *fastHttpServerTransport) Register(fun string, handler func(in []byte) ([]byte, error)) {
	f.lock.Lock()
	f.hds[fun] = handler
	f.lock.Unlock()
}

func byte2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type fastHttpClientTransport struct {
	url string
	cli *fasthttp.Client
}

func (c *fastHttpClientTransport) Call(fun string, data []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(c.url + fun)
	req.Header.SetMethod("POST")
	req.SetBodyRaw(data)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.cli.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%v %v", resp.StatusCode(), byte2string(resp.Body()))
	}

	return resp.Body(), nil
}

type FasthttpTransportFactory struct {
}

func (f *FasthttpTransportFactory) NewServerTransport(ls net.Listener) (fork.ServerTransport, error) {
	return NewFastHttpTransport(ls), nil
}

func (f *FasthttpTransportFactory) NewClientTransport(serverAddr string) (fork.ClientTransport, error) {
	return &fastHttpClientTransport{
		url: fmt.Sprintf("http://%s", serverAddr),
		cli: &fasthttp.Client{},
	}, nil
}
