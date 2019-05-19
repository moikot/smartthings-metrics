package extractors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

const (
	Temperature = "temperature"
)

const (
	UnitLabel = "unit"
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

func (c *temperatureMeasurement) Labels() []string {
	return []string{UnitLabel}
}

func (c *temperatureMeasurement) Extract(status models.CapabilityStatus) (*AttributeValue, error) {
	attr, ok := status[Temperature]
	if !ok {
		return nil, fmt.Errorf("expected attribute '%s' does not exist",
			Temperature)
	}

	valueNum, ok := (attr.Value).(json.Number)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%s', expected '%s'",
			reflect.TypeOf(attr.Value), "json.Number")
	}

	value, err := valueNum.Float64()
	if err != nil {
		return nil, fmt.Errorf("unable to convert '%s' to a float number",
			valueNum.String())
	}

	labels := Labels{}
	labels[UnitLabel] = attr.Unit

	return &AttributeValue{
		Name:   Temperature,
		Labels: labels,
		Value:  value,
	}, nil
}
