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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/moikot/smartthings-go/models"
	"github.com/pkg/errors"
)

type numberValueExtractor struct {
}

func NewNumberValueExtractor() AttributeValueExtractor {
	return &numberValueExtractor{}
}

func (c *numberValueExtractor) Extract(state models.AttributeState) (*ExtractedValue, error) {
	valueNum, ok := (state.Value).(json.Number)
	if !ok {
		return nil, fmt.Errorf("value has unexpected type '%v', expected '%s'",
			reflect.TypeOf(state.Value), "json.Number")
	}

	value, err := valueNum.Float64()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to convert '%s' to a float number", valueNum.String())
	}

	return &ExtractedValue{
		Value: value,
		Unit:  state.Unit,
	}, nil
}
