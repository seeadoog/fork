package fork

type ServerPipe struct {
	serverTransport ServerTransport
	marshaller      Marshaller
}

type ClientPipe struct {
	clientTransport ClientTransport
	marshaller      Marshaller
}

func RegisterHandler[In any, Out any](p *ServerPipe, fun string, handler func(in In) (out Out, err error)) {

	p.serverTransport.Register(fun, func(in []byte) ([]byte, error) {
		input := new(In)
		err := p.marshaller.Decode(in, input)
		if err != nil {
			return nil, err
		}
		out, err := handler(*input)
		if err != nil {
			return nil, err
		}
		return p.marshaller.Encode(out)
	})
}

func Call(p *ClientPipe, fun string, in interface{}, out interface{}) error {
	input, err := p.marshaller.Encode(in)
	if err != nil {
		return err
	}
	output, err := p.clientTransport.Call(fun, input)
	if err != nil {
		return err
	}

	return p.marshaller.Decode(output, out)
}

func CallR[R any](p *ClientPipe, fun string, in interface{}) (r R, err error) {
	out := new(R)
	err = Call(p, fun, in, out)
	if err != nil {
		return
	}
	return *out, nil
}
