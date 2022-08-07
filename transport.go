package fork

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ServerTransport interface {
	Register(fun string, handler func(in []byte) ([]byte, error))
}

type ClientTransport interface {
	Call(fun string, data []byte) ([]byte, error)
}

type HttpServerTransport struct {
	sv      *http.ServeMux
	callUrl string
}

func NewHttpServerTransport(sv *http.ServeMux) ServerTransport {
	return &HttpServerTransport{
		sv: sv,
	}
}

func (h *HttpServerTransport) Register(fun string, handler func(in []byte) ([]byte, error)) {
	h.sv.HandleFunc(fun, func(writer http.ResponseWriter, request *http.Request) {
		f := func() error {
			input, err := ioutil.ReadAll(request.Body)
			if err != nil {
				return err
			}
			request.Body.Close()

			output, err := handler(input)
			if err != nil {
				return err
			}
			writer.Write(output)
			return nil
		}

		err := f()
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
		}
	})
}

type HttpClientTransport struct {
	callUrl string
}

func NewHttpClientTransport(callUrl string) ClientTransport {
	return &HttpClientTransport{
		callUrl: callUrl,
	}
}

func (h *HttpClientTransport) Call(fun string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", h.callUrl+fun, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
