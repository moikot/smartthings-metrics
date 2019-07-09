package main

import (
	"regexp"
	"strings"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
	"github.com/moikot/smartthings-metrics/extractors"
)

type Extractors map[string]extractors.ValueExtractor

func (exts Extractors) Add(e extractors.ValueExtractor) {
	exts[e.ID()] = e
}

type StatusProcessor interface {
	Process(device *models.Device, status *models.DeviceStatus) ([]*extractors.AttributeValue, error)
}

func NewStatusProcessor() StatusProcessor {
	exts := Extractors{}

	exts.Add(extractors.NewTemperatureMeasurement())
	exts.Add(extractors.NewBattery())
	exts.Add(extractors.NewRelativeHumidityMeasurement())
	exts.Add(extractors.NewIlluminanceMeasurement())
	exts.Add(extractors.NewMotionSensor())
	exts.Add(extractors.NewWaterSensor())
	exts.Add(extractors.NewPowerMeter())
	exts.Add(extractors.NewOutlet())

	return &statusProcessor{
		extractors: exts,
	}
}

type statusProcessor struct {
	extractors Extractors
}

func (b statusProcessor) Process(device *models.Device, status *models.DeviceStatus) ([]*extractors.AttributeValue, error) {
	var res []*extractors.AttributeValue
	for component, componentStatus := range status.Components {
		for capability, capabilityStatus := range componentStatus {
			if extractor, ok := b.extractors[capability]; ok {
				val, err := extractor.Extract(capabilityStatus)
				if err != nil {
					return nil, err
				}

				val.Name = toMetricName(capability) + "_" + toMetricName(val.Name)

				val.Labels["name"] = device.Name
				val.Labels["label"] = device.Label
				val.Labels["device_type_name"] = device.DeviceTypeName
				val.Labels["component"] = component

				res = append(res, val)
			}
		}
	}
	return res, nil
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
