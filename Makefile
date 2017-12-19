GO = go
BIN_DIR ?= $(shell pwd)

info:
	@echo "build: Go build"
	@echo "docker: build and run in docker container"
	@echo "gotest: run go tests and reformats"

build:
	go get -u "github.com/prometheus/client_golang/prometheus"
	go get -u "github.com/prometheus/client_golang/prometheus/promhttp"
	go get -u "github.com/prometheus/common/log"
	go get -u "github.com/prometheus/common/version"
	go get -u "github.com/heketi/heketi/client/api/go-client"
	go get -u "github.com/heketi/heketi/pkg/glusterfs/api"
	$(GO) build -o heketi-metrics-exporter

docker: build
	docker build -t heketi-metrics-exporter .
