package prometheushandler

type GrpcConn struct {
	currentStreams int32
	caller         func() error
}

func (g *GrpcConn) Call(f func() error) {

}
