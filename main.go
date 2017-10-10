// Copyright 2015 Oliver Fesseler
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for 	the specific language governing permissions and
// limitations under the License.

// Gluster exporter, exports metrics from gluster commandline tool.
package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"net/http"
	"os"
)

const (
	namespace = "gluster"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Gluster successful.",
		nil, nil,
	)
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	hostname string
}

// Describe all the metrics exported by Gluster exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Collect metrics from volume info
	topinfo, err := TopInfo()
	//TopologyInfo()
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		log.Errorf("couldn't get topology info: %v", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}
	log.Info(topinfo)
}

// NewExporter initialises exporter
func NewExporter(hostname string) (*Exporter, error) {
	return &Exporter{
		hostname: hostname,
	}, nil
}

func versionInfo() {
	fmt.Println(version.Print("gluster_exporter"))
	os.Exit(0)
}

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}

func main() {

	// commandline arguments
	var (
		metricPath    = flag.String("metrics-path", "/metrics", "URL Endpoint for metrics")
		listenAddress = flag.String("listen-address", ":9189", "The address to listen on for HTTP requests.")
		showVersion   = flag.Bool("version", false, "Prints version information")
	)
	flag.Parse()

	if *showVersion {
		versionInfo()
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("While trying to get Hostname error happened: %v", err)
	}
	exporter, err := NewExporter(hostname)
	if err != nil {
		log.Errorf("Creating new Exporter went wrong, ... \n%v", err)
	}
	prometheus.MustRegister(exporter)

	log.Info("Heketi Metrics Exporter v", version.Version)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>GlusterFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>GlusterFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
