package health


import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/strfmt"
)

// NewGetDeviceHealthParams creates a new GetDeviceHealthParams object
// with the default values initialized.
func NewGetDeviceHealthParams() *GetDeviceHealthParams {
	var ()
	return &GetDeviceHealthParams{

		timeout: cr.DefaultTimeout,
	}
}

/*GetDeviceHealthParams contains all the parameters to send to the API endpoint
for the get device status operation typically these are written to a http.Request
*/
type GetDeviceHealthParams struct {

	/*Authorization
	  OAuth token

	*/
	Authorization string
	/*DeviceID
	  the device ID

	*/
	DeviceID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDeviceID adds the deviceID to the get device status params
func (o *GetDeviceHealthParams) WithDeviceID(deviceID string) *GetDeviceHealthParams {
	o.SetDeviceID(deviceID)
	return o
}

// SetDeviceID adds the deviceId to the get device status params
func (o *GetDeviceHealthParams) SetDeviceID(deviceID string) {
	o.DeviceID = deviceID
}

// WriteToRequest writes these params to a swagger request
func (o *GetDeviceHealthParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param Authorization
	if err := r.SetHeaderParam("Authorization", o.Authorization); err != nil {
		return err
	}

	// path param deviceId
	if err := r.SetPathParam("deviceId", o.DeviceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
