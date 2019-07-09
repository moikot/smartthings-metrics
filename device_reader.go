package main

import (
	"net/http"
	"time"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/client"
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/client/devices"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

type DeviceStatus struct {
	Device *models.Device
	Status *models.DeviceStatus
}

type DeviceReader interface {
	Read() ([]*DeviceStatus, error)
}

func NewDeviceReader(token string) DeviceReader {
	httpClient := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	cfg := client.DefaultTransportConfig()

	rtm := httptransport.NewWithClient(
		cfg.Host,
		cfg.BasePath,
		cfg.Schemes,
		httpClient,
	)

	return &statusReader{
		Client: client.New(rtm, nil).Devices,
		auth:   NewAuthInfoWriter(token),
	}
}

type statusReader struct {
	*devices.Client
	auth runtime.ClientAuthInfoWriter
}

// TODO: Request statuses for devices with supported capabilities only
func (c *statusReader) Read() ([]*DeviceStatus, error) {
	devs, err := c.GetDevices(nil, c.auth)
	if err != nil {
		return nil, err
	}
	var res []*DeviceStatus
	for _, dev := range devs.Payload.Items {
		params := devices.NewGetDeviceStatusParams()
		params.DeviceID = dev.DeviceID

		status, err := c.GetDeviceStatus(params, c.auth)
		if err != nil {
			return nil, err
		}

		devStatus := &DeviceStatus{
			Device: dev,
			Status: status.Payload,
		}
		res = append(res, devStatus)
	}
	return res, nil
}
