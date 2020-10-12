package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricRecorder interface {
	Record(attr *Measurement)
}

func NewMetricRecorder(log logrus.FieldLogger) MetricRecorder {
	return &metricRecorder{
		log: log,
		gaugeVecs: map[string]*prometheus.GaugeVec{},
	}
}

type metricRecorder struct {
	log logrus.FieldLogger
	gaugeVecs map[string]*prometheus.GaugeVec
}

func printLabels(labels prometheus.Labels) string {
	var sb []string
	for k, v := range labels {
		sb = append(sb, fmt.Sprintf("%s='%s'", k, v))
	}
	return strings.Join(sb, ", ")
}

func (r *metricRecorder) Record(attr *Measurement) {
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

	r.log.Infof("gauge: '%s' value: %f labels: %s", attr.Name, attr.Value, printLabels(labels))
}

func newGaugeVec(attr *Measurement) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: attr.Name,
	}

	var labelNames []string
	for name, _ := range attr.Labels {
		labelNames = append(labelNames, name)
	}

	return promauto.NewGaugeVec(opts, labelNames)
}
