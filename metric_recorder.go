package main

import (
	"fmt"
	"log"
	"strings"

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

func printLabels(labels prometheus.Labels) string {
	var sb []string
	for k, v := range labels {
		sb = append(sb, fmt.Sprintf("%s='%s'", k, v))
	}
	return strings.Join(sb, ", ")
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

	log.Printf("attr: '%s' value: %f labels: %s", attr.Name, attr.Value, printLabels(labels))
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
