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
	"github.com/moikot/smartthings-go/models"
	"github.com/moikot/smartthings-metrics/caching"
	"github.com/patrickmn/go-cache"

	"time"
)

type ExtractorFactory struct {
	cache *cache.Cache
}

func NewExtractorFactory() *ExtractorFactory {
	return &ExtractorFactory{
		cache: cache.New(24*time.Hour, 24*time.Hour),
	}
}

func (c *ExtractorFactory) GetValueExtractors(capability *models.Capability) map[string]AttributeValueExtractor {
	extractors := caching.GetOrSet(c.cache, *capability.ID, func() interface{} {
		var ves = map[string]AttributeValueExtractor{}
		for attrName, attribute := range capability.Attributes {
			extractor := c.getValueExtractor(attribute.Schema.Properties)
			if extractor != nil {
				ves[attrName] = extractor
			}
		}
		return ves
	})
	return extractors.(map[string]AttributeValueExtractor)
}

func (c *ExtractorFactory) getValueExtractor(props *models.AttributeProperties) AttributeValueExtractor {
	if props.Value.Type == "string" && props.Value.Enum != nil {
		values := map[string]float64{}
		for i, v := range props.Value.Enum {
			values[v] = float64(i)
		}
		return NewEnumValueExtractor(values)
	} else if props.Value.Type == "number" || props.Value.Type == "integer" {
		return NewNumberValueExtractor()
	}
	return nil
}
