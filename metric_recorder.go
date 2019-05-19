package main

import (
	"log"

	"github.com/moikot/smartthings-metrics/extractors"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricRecorder interface {
	Record(attr *extractors.AttributeValue)
}

func NewMetricRecorder() MetricRecorder {
	return &metricRecorder{
		gaugeVecs: map[string]*prometheus.GaugeVec{},
	}
}

type metricRecorder struct {
	gaugeVecs map[string]*prometheus.GaugeVec
}

func (r *metricRecorder) Record(attr *extractors.AttributeValue) {
	gaugeVec, ok := r.gaugeVecs[attr.Name]
	if !ok {
		gaugeVec = newGaugeVec(attr)
		r.gaugeVecs[attr.Name] = gaugeVec
	}

	labels := prometheus.Labels{}

	for label, val := range attr.Labels {
		labels[label] = val
	}

	gaugeVec.With(labels).Set(attr.Value)
	log.Printf("val %f", attr.Value)
}

func newGaugeVec(attr *extractors.AttributeValue) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: attr.Name,
	}

	var labelNames []string
	for name, _ := range attr.Labels {
		labelNames = append(labelNames, name)
	}

	return prometheus.NewGaugeVec(opts, labelNames)
}
