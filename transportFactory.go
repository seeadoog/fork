package fork

import (
	"fmt"
	"net"
	"net/http"
)

type TransportFactory interface {
	NewServerTransport(ls net.Listener) (ServerTransport, error)
	NewClientTransport(serverAddr string) (ClientTransport, error)
}

type defaultTransportFactory struct {
}

func (d *defaultTransportFactory) NewServerTransport(ls net.Listener) (ServerTransport, error) {
	mu := http.NewServeMux()
	go func() {
		err := http.Serve(ls, mu)
		if err != nil {
			panic(err)
		}
	}()
	return NewHttpServerTransport(mu), nil
}

func (d *defaultTransportFactory) NewClientTransport(serverAddr string) (ClientTransport, error) {
	return NewHttpClientTransport(fmt.Sprintf("http://%s", serverAddr)), nil
}

var (
	defaultTransFactory = &defaultTransportFactory{}
)
