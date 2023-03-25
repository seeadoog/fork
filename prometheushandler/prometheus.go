package prometheushandler

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/seeadoog/fork"
	"net/http"
	"runtime"
	"sync"
)

type MetricInput struct {
}

type MetricsOutput struct {
	Metrics []*io_prometheus_client.MetricFamily `json:"metrics"`
}

var (
	bufferPool = &sync.Pool{}
)

func RegisterClientMetricsHandler(s *fork.ServerPipe, ga prometheus.Gatherer, handler string) {
	fork.RegisterHandler(s, handler, func(in MetricInput) (out MetricsOutput, err error) {

		ms, err := ga.Gather()
		if err != nil {
			return out, err
		}
		out.Metrics = ms
		return out, nil
	})
}

func CallMetrics(c *fork.ClientPipe, handler string) ([]*io_prometheus_client.MetricFamily, error) {
	res, err := fork.CallR[*MetricsOutput](c, handler, MetricInput{})
	if err != nil {
		return nil, err
	}

	return res.Metrics, nil
}

func ServerMetricsHandler(mt *fork.MasterTool, handler string) http.HandlerFunc {
	fmtt := expfmt.Negotiate(http.Header{})

	return func(writer http.ResponseWriter, request *http.Request) {

		metrifs := make([][]*io_prometheus_client.MetricFamily, 0, runtime.NumCPU()/2)
		mt.RangeChildren(func(cmd *fork.Cmd) bool {
			res, err := CallMetrics(cmd.ClientPipe(), handler)
			if err != nil {
				return true
			}
			metrifs = append(metrifs, res)
			return true
		})

		enc := expfmt.NewEncoder(writer, fmtt)
		for _, family := range merge(metrifs...) {
			enc.Encode(family)
		}
	}
}

type metricWrapper struct {
	mec *io_prometheus_client.Metric
	fm  *io_prometheus_client.MetricFamily
}

func merge(mfss ...[]*io_prometheus_client.MetricFamily) []*io_prometheus_client.MetricFamily {
	counterMap := map[string]map[string]*metricWrapper{}
	histormMap := map[string]map[string]*metricWrapper{}
	//gauge := map[string]map[string]*metricWrapper{}
	//sumary := map[string]map[string]*metricWrapper{}

	mgf := func(dst map[string]map[string]*metricWrapper, mf *io_prometheus_client.MetricFamily, ff func(key string, dst *io_prometheus_client.Metric, src *io_prometheus_client.Metric)) {
		for _, metric := range mf.GetMetric() {
			fm := dst[mf.GetName()]
			if fm == nil {
				fm = map[string]*metricWrapper{}
				dst[mf.GetName()] = fm
			}
			key2 := fmt.Sprintf("%v", metric.GetLabel())
			mec := fm[key2]
			if mec == nil {
				mec = &metricWrapper{
					mec: metric,
					fm:  mf,
				}
				fm[key2] = mec
			} else {
				ff(mf.GetName()+":"+key2, mec.mec, metric)
			}

		}

	}

	gagueCount := make(map[string]int)
	qq := make([]*io_prometheus_client.MetricFamily, 0)
	for _, mfs := range mfss {
		for _, mf := range mfs {
			if mf.GetType() == io_prometheus_client.MetricType_COUNTER {
				mgf(counterMap, mf, func(key string, dst *io_prometheus_client.Metric, src *io_prometheus_client.Metric) {
					v := dst.GetCounter().GetValue() + src.Counter.GetValue()
					dst.GetCounter().Value = &v
				})
			}
			if mf.GetType() == io_prometheus_client.MetricType_HISTOGRAM {
				mgf(histormMap, mf, func(key string, dst *io_prometheus_client.Metric, src *io_prometheus_client.Metric) {
					br := dst.GetHistogram().Bucket
					bt := src.GetHistogram().Bucket

					if len(br) == len(bt) {
						for i, bucket := range br {
							v := bucket.GetExemplar().GetValue() + bt[i].GetExemplar().GetValue()
							bucket.Exemplar.Value = &v
						}
						sc := dst.Histogram.GetSampleCount() + src.GetHistogram().GetSampleCount()
						ss := dst.Histogram.GetSampleSum() + src.GetHistogram().GetSampleSum()
						dst.GetHistogram().SampleSum = &ss
						dst.GetHistogram().SampleCount = &sc
					}
				})
			}
			if mf.GetType() == io_prometheus_client.MetricType_GAUGE {
				mgf(histormMap, mf, func(key string, dst *io_prometheus_client.Metric, src *io_prometheus_client.Metric) {
					gagueCount[key]++
					dg := dst.GetGauge()
					sg := src.GetGauge()
					m := float64(gagueCount[key])
					v := (dg.GetValue()*m + sg.GetValue()) / (m + 1)
					dg.Value = &v
				})
			}
			if mf.GetType() == io_prometheus_client.MetricType_SUMMARY {
				qq = append(qq, mf)
			}
		}

	}

	format := func(m map[string]map[string]*metricWrapper) []*io_prometheus_client.MetricFamily {
		res := []*io_prometheus_client.MetricFamily{}

		for _, val := range m {
			mf := &io_prometheus_client.MetricFamily{}
			for _, metric := range val {
				mf.Metric = append(mf.Metric, metric.mec)
				mf.Type = metric.fm.Type
				mf.Help = metric.fm.Help
				mf.Name = metric.fm.Name
				//mf.Metric = append(mf.Metric, metric.mec)
			}
			res = append(res, mf)
		}
		return res
	}
	qq = append(qq, format(counterMap)...)
	qq = append(qq, format(histormMap)...)
	return qq
}
