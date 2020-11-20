package main

import (
	"github.com/moikot/smartthings-metrics/extractors"
	"github.com/moikot/smartthings-metrics/health"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type Extractors map[string][]extractors.ValueExtractor

func (exts Extractors) Add(e extractors.ValueExtractor) {
	exts[e.ID()] = append(exts[e.ID()], e)
}

var switchValues = map[string]float64{
	"off": 0.0,
	"on":  1.0,
}

var waterSensorValues = map[string]float64{
	"dry": 0.0,
	"wet": 1.0,
}

var motionSensorValues = map[string]float64{
	"inactive": 0.0,
	"active":   1.0,
}

var buttonValues = map[string]float64{
	"pushed": 0.0,
	"held": 1.0,
	"double": 2.0,
	"pushed_2x": 3.0,
	"pushed_3x": 4.0,
	"pushed_4x": 5.0,
	"pushed_5x": 6.0,
	"pushed_6x": 7.0,
	"down": 8.0,
	"down_2x": 9.0,
	"down_3x": 10.0,
	"down_4x": 11.0,
	"down_5x": 12.0,
	"down_6x": 13.0,
	"down_hold": 14.0,
	"up": 15.0,
	"up_2x": 16.0,
	"up_3x": 17.0,
	"up_4x": 18.0,
	"up_5x": 19.0,
	"up_6x": 20.0,
	"up_hold": 21.0,
}

var holdableButtonValues = map[string]float64{
	"held": 0.0,
	"pushed": 1.0,
}

var unitsMap = map[string]string {
	"%": "percent",
	"lux": "lux",
	"s": "seconds",
	"W": "watts",
	"C": "degrees_celsius",
	"K": "degrees_kelvin",
}

var healthStateMap = map[string]float64 {
	"OFFLINE": 0,
	"UNHEALTHY": 1.0,
	"ONLINE": 2.0,
}

type Measurement struct {
	Name   string
	Labels prometheus.Labels
	Value  float64
}

type StatusProcessor interface {
	Process(statuses []*DeviceStatus) []*Measurement
}

func NewStatusProcessor(log logrus.FieldLogger) StatusProcessor {
	exts := Extractors{}

	exts.Add(extractors.NewEnumValue("button", "button", buttonValues))
	exts.Add(extractors.NewEnumValue("holdableButton", "button", holdableButtonValues))
	exts.Add(extractors.NewValueWithUnit("battery", "battery"))
	exts.Add(extractors.NewValueWithUnit("temperatureMeasurement", "temperature"))
	exts.Add(extractors.NewValueWithUnit("relativeHumidityMeasurement", "humidity"))
	exts.Add(extractors.NewValueWithUnit("illuminanceMeasurement", "illuminance"))
	exts.Add(extractors.NewEnumValue("motionSensor", "motion", motionSensorValues))
	exts.Add(extractors.NewEnumValue("waterSensor", "water", waterSensorValues))
	exts.Add(extractors.NewValueWithUnit("powerMeter", "power"))
	exts.Add(extractors.NewEnumValue("outlet", "switch", switchValues))
	exts.Add(extractors.NewEnumValue("switch", "switch", switchValues))
	exts.Add(extractors.NewValueWithUnit("switchLevel", "level"))
	exts.Add(extractors.NewEnumValue("light", "switch", switchValues))
	exts.Add(extractors.NewValueWithUnit("healthCheck", "checkInterval"))
	exts.Add(extractors.NewValueWithUnit("colorTemperature", "colorTemperature"))
	exts.Add(extractors.NewValueWithUnit("colorControl", "hue"))
	exts.Add(extractors.NewValueWithUnit("colorControl", "saturation"))

	return &statusProcessor{
		extractors: exts,
		log: log,
	}
}

type statusProcessor struct {
	extractors Extractors
	log logrus.FieldLogger
}

func (c *statusProcessor) Process(statuses []*DeviceStatus) ([]*Measurement) {
	var res []*Measurement

	for _, status := range statuses {

		measurement := c.GetHealthMeasurement(status.Device, status.Health)
		if measurement != nil {
			res = append(res, measurement)
		}

		for component, componentStatus := range status.Status.Components {
			for capability, capabilityStatus := range componentStatus {
				if len(capabilityStatus) == 0 {
					continue
				}

				measurement := c.GetCapabilityMeasurement(status.Device, component, capability, capabilityStatus)
				if measurement != nil {
					res = append(res, measurement)
				}
			}
		}
	}

	return res
}

func (c *statusProcessor) GetHealthMeasurement(device *models.Device, health *health.DeviceHealth) *Measurement {
	value, ok := healthStateMap[health.State]
	if !ok {
		c.log.Errorf("health state has unsupported value '%s'", health.State)
		return nil
	}

	measurement := &Measurement{Labels: prometheus.Labels{}}
	measurement.Name = "smartthings_health_state"
	measurement.Value = value

	addDeviceLabels(measurement, device)

	return measurement;
}

func (c *statusProcessor) GetCapabilityMeasurement(device *models.Device, component string, capability string, capabilityStatus models.CapabilityStatus) *Measurement {

	if extrs, ok := c.extractors[capability]; ok {
		for _, extractor := range extrs {
			val, err := extractor.Extract(capabilityStatus)
			if err != nil {
				c.log.Errorf("failed to get attribute '%s' of capability '%s', error: %v", extractor.Attribute(), capability, err)
				continue
			}

			unitSuffix := ""

			if len(val.Unit()) > 0 {
				unit, ok := unitsMap[val.Unit()]
				if !ok {
					c.log.Errorf("attribute '%s' of capability '%s' has unsupported unit '%s'", extractor.Attribute(), capability, val.Unit())
					continue
				}
				unitSuffix = "_" + unit
			}

			measurement := &Measurement{Labels: prometheus.Labels{}}
			measurement.Name = "smartthings_" + toMetricName(capability) + "_" + toMetricName(extractor.Attribute()) + unitSuffix
			measurement.Value = val.Value()
			measurement.Labels["component"] = component

			addDeviceLabels(measurement, device)

			return measurement;
		}
	} else {
		c.log.Errorf("capability '%s' is not supported", capability)
	}
	return nil
}

var (
	notAlphaNum   = regexp.MustCompile("[^a-zA-Z0-9]+")
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toMetricName(deviceName string) string {
	deviceName = strings.ToLower(toSnakeCase(deviceName))
	return notAlphaNum.ReplaceAllString(deviceName, "_")
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func addDeviceLabels(measurement *Measurement, device *models.Device) {
	measurement.Labels["name"] = device.Name
	measurement.Labels["label"] = device.Label
	if device.Dth != nil {
		measurement.Labels["device_type_name"] = device.Dth.DeviceTypeName
		measurement.Labels["device_type_id"] = device.Dth.DeviceTypeID
		measurement.Labels["device_network_type"] = device.Dth.DeviceNetworkType
	} else {
		measurement.Labels["device_type_name"] = ""
		measurement.Labels["device_type_id"] = ""
		measurement.Labels["device_network_type"] = ""
	}
	measurement.Labels["device_id"] = device.DeviceID
	measurement.Labels["location_id"] = device.LocationID
	measurement.Labels["device_manufacturer_code"] = device.DeviceManufacturerCode
}
