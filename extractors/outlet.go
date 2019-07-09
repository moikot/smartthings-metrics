package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type outlet struct {
}

func NewOutlet() ValueExtractor {
	return &outlet{}
}

func (c *outlet) ID() string {
	return "outlet"
}

func (c *outlet) Name() string {
	return "Power Meter"
}

func (c *outlet) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getEnumValue(status, switchAttribute, switchValues)
}
