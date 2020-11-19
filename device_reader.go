package main

import (
	"github.com/moikot/smartthings-metrics/health"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/client"
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/client/devices"
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

type DeviceStatus struct {
	Device *models.Device
	Status *models.DeviceStatus
	Health *health.DeviceHealth
}

type DeviceReader interface {
	ReadStatuses() ([]*DeviceStatus, error)
}

func NewDeviceReader(token string, log logrus.FieldLogger) DeviceReader {
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
		SmartThings: client.New(rtm, nil),
		Client:      health.New(rtm, nil),
		auth:        NewAuthInfoWriter(token),
		log:         log,
	}
}

type statusReader struct {
	*client.SmartThings
	*health.Client
	auth runtime.ClientAuthInfoWriter
	log logrus.FieldLogger
}

func (c *statusReader) ReadStatuses() ([]*DeviceStatus, error) {
	devs, err := c.Devices.GetDevices(nil, c.auth)
	if err != nil {
		return nil, err
	}
	var res []*DeviceStatus
	for _, dev := range devs.Payload.Items {
		p := health.NewGetDeviceHealthParams()
		p.DeviceID = dev.DeviceID
		health, err := c.GetDeviceHealth(p, c.auth)
		if err != nil {
			return nil, err
		}

		params := devices.NewGetDeviceStatusParams()
		params.DeviceID = dev.DeviceID
		status, err := c.Devices.GetDeviceStatus(params, c.auth)
		if err != nil {
			return nil, err
		}

		devStatus := &DeviceStatus{
			Device: dev,
			Health: health.Payload,
			Status: status.Payload,
		}
		res = append(res, devStatus)
	}
	return res, nil
}
