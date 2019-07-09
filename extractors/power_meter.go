package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type powerMeter struct {
}

func NewPowerMeter() ValueExtractor {
	return &powerMeter{}
}

func (c *powerMeter) ID() string {
	return "powerMeter"
}

func (c *powerMeter) Name() string {
	return "Power Meter"
}

func (c *powerMeter) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, powerAttribute)
}
