package gometrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

var (
	MetricStore = make(map[string]MetricParams)
)

type Metrics struct {
	Port             int            `json:"port"`
	ProcessingTotal  []MetricParams `yaml:"processing_total"`
	ErrorTotal       []MetricParams `yaml:"error_total"`
	WarningTotal     []MetricParams `yaml:"warning_total"`
	ProcessingMetric *prometheus.CounterVec
	ErrorMetric      *prometheus.CounterVec
	WarningMetric    *prometheus.CounterVec
}

type MetricParams struct {
	Type         string `yaml:"type"`
	Process      string `yaml:"process"`
	TriggerTime  string `yaml:"trigger_time"`
	TriggerCount string `yaml:"trigger_count"`
}

func (m *Metrics) loadConfig(path string) {
	data, err := os.ReadFile(filepath.Join(path))
	if err != nil {
		slog.Error(fmt.Sprintf("error reading the file: %v", err))
	}
	err = yaml.Unmarshal(data, m)
	if err != nil {
		slog.Error(fmt.Sprintf("error unmarshalling the file: %v", err))
	}
}

func (m *Metrics) Init(path string) {
	m.loadConfig(path)
	m.ProcessingMetric = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "processing_total"},
		[]string{"type", "process", "trigger_time", "trigger_count"},
	)
	m.ErrorMetric = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "error_total"},
		[]string{"type", "process", "trigger_time", "trigger_count"},
	)
	m.WarningMetric = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "warning_total"},
		[]string{"type", "process", "trigger_time", "trigger_count"},
	)

	m.initMetrics()
	prometheus.MustRegister(m.ProcessingMetric, m.ErrorMetric, m.WarningMetric)
	m.IncMetric("processing", "new_test.testing", 10)

	// go func() {
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", m.Port), nil))
	// }()
}

func GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}
	
func (m *Metrics) initMetrics() {
	for _, mp := range m.ProcessingTotal {
		m.ProcessingMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(0)
		MetricStore[fmt.Sprintf("%s.%s", mp.Type, mp.Process)] = mp
	}
	for _, mp := range m.ErrorTotal {
		m.ErrorMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(0)
	}
	for _, mp := range m.WarningTotal {
		m.WarningMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(0)
	}
}

func (m *Metrics) IncMetric(metricType, key string, value float64) {
	if mp, ok := MetricStore[key]; ok {
		switch metricType {
		case "processing":
			m.ProcessingMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(value)
		case "error":
			m.ErrorMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(value)
		case "warning":
			m.WarningMetric.WithLabelValues(mp.Type, mp.Process, mp.TriggerTime, mp.TriggerCount).Add(value)
		}
	} else {
		slog.Error(fmt.Sprintf("unknown metric key: %s", key))
	}
}
