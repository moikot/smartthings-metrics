package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type motionSensor struct {
}

func NewMotionSensor() ValueExtractor {
	return &motionSensor{}
}

func (c *motionSensor) ID() string {
	return "motionSensor"
}

func (c *motionSensor) Name() string {
	return "Motion Sensor"
}

func (c *motionSensor) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	return getEnumValue(status, motionAttribute, motionSensorValues)
}
