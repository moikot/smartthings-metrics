/*
Copyright (c) 2020 Sergey Anisimov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package extracting

import (
	"fmt"
	"reflect"

	"github.com/moikot/smartthings-go/models"
)

type enumValueExtractor struct {
	enumVals map[string]float64
}

func NewEnumValueExtractor(enumVals map[string]float64) AttributeValueExtractor {
	return &enumValueExtractor{
		enumVals: enumVals,
	}
}

func (c *enumValueExtractor) Extract(state models.AttributeState) (*ExtractedValue, error) {
	valueStr, ok := (state.Value).(string)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%s', expected '%s'",
			reflect.TypeOf(state.Value), "string")
	}

	value, ok := c.enumVals[valueStr]
	if !ok {
		return nil, fmt.Errorf("unexpected value '%s'", valueStr)
	}

	return &ExtractedValue{
		Value: value,
		Unit:  "",
	}, nil
}
