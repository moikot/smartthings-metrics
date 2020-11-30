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
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/moikot/smartthings-go/client"
	"github.com/moikot/smartthings-go/client/devices"
	"github.com/moikot/smartthings-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Component struct {
	Capabilities map[string]*models.Capability
}

type Schema struct {
	Components map[string]*Component
}

type DeviceStatus struct {
	Device *models.Device
	Status *models.DeviceStatus
	Health *DeviceHealth
	Schema *Schema
}

type DeviceReader interface {
	ReadStatuses() ([]*DeviceStatus, error)
}

func NewDeviceReader(token string, log logrus.FieldLogger) DeviceReader {
	auth := NewAuthInfoWriter(token)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	cfg := client.DefaultTransportConfig()

	rtm := httptransport.NewWithClient(
		cfg.Host,
		cfg.BasePath,
		cfg.Schemes,
		httpClient,
	)

	return &statusReader{
		SmartThingsAPI:   client.New(rtm, nil),
		DeviceHealthAPI:  NewDeviceHealthAPI(rtm, nil),
		CapabilityReader: NewCapabilityReader(rtm, nil),
		auth:             auth,
		log:              log,
	}
}

type statusReader struct {
	*client.SmartThingsAPI
	*DeviceHealthAPI
	*CapabilityReader
	auth runtime.ClientAuthInfoWriter
	log  logrus.FieldLogger
}

func (r *statusReader) ReadStatuses() ([]*DeviceStatus, error) {
	devs, err := r.Devices.GetDevices(nil, r.auth)
	if err != nil {
		return nil, err
	}
	var res []*DeviceStatus
	for _, dev := range devs.Payload.Items {

		p := NewGetDeviceHealthParams()
		p.DeviceID = *dev.DeviceID
		health, err := r.GetDeviceHealth(p, r.auth)
		if err != nil {
			r.log.Errorf("failed to health status of device '%s': %v", dev.DeviceID, err)
			continue
		}

		params := devices.NewGetDeviceStatusParams()
		params.DeviceID = *dev.DeviceID
		status, err := r.Devices.GetDeviceStatus(params, r.auth)
		if err != nil {
			r.log.Errorf("failed to read status of device '%s': %v", dev.DeviceID, err)
			continue
		}

		schema := &Schema{
			Components: make(map[string]*Component),
		}

		for _, component := range dev.Components {
			comp := schema.Components[*component.ID]
			if comp == nil {
				comp = &Component{
					Capabilities: make(map[string]*models.Capability),
				}
				schema.Components[*component.ID] = comp
			}
			for _, capability := range component.Capabilities {
				c, err := r.ReadCapability(*capability.ID, capability.Version, r.auth)
				if err != nil {
					r.log.Errorf("failed to read capability '%s' of device '%s': %v", capability.ID, dev.DeviceID, err)
					continue
				}
				comp.Capabilities[*capability.ID] = c
			}
		}

		devStatus := &DeviceStatus{
			Device: dev,
			Health: health.Payload,
			Status: status.Payload,
			Schema: schema,
		}
		res = append(res, devStatus)
	}
	return res, nil
}
