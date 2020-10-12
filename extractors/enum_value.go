package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type enumValue struct {
	attr string
	capability string
	enumVals map[string]float64
}

func NewEnumValue(capability string, attr string, enumVals map[string]float64) ValueExtractor {
	return &enumValue{
		capability: capability,
		attr: attr,
		enumVals: enumVals,
	}
}

func (c *enumValue) ID() string {
	return c.capability
}

func (c *enumValue) Attribute() string {
	return c.attr
}

func (c *enumValue) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getEnumValue(status, c.attr, c.enumVals)
}
