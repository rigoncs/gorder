package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PrometheusMetricsClient struct {
	registery *prometheus.Registry
}

var dynamicCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "dynamic_counter",
		Help: "Count custom keys",
	},
	[]string{"key"},
)

type PrometheusMetricsClientConfig struct {
	Host        string
	ServiceName string
}

func NewPrometheusMetricsClient(config *PrometheusMetricsClientConfig) *PrometheusMetricsClient {
	client := &PrometheusMetricsClient{}
	client.initPrometheus(config)
	return client
}

func (p *PrometheusMetricsClient) initPrometheus(conf *PrometheusMetricsClientConfig) {
	p.registery = prometheus.NewRegistry()
	p.registery.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// custom collectors:
	p.registery.Register(dynamicCounter)

	// metadata wrap
	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": conf.ServiceName}, p.registery)

	// export
	http.Handle("/metrics", promhttp.HandlerFor(p.registery, promhttp.HandlerOpts{}))
	go func() {
		logrus.Fatalf("failed to start prometheus metrics endpoint, err=%v", http.ListenAndServe(conf.Host, nil))
	}()
}

func (p *PrometheusMetricsClient) Inc(key string, value int) {
	dynamicCounter.WithLabelValues(key).Add(float64(value))
}
