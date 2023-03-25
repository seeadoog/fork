package transports

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"testing"
)

var (
	f = &FasthttpTransportFactory{}
)

func TestFastHttpServerTransport_Register(t *testing.T) {
	runtime.GOMAXPROCS(1)
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

func BenchmarkFastHttpClientTransport_Call(b *testing.B) {
	c, err := f.NewClientTransport("127.0.0.1:9821")
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	reqq := []byte("gggggg")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.Call("/ha", reqq)
			if err != nil {
				panic(err)
			}
		}

	})

}

func TestServer(t *testing.T) {
	ls, err := net.Listen("tcp", ":8821")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			io.Copy(conn, conn)
		}()
	}
}

func BenchmarkClient(b *testing.B) {
	conn, err := net.Dial("tcp", "127.0.0.1:8821")
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	bs := []byte("hello world")
	res := make([]byte, 20)
	for i := 0; i < b.N; i++ {
		conn.Write(bs)
		conn.Read(res)
	}
}

func BenchmarkParallel(b *testing.B) {
	//runtime.GOMAXPROCS(2)
	b.SetParallelism(10)
	fmt.Println(runtime.NumGoroutine())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

		}
	})
}
