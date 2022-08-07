package fork

import (
	"fmt"
	"net/http"
	"testing"
)

type input struct {
	A int
	B int
}

type Output struct {
	Sum int
}

func TestNewHttpServerTransport(t *testing.T) {
	go func() {
		err := http.ListenAndServe("unix://./abc.sock", nil)
		panic(err)
	}()
	s := NewHttpServerTransport(http.DefaultServeMux)
	p := &ServerPipe{
		serverTransport: s,
		marshaller:      DefaultMarshal,
	}
	RegisterHandler(p, "/add", func(in input) (out Output, err error) {
		out = Output{Sum: in.A + in.B}
		return out, nil
	})
	select {}
}

func TestNewHttpClientTransport(t *testing.T) {
	p := &ClientPipe{
		clientTransport: NewHttpClientTransport("http://127.0.0.1:8080"),
		marshaller:      DefaultMarshal,
	}
	out := &Output{}
	err := Call(p, "/add", input{
		1, 2,
	}, out)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

func TestName(t *testing.T) {

}
