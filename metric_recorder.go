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
		measurements: map[string]*Measurement{},
	}
}

type metricRecorder struct {
	log logrus.FieldLogger
	gaugeVecs map[string]*prometheus.GaugeVec
	measurements map[string]*Measurement
}

func printLabels(labels prometheus.Labels) string {
	var sb []string
	for k, v := range labels {
		sb = append(sb, fmt.Sprintf("%s='%s'", k, v))
	}
	return strings.Join(sb, ", ")
}

func (r *metricRecorder) Record(measurements []*Measurement) {
	observedMeasurements := map[string]*Measurement{}

	for _, measurement := range measurements {
		gaugeVec, ok := r.gaugeVecs[measurement.Name]
		if !ok {
			gaugeVec = newGaugeVec(measurement)
			r.gaugeVecs[measurement.Name] = gaugeVec
		}

		measurementKey := measurement.Name + " " + fmt.Sprint(measurement.Labels)
		observedMeasurements[measurementKey] = measurement

		gaugeVec.With(measurement.Labels).Set(measurement.Value)

		r.log.Infof("gauge: '%s' value: %f labels: %s", measurement.Name, measurement.Value, printLabels(measurement.Labels))
	}

	// Delete not observed vectors
	for key, measurement := range r.measurements {
		if _, ok := observedMeasurements[key]; !ok {
			if gaugeVec, ok := r.gaugeVecs[measurement.Name]; ok {
				gaugeVec.Delete(measurement.Labels)
				r.log.Infof("removed gauge: '%s' labels: %s", measurement.Name, printLabels(measurement.Labels))
			}
		}
	}

	r.measurements = observedMeasurements
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
