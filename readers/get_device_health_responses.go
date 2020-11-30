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
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/moikot/smartthings-go/models"
)

// GetDeviceHealthReader is a Reader for the GetDeviceHealth structure.
type GetDeviceHealthReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceHealthReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetDeviceHealthOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewGetDeviceHealthBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewGetDeviceHealthUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 403:
		result := NewGetDeviceHealthForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 429:
		result := NewGetDeviceHealthTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewGetDeviceHealthDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDeviceHealthOK creates a GetDeviceHealthOK with default headers values
func NewGetDeviceHealthOK() *GetDeviceHealthOK {
	return &GetDeviceHealthOK{}
}

/*GetDeviceHealthOK handles this case with default header values.

successful return of current status of device attributes
*/
type GetDeviceHealthOK struct {
	Payload *DeviceHealth
}

func (o *GetDeviceHealthOK) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealthOK  %+v", 200, o.Payload)
}

func (o *GetDeviceHealthOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(DeviceHealth)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceHealthBadRequest creates a GetDeviceHealthBadRequest with default headers values
func NewGetDeviceHealthBadRequest() *GetDeviceHealthBadRequest {
	return &GetDeviceHealthBadRequest{}
}

/*GetDeviceHealthBadRequest handles this case with default header values.

Bad request
*/
type GetDeviceHealthBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *GetDeviceHealthBadRequest) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealthBadRequest  %+v", 400, o.Payload)
}

func (o *GetDeviceHealthBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceHealthUnauthorized creates a GetDeviceHealthUnauthorized with default headers values
func NewGetDeviceHealthUnauthorized() *GetDeviceHealthUnauthorized {
	return &GetDeviceHealthUnauthorized{}
}

/*GetDeviceHealthUnauthorized handles this case with default header values.

Not authenticated
*/
type GetDeviceHealthUnauthorized struct {
}

func (o *GetDeviceHealthUnauthorized) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealthUnauthorized ", 401)
}

func (o *GetDeviceHealthUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetDeviceHealthForbidden creates a GetDeviceHealthForbidden with default headers values
func NewGetDeviceHealthForbidden() *GetDeviceHealthForbidden {
	return &GetDeviceHealthForbidden{}
}

/*GetDeviceHealthForbidden handles this case with default header values.

Not authorized
*/
type GetDeviceHealthForbidden struct {
}

func (o *GetDeviceHealthForbidden) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealthForbidden ", 403)
}

func (o *GetDeviceHealthForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetDeviceHealthTooManyRequests creates a GetDeviceHealthTooManyRequests with default headers values
func NewGetDeviceHealthTooManyRequests() *GetDeviceHealthTooManyRequests {
	return &GetDeviceHealthTooManyRequests{}
}

/*GetDeviceHealthTooManyRequests handles this case with default header values.

Too many requests
*/
type GetDeviceHealthTooManyRequests struct {
	Payload *models.ErrorResponse
}

func (o *GetDeviceHealthTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealthTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetDeviceHealthTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceHealthDefault creates a GetDeviceHealthDefault with default headers values
func NewGetDeviceHealthDefault(code int) *GetDeviceHealthDefault {
	return &GetDeviceHealthDefault{
		_statusCode: code,
	}
}

/*GetDeviceHealthDefault handles this case with default header values.

Unexpected error
*/
type GetDeviceHealthDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get device status default response
func (o *GetDeviceHealthDefault) Code() int {
	return o._statusCode
}

func (o *GetDeviceHealthDefault) Error() string {
	return fmt.Sprintf("[GET /health/{deviceId}/health][%d] getDeviceHealth default  %+v", o._statusCode, o.Payload)
}

func (o *GetDeviceHealthDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
