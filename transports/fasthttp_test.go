package transports

import (
	"fmt"
	"net"
	"testing"
)

var (
	f = &FasthttpTransportFactory{}
)

func TestFastHttpServerTransport_Register(t *testing.T) {
	ls, err := net.Listen("tcp4", ":9821")
	if err != nil {
		panic(err)
	}

	s, err := f.NewServerTransport(ls)
	if err != nil {
		panic(err)
	}
	s.Register("/ha", func(in []byte) ([]byte, error) {
		return in, nil
	})

	select {}
}

func TestFastHttpClientTransport_Call(t *testing.T) {
	c, err := f.NewClientTransport("127.0.0.1:9821")
	if err != nil {
		panic(err)
	}

	resp, err := c.Call("/ha", []byte("gggggg"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))
}