package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type valueWithUnit struct {
	attr string
	capability string
}

func NewValueWithUnit(capability string, attr string) ValueExtractor {
	return &valueWithUnit{
		capability: capability,
		attr: attr,
	}
}

func (c *valueWithUnit) ID() string {
	return c.capability
}

func (c *valueWithUnit) Attribute() string {
	return c.attr
}

func (c *valueWithUnit) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getNumberWithUnit(status, c.attr)
}
