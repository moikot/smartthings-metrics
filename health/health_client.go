package health

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new health API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	if formats == nil {
		formats = strfmt.Default
	}

	return &Client{transport: transport, formats: formats}
}

/*
Client for health API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

func (c *Client) GetDeviceHealth(params *GetDeviceHealthParams, authInfo runtime.ClientAuthInfoWriter) (*GetDeviceHealthOK, error) {
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