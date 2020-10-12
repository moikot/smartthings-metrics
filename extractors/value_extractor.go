package extractors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type AttributeValue struct {
	value  float64
	unit   string
}

func (a AttributeValue) Value() float64 {
	return a.value;
}

func (a AttributeValue) Unit() string {
	return a.unit;
}

func NewAttributeValue(value float64, unit string) *AttributeValue {
	return &AttributeValue{
		value: value,
		unit: unit,
	}
}

type ValueExtractor interface {
	ID() string
	Attribute() string
	Extract(status models.CapabilityStatus) (*AttributeValue, error)
}

func getNumberWithUnit(status models.CapabilityStatus, attr string) (*AttributeValue, error) {
	attribute, ok := status[attr]
	if !ok {
		return nil, fmt.Errorf("expected attribute '%s' does not exist", attr)
	}

	valueNum, ok := (attribute.Value).(json.Number)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%v', expected '%s'",
			reflect.TypeOf(attribute.Value), "json.Number")
	}

	value, err := valueNum.Float64()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to convert '%s' to a float number",
			valueNum.String())
	}

	return NewAttributeValue(value, attribute.Unit), nil
}

func getEnumValue(status models.CapabilityStatus, attr string, enumVals map[string]float64) (*AttributeValue, error) {
	attribute, ok := status[attr]
	if !ok {
		return nil, fmt.Errorf("expected attribute '%s' does not exist", attr)
	}

	valueStr, ok := (attribute.Value).(string)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%s', expected '%s'",
			reflect.TypeOf(attribute.Value), "string")
	}

	value, ok := enumVals[valueStr]
	if !ok {
		return nil, fmt.Errorf("unexpected value '%s' of attribute '%s'",
			valueStr, attr)
	}

	return NewAttributeValue(value, ""), nil
}
