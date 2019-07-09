package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type battery struct {
}

func NewBattery() ValueExtractor {
	return &battery{}
}

func (c *battery) ID() string {
	return "battery"
}

func (c *battery) Name() string {
	return "Battery"
}

func (c *battery) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, batteryAttribute)
}
