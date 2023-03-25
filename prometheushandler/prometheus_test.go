package prometheushandler

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	grpc "google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestGather(t *testing.T) {

	promhttp.Handler()
	r := prometheus.NewRegistry()
	pro := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "default",
		Subsystem:   "",
		Name:        "sms",
		Help:        "",
		ConstLabels: nil,
	}, []string{"namespace", "method"})
	pro.WithLabelValues("ddd", "GET").Inc()
	pro.WithLabelValues("abc", "POST").Inc()
	r.Register(pro)

	r2 := prometheus.NewRegistry()
	pro2 := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "default",
		Subsystem:   "",
		Name:        "sms",
		Help:        "",
		ConstLabels: nil,
	}, []string{"namespace", "method"})
	pro2.WithLabelValues("ddd", "GET").Inc()
	pro2.WithLabelValues("abc", "POST").Inc()

	ssv := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   "default",
		Subsystem:   "",
		Name:        "gas",
		Help:        "",
		ConstLabels: nil,
		Objectives: map[float64]float64{
			0.5:    0.05,
			0.9:    0.01,
			0.9999: 0.00001,
			0.999:  0.00001,
			0.99:   0.00001,
			0.95:   0.00001,
		},
		MaxAge:     0,
		AgeBuckets: 0,
		BufCap:     0,
	}, []string{"namespace", "method"})
	hsv := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "default",
		Subsystem:   "",
		Name:        "sss",
		Help:        "",
		ConstLabels: nil,
		Buckets:     nil,
	}, []string{"namespace", "method"})

	r.MustRegister(ssv, hsv)

	for i := 0; i < 10000; i++ {
		ssv.WithLabelValues("abc", "POST").Observe(float64(i))

	}

	hsv.WithLabelValues("abc", "POST").Observe(9)
	hsv.WithLabelValues("abc", "POST").Observe(22)

	r2.Register(pro2)
	mfs, err := r.Gather()
	if err != nil {
		panic(err)
	}
	mfs2, err := r2.Gather()
	if err != nil {
		panic(err)
	}
	//for _, mf := range merge(append(mfs2, mfs...)) {
	//	fmt.Println(mf.GetType(), mf.GetName(), mf.GetMetric(), len(mf.GetMetric()))
	//}
	mfs = merge(mfs, mfs2)
	enc := expfmt.NewEncoder(os.Stdout, expfmt.Negotiate(http.Header{}))
	for _, mf := range mfs {
		enc.Encode(mf)
	}

}

type grpcServer struct {
}

func (g *grpcServer) CreateUser(context context.Context, request *Request) (*Response, error) {
	return &Response{
		Code: 0,
		Sid:  "haha",
	}, nil
}

func TestGrpcs(t *testing.T) {
	s := grpc.NewServer(grpc.MaxConcurrentStreams(1000), grpc.WriteBufferSize(1024*1024*32))
	RegisterRpcServerServer(s, &grpcServer{})
	ls, err := net.Listen("tcp", ":8021")
	if err != nil {
		panic(err)
	}

	s.Serve(ls)
}

func TestGrpcc(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:8021", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli := NewRpcServerClient(conn)
	resp, err := cli.CreateUser(context.TODO(), &Request{})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
func BenchmarkGrpc(b *testing.B) {
	fmt.Println("test")
	conn, err := grpc.Dial("127.0.0.1:8021", grpc.WithInsecure(), grpc.WithReadBufferSize(32*1024), grpc.WithWriteBufferSize(1024*1024*4), grpc.WithDialer(func(s string, duration time.Duration) (net.Conn, error) {
		fmt.Println("dial")
		return net.DialTimeout("tcp", s, duration)
	}))
	if err != nil {
		panic(err)
	}

	cli := NewRpcServerClient(conn)
	name := string(make([]byte, 1024*30))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cli.CreateUser(context.TODO(), &Request{
				Name: name,
			})
		}
	})

}
func TestFast(t *testing.T) {
	fmt.Println("test")
	conn, err := grpc.Dial("127.0.0.1:8021", grpc.WithInsecure(), grpc.WithDialer(func(s string, duration time.Duration) (net.Conn, error) {
		fmt.Println("dial")
		return net.DialTimeout("tcp", s, duration)
	}))
	if err != nil {
		panic(err)
	}

	cli := NewRpcServerClient(conn)
	for i := 0; i < 300; i++ {
		go func() {
			for {
				_, err := cli.CreateUser(context.TODO(), &Request{})
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	select {}
}
