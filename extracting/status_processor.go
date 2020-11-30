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

package extracting

import (
	"encoding/json"
	"fmt"
	"github.com/moikot/smartthings-go/models"
	"github.com/moikot/smartthings-metrics/readers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const unitsFileName = "units.json"

var healthStateMap = map[string]float64{
	"OFFLINE":   0,
	"UNHEALTHY": 1.0,
	"ONLINE":    2.0,
}

type Measurement struct {
	Name   string
	Labels prometheus.Labels
	Value  float64
}

type StatusProcessor interface {
	GetMeasurements(statuses []*readers.DeviceStatus) []*Measurement
}

func NewStatusProcessor(log logrus.FieldLogger) StatusProcessor {
	return &statusProcessor{
		valueExtractor: NewValueExtractor(log),
		log:            log,
		unitMap:        readUnitMap(unitsFileName),
	}
}

func readUnitMap(file string) map[string]string {
	units, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var unitMap map[string]string
	if err := json.Unmarshal(units, &unitMap); err != nil {
		panic(fmt.Sprintf("failed to unmarshal '%s': %v", file, err))
	}
	return unitMap
}

type statusProcessor struct {
	log            logrus.FieldLogger
	valueExtractor *ValueExtractor
	unitMap        map[string]string
}

func (p *statusProcessor) GetMeasurements(statuses []*readers.DeviceStatus) []*Measurement {
	var res []*Measurement

	for _, status := range statuses {

		measurement, err := p.getHealthMeasurement(status.Health)
		if err != nil {
			p.log.Errorf("failed to get device '%s' health; %v", *status.Device.DeviceID, err)
		} else {
			setDeviceLabels(measurement, status.Device)
			res = append(res, measurement)
		}

		attributeValues := p.valueExtractor.GetAttributeValues(status.Status, status.Schema)
		for _, attributeValue := range attributeValues {
			measurement := p.getAttributeValueMeasurement(attributeValue)
			setDeviceLabels(measurement, status.Device)
			res = append(res, measurement)
		}
	}

	return res
}

func (p *statusProcessor) getHealthMeasurement(health *readers.DeviceHealth) (*Measurement, error) {
	value, ok := healthStateMap[health.State]
	if !ok {
		return nil, fmt.Errorf("health state has unsupported value '%s'", health.State)
	}

	measurement := &Measurement{Labels: prometheus.Labels{}}
	measurement.Name = "smartthings_health_state"
	measurement.Value = value

	return measurement, nil
}

func (p *statusProcessor) getAttributeValueMeasurement(attrValue AttributeValue) *Measurement {
	measurement := &Measurement{Labels: prometheus.Labels{}}
	measurement.Name = getMetricName(
		attrValue.ComponentName,
		attrValue.CapabilityName,
		attrValue.AttributeName,
		p.getUnitSuffix(attrValue.Unit),
	)
	measurement.Value = attrValue.Value

	return measurement
}

func (p *statusProcessor) getUnitSuffix(unit string) string {

	if len(unit) > 0 {
		unit, ok := p.unitMap[unit]
		if ok {
			return unit
		} else {
			p.log.Warnf("unit '%s' is not supported", unit)
			return ""
		}
	}

	return ""
}

func setDeviceLabels(measurement *Measurement, device *models.Device) {
	measurement.Labels["name"] = device.Name
	measurement.Labels["label"] = device.Label
	if device.Dth != nil {
		measurement.Labels["device_type_name"] = *device.Dth.DeviceTypeName
		measurement.Labels["device_type_id"] = *device.Dth.DeviceTypeID
		measurement.Labels["device_network_type"] = device.Dth.DeviceNetworkType
	} else {
		measurement.Labels["device_type_name"] = ""
		measurement.Labels["device_type_id"] = ""
		measurement.Labels["device_network_type"] = ""
	}
	measurement.Labels["device_id"] = *device.DeviceID
	measurement.Labels["location_id"] = device.LocationID
	measurement.Labels["device_manufacturer_code"] = device.DeviceManufacturerCode
}
