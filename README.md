# SmartThings metrics
![ci](https://github.com/moikot/smartthings-metrics/workflows/ci/badge.svg)

A micro-service that provides SmartThings metrics to Prometheus.

## Run

For this service to have access to SmartThings API, you need to provide it with a personal access token (PAT). To generate a PAT, do the following:

1. Open SmartThings [Personal access tokens](https://account.smartthings.com/tokens) page.
2. Click | GENERATE NEW TOKEN|  button.
3. Give it a name and enable | Devices/List all devices|  and | Devices/See all devices|  scopes.
4. Click | GENERATE TOKEN|  button.

### Run as a standalone app

**Prerequisites:**
  * [Golang >=1.14](https://golang.org/doc/install)

```bash
$ go get github.com/moikot/smartthings-metrics
$ smartthings-metrics -token [Smarthings-API-token]
$ curl localhost:9153/metrics
```

**Note:** Using `-interval` you can define the refresh interval in seconds. The default value of the refresh interval is 60 seconds.

### Run as a Docker container

**Prerequisites:**
  * [Docker](https://docs.docker.com/get-docker/)

```bash
$ docker run -d --rm -p 9153:9153 moikot/smartthings-metrics -token [Smarthings-API-token]
$ curl localhost:9153/metrics
```

### Deploy to a Kubernetes cluster

**Prerequisites:**
  * [Kuberentes](https://kubernetes.io/)
  * [Helm 3](https://helm.sh)

SmartThing metrics service is installed to Kubernetes via its [Helm chart](https://github.com/moikot/helm-charts/tree/master/charts/smartthings-metrics).

```
$ helm repo add moikot https://moikot.github.io/helm-charts
$ helm install smartthings-metrics moikot/smartthings-metrics --create-namespace --namespace smartthings --set token=[Smarthings-API-token]
```

## How it works

The service uses [SmartThings API](https://smartthings.developer.samsung.com/docs/api-ref/st-api.html) to obtain the current status of all connected devices periodically. It exposes the metrics received at `localhost:9153/metrics` so that [Prometheus](https://prometheus.io/) could scrape them.

### Gauges

Metric are exposed as Prometheus [gauges](https://prometheus.io/docs/concepts/metric_types/#gauge) their names are formed using pattern `smartthings_[component_name]_[capability_name]_[attribute_name]_[measurement_unit]` with all individual names converted to the snake case.

**Examples:**
 * `smartthings_motion_sensor_motion`
 * `smartthings_battery_battery_percent`
 * `smartthings_power_meter_power_watt`

**Note:**  
  * The component name is used unless it the `main` component.
  * The measurement unit is added for the values with units only.
  * Gauge with the name  `smartthings_health_state` is used for the health probe.

### Values

All the attributes of type `number` and `integer` are supported. Attributes of type `string` are supported only for enumerations. The service maps all of the enumeration identifiers to numbers and then uses them to translate a string value to a number. For example [switch](https://smartthings.developer.samsung.com/docs/api-ref/capabilities.html#Switch) capability defines two enumeration identifiers `on` and `off`. They will be mapped to `0.0` and `1.0`.

**Note:** For `smartthings_health_state` gauge uses this static mapping:

| Enumeration value | Numeric value |
|-------------------|---------------|
| OFFLINE |   0.0 |
| UNHEALTHY | 1.0 |
| ONLINE |    2.0 |

### Labels

Several labels are added for all measurements:

| Label | Description |
|-------|-------------|
| `name` | The name of the device. |
| `label` | The device name provided by a user (the device name by default). |
| `device_id` | The unique identifier for the device instance. |
| `location_id` | The identifier of the location the device is associated with. |
| `device_manufacturer_code` | The device manufacturer code. |

For the devices with a Device Type Handler (DTH), several additional labels are provided:

| Label | Description |
|-------|-------------|
| `device_type_name` | The name for the device type. |
| `device_type_id` | The identifier of the device type. |
| `device_network_type` | The device network type. |

### Units

The units of measurement are mapped according to this table:

| Symbol | Unit name |
|-------|-------------|
| % | percent |
| lux | lux |
| s | second |
| W | watt |
| C | degree_celsius |
| K | degree_kelvin |
| F | degree_fahrenheit |
| V | volt |
| kg | kilogram |
| lbs | pound |
| 斤 | catty |
| CAQI | caqi |
| dBm | decibel_milliwatt |
| μg/m^3 | microgram_per_cubic_meter |
| mg/m^3 | milligram_per_cubic_meter |
| kg/m^2 | kilogram_per_square_meter |
| kWh | kilowatt_hour |
| ppm | parts_per_million |
