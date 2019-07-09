package extractors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

const (
	unitLabel = "unit"
)

type AttributeValue struct {
	Name   string
	Labels Labels
	Value  float64
}

type ValueExtractor interface {
	ID() string
	Name() string
	Extract(status models.CapabilityStatus) (*AttributeValue, error)
}

type Labels map[string]string

func getNumberWithUnit(status models.CapabilityStatus, attr string) (*AttributeValue, error) {
	attribute, ok := status[attr]
	if !ok {
		return nil, fmt.Errorf("expected attribute '%s' does not exist", attr)
	}

	valueNum, ok := (attribute.Value).(json.Number)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%s', expected '%s'",
			reflect.TypeOf(attribute.Value), "json.Number")
	}

	value, err := valueNum.Float64()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to convert '%s' to a float number",
			valueNum.String())
	}

	return &AttributeValue{
		Name: attr,
		Labels: Labels{
			unitLabel: attribute.Unit,
		},
		Value: value,
	}, nil
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

	return &AttributeValue{
		Name:   attr,
		Labels: Labels{},
		Value:  value,
	}, nil
}
