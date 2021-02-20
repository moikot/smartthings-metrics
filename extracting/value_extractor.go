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

	"github.com/moikot/smartthings-go/models"
	"github.com/moikot/smartthings-metrics/readers"
	log "github.com/sirupsen/logrus"
)

type AttributeValue struct {
	// the component name
	ComponentName string

	// the capability name
	CapabilityName string

	// the attribute name
	AttributeName string

	// the measured value
	Value float64

	// the unit of measurement
	Unit string
}

type ValueExtractor struct {
	extFactory *ExtractorFactory
}

func NewValueExtractor() *ValueExtractor {
	return &ValueExtractor{
		extFactory: NewExtractorFactory(),
	}
}

func (c *ValueExtractor) GetAttributeValues(status *models.DeviceStatus, schema *readers.Schema) []AttributeValue {
	var vals []AttributeValue

	for compName, componentStatus := range status.Components {
		for capName, capabilityStatus := range componentStatus {
			exts, err := c.getExtractors(compName, capName, schema)
			if err != nil {
				log.Errorf("unable to process capability '%s': %v", capName, err)
				continue
			}
			for attrName, extractor := range exts {
				cStatus, ok := capabilityStatus[attrName]
				if !ok {
					log.Errorf("expected attribute '%s' is not found", attrName)
					continue
				}

				if cStatus.Value == nil {
					log.Debugf("value of attribute '%s' is not defined", attrName)
					continue
				}

				val, err := extractor.Extract(cStatus)
				if err != nil {
					log.Errorf("failed to get the value of attribute '%s': %v", attrName, err)
					continue
				}

				attrVal := AttributeValue{
					Value:          val.Value,
					Unit:           val.Unit,
					ComponentName:  compName,
					CapabilityName: capName,
					AttributeName:  attrName,
				}

				vals = append(vals, attrVal)
			}
		}
	}

	return vals
}

func (c *ValueExtractor) getExtractors(compName, capName string, schema *readers.Schema) (map[string]AttributeValueExtractor, error) {
	compSchema := schema.Components[compName]
	if compSchema == nil {
		return nil, fmt.Errorf("schema of component '%s' is not defined", compName)
	}

	capSchema := compSchema.Capabilities[capName]
	if capSchema == nil {
		return nil, fmt.Errorf("schema of capability '%s' is not defined in component '%s'", capName, compName)
	}

	return c.extFactory.GetValueExtractors(capSchema), nil
}
