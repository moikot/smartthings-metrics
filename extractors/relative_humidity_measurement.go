package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type relativeHumidityMeasurement struct {
}

func NewRelativeHumidityMeasurement() ValueExtractor {
	return &relativeHumidityMeasurement{}
}

func (c *relativeHumidityMeasurement) ID() string {
	return "relativeHumidityMeasurement"
}

func (c *relativeHumidityMeasurement) Name() string {
	return "Relative Humidity Measurement"
}

func (c *relativeHumidityMeasurement) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, humidityAttribute)
}
