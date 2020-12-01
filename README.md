# SmartThings metrics
![ci](https://github.com/moikot/smartthings-metrics/workflows/ci/badge.svg)

A micro-service that provides SmartThings metrics to Prometheus.

## Run 

For this service to have access to SmartThings API you need to provide it with a personal access token (PAT). To generate a PAT do the following:


1. Open SmartThings [Personal access tokens](https://account.smartthings.com/tokens) page.
2. Click "GENERATE NEW TOKEN" button.
3. Give it a name and enable "Devices/List all devices" and "Devices/See all devices" scopes.
4. Click "GENERATE TOKEN" button.

### Run as a standalone app

**Prerequisites:** 
  * [Golang >=1.14](https://golang.org/doc/install)

```bash
$ go get github.com/moikot/smartthings-metrics
$ smartthings-metrics -token [Smarthings-API-token]
```

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
$ helm install smartthings-metrics moikot/smartthings-metrics --create-namespace --names
pace smartthings --set token=[Smarthings-API-token] 
```



 

