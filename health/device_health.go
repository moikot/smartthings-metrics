package health

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DeviceHealth The status of a Device.
// swagger:model DeviceHealth
type DeviceHealth struct {
	DeviceId string `json:"deviceId,omitempty"`

	State string `json:"state,omitempty"`

	LastUpdatedDate strfmt.DateTime `json:"lastUpdatedDate,omitempty"`
}

// Validate validates this device status
func (m *DeviceHealth) Validate(_ strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeviceHealth) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceHealth) UnmarshalBinary(b []byte) error {
	var res DeviceHealth
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
