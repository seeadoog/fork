package prometheushandler

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/seeadoog/fork"
	"net/http"
	"sync"
)

type MetricInput struct {
}
type MetricsOutput struct {
	Metrics []byte `json:"metrics"`
}

var (
	bufferPool = &sync.Pool{}
)

func RegisterClientMetricsHandler(s *fork.ServerPipe, ga prometheus.Gatherer, handler string) {
	hd := http.Header{}
	fork.RegisterHandler(s, handler, func(in MetricInput) (out MetricsOutput, err error) {
		bf, ok := bufferPool.Get().(*bytes.Buffer)
		if ok {
			bf.Reset()
		} else {
			bf = bytes.NewBuffer(make([]byte, 0, 1024*32))
		}
		enc := expfmt.NewEncoder(bf, expfmt.Negotiate(hd))
		ms, err := ga.Gather()
		if err != nil {
			return out, err
		}
		for _, m := range ms {
			err = enc.Encode(m)
			if err != nil {
				return out, err
			}
		}
		out.Metrics = bf.Bytes()
		return out, nil
	})
}

func CallMetrics(c *fork.ClientPipe, handler string) ([]byte, error) {
	res, err := fork.CallR[MetricsOutput](c, handler, MetricInput{})
	if err != nil {
		return nil, err
	}
	return res.Metrics, nil
}

func ServerMetricsHandler(mt *fork.MasterTool, handler string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mt.RangeChildren(func(cmd *fork.Cmd) bool {
			res, err := CallMetrics(cmd.ClientPipe(), handler)
			if err == nil {
				_, err := writer.Write(res)
				if err != nil {
					return false
				}
			}
			return true
		})
	}
}
