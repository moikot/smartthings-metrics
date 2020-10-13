package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricRecorder interface {
	Record(attr []*Measurement)
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

func (r *metricRecorder) Record(measurements []*Measurement) {
	observedMeasurements := map[string]struct{}{}

	for _, measurement := range measurements {
		gaugeVec, ok := r.gaugeVecs[measurement.Name]
		if !ok || gaugeVec == nil {
			gaugeVec = newGaugeVec(measurement)
			r.gaugeVecs[measurement.Name] = gaugeVec
		}
		// Add to the observed measurements
		observedMeasurements[measurement.Name] = struct{}{}

		labels := prometheus.Labels{}

		for label, val := range measurement.Labels {
			labels[label] = val
		}

		gaugeVec.With(labels).Set(measurement.Value)

		r.log.Infof("gauge: '%s' value: %f labels: %s", measurement.Name, measurement.Value, printLabels(labels))
	}

	// Clean not observed measurements
	for measurement, gaugeVec := range r.gaugeVecs {
		if _, ok := observedMeasurements[measurement]; !ok && gaugeVec != nil {
			prometheus.Unregister(gaugeVec)
			r.gaugeVecs[measurement] = nil
			r.log.Infof("gauge: '%s' is not observed", measurement)
		}
	}
}

func newGaugeVec(measurement *Measurement) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: measurement.Name,
	}

	var labelNames []string
	for name, _ := range measurement.Labels {
		labelNames = append(labelNames, name)
	}

	return promauto.NewGaugeVec(opts, labelNames)
}
