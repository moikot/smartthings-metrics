/*
Copyright (c) 2020 Sergey Anisimov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package recording

import (
	"fmt"
	"strings"

	"github.com/moikot/smartthings-metrics/extracting"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

type MetricRecorder interface {
	Record(attr []*extracting.Measurement)
}

func NewMetricRecorder() MetricRecorder {
	return &metricRecorder{
		gaugeVecs:    map[string]*prometheus.GaugeVec{},
		measurements: map[string]*extracting.Measurement{},
	}
}

type metricRecorder struct {
	gaugeVecs    map[string]*prometheus.GaugeVec
	measurements map[string]*extracting.Measurement
}

func printLabels(labels prometheus.Labels) string {
	var sb []string
	for k, v := range labels {
		sb = append(sb, fmt.Sprintf("%s='%s'", k, v))
	}
	return strings.Join(sb, ", ")
}

func (r *metricRecorder) Record(measurements []*extracting.Measurement) {
	observedMeasurements := map[string]*extracting.Measurement{}

	for _, measurement := range measurements {
		gaugeVec, ok := r.gaugeVecs[measurement.Name]
		if !ok {
			gaugeVec = newGaugeVec(measurement)
			r.gaugeVecs[measurement.Name] = gaugeVec
		}

		measurementKey := measurement.Name + " " + fmt.Sprint(measurement.Labels)
		observedMeasurements[measurementKey] = measurement

		gaugeVec.With(measurement.Labels).Set(measurement.Value)

		log.Debugf("gauge: '%s' value: %f labels: %s", measurement.Name, measurement.Value, printLabels(measurement.Labels))
	}

	// Delete not observed vectors
	for key, measurement := range r.measurements {
		if _, ok := observedMeasurements[key]; !ok {
			if gaugeVec, ok := r.gaugeVecs[measurement.Name]; ok {
				gaugeVec.Delete(measurement.Labels)
				log.Debugf("removed gauge: '%s' labels: %s", measurement.Name, printLabels(measurement.Labels))
			}
		}
	}

	r.measurements = observedMeasurements
}

func newGaugeVec(measurement *extracting.Measurement) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: measurement.Name,
	}

	var labelNames []string
	for name := range measurement.Labels {
		labelNames = append(labelNames, name)
	}

	return promauto.NewGaugeVec(opts, labelNames)
}
