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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new health API client.
func NewDeviceHealthAPI(transport runtime.ClientTransport, formats strfmt.Registry) *DeviceHealthAPI {
	if formats == nil {
		formats = strfmt.Default
	}

	return &DeviceHealthAPI{transport: transport, formats: formats}
}

/*
Client for health API
*/
type DeviceHealthAPI struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

func (c *DeviceHealthAPI) GetDeviceHealth(params *GetDeviceHealthParams, authInfo runtime.ClientAuthInfoWriter) (*GetDeviceHealthOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetDeviceHealthParams()
	}

	result, err := c.transport.Submit(&runtime.ClientOperation{
		ID:                 "getDeviceStatus",
		Method:             "GET",
		PathPattern:        "/devices/{deviceId}/health",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetDeviceHealthReader{formats: c.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})

	if err != nil {
		return nil, err
	}
	return result.(*GetDeviceHealthOK), nil
}
