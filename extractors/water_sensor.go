package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type waterSensor struct {
}

func NewWaterSensor() ValueExtractor {
	return &waterSensor{}
}

func (c *waterSensor) ID() string {
	return "waterSensor"
}

func (c *waterSensor) Name() string {
	return "Water Sensor"
}

func (c *waterSensor) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getEnumValue(status, waterAttribute, waterSensorValues)
}
