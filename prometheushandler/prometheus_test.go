package prometheushandler

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"os"
	"testing"
)

func TestGather(t *testing.T) {

	promhttp.Handler()
	r := prometheus.NewRegistry()
	pro := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "default",
		Subsystem:   "",
		Name:        "",
		Help:        "",
		ConstLabels: nil,
	}, []string{"namespace"})
	pro.WithLabelValues("ddd").Inc()
	r.Register(pro)

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		panic(err)
	}
	enc := expfmt.NewEncoder(os.Stdout, expfmt.Negotiate(http.Header{}))
	for _, mf := range mfs {
		enc.Encode(mf)
	}
	fmt.Println(mfs)
}
