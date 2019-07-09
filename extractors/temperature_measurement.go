package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type temperatureMeasurement struct {
}

func NewTemperatureMeasurement() ValueExtractor {
	return &temperatureMeasurement{}
}

func (c *temperatureMeasurement) ID() string {
	return "temperatureMeasurement"
}

func (c *temperatureMeasurement) Name() string {
	return "Temperature Measurement"
}

func (c *temperatureMeasurement) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, temperatureAttribute)
}
