package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type illuminanceMeasurement struct {
}

func NewIlluminanceMeasurement() ValueExtractor {
	return &illuminanceMeasurement{}
}

func (c *illuminanceMeasurement) ID() string {
	return "illuminanceMeasurement"
}

func (c *illuminanceMeasurement) Name() string {
	return "Illuminance Measurement"
}

func (c *illuminanceMeasurement) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, illuminanceAttribute)
}
