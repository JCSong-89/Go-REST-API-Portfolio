package prometheus

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"strconv"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObeserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter //request total
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec //request time
	Memory    *prometheus.GaugeVec
}

func (p *PrometheusMetrics) IncHits(status int, method, path string) {
	p.HitsTotal.Inc()
	p.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (p *PrometheusMetrics) ObeserveResponseTime(status int, method, path string, observeTime float64) {
	p.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

func (p *PrometheusMetrics) IncMemory(status int, method, path string) {
	p.Memory.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func CreateMetrics(address string, name string) (Metrics, error) {
	var pm PrometheusMetrics

	pm.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + " hits total",
	})
	if err := prometheus.Register(pm.HitsTotal); err != nil {
		return nil, err
	}

	pm.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + " hits",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(pm.Hits); err != nil {
		return nil, err
	}

	pm.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name + " times",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(pm.Times); err != nil {
		return nil, err
	}

	pm.Memory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name + " memory",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(pm.Memory); err != nil {
		return nil, err
	}

	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler())) // http handler wrap echo handler
		log.Printf("메트릭스 서버 RUN On Port: %s", address)
		if err := router.Start(address); err != nil {
			log.Fatalf("s메트릭스 서버 에러: %v", err)
		}
	}()

	return &pm, nil
}