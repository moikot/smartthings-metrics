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

package readers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/moikot/smartthings-go/client/capabilities"
	"github.com/moikot/smartthings-go/models"
	"github.com/moikot/smartthings-metrics/caching"
	"github.com/patrickmn/go-cache"
	"time"
)

func NewCapabilityReader(rtm runtime.ClientTransport, formats strfmt.Registry) *CapabilityReader {
	return &CapabilityReader{
		ClientService: capabilities.New(rtm, formats),
		cache:         cache.New(24*time.Hour, 24*time.Hour),
	}
}

type CapabilityReader struct {
	capabilities.ClientService
	cache *cache.Cache
}

func (r *CapabilityReader) ReadCapability(capability string, version int64, auth runtime.ClientAuthInfoWriter) (*models.Capability, error) {
	key := fmt.Sprintf("%s/%v", capability, version)
	c, err := caching.GetOrSetE(r.cache, key, func() (interface{}, error) {
		params := capabilities.NewGetCapabilityParams()
		params.CapabilityID = capability
		params.CapabilityVersion = version

		c, err := r.GetCapability(params, auth)
		if err != nil {
			return nil, err
		}

		return c.Payload, nil
	})
	if err != nil {
		return nil, err
	}
	return c.(*models.Capability), nil
}
