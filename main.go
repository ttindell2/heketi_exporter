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

// heketi exporter, exports metrics using Heketi Go Client
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
		"Was the last query of heketi successful.",
		nil, nil,
	)
	clusterCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_count"),
		"Number of clusters at last query.",
		nil, nil,
	)
	volumesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volumes_count"),
		"How many volumes were up at the last query.",
		[]string{"cluster"}, nil,
	)
	nodesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nodes_count"),
		"How many Nodes were up at the last query.",
		[]string{"cluster"}, nil,
	)
	deviceCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "device_count"),
		"How many Devices were up at the last query.",
		[]string{"cluster","hostname"}, nil,
	)
	deviceSize = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "device_size"),
		"How many Devices were up at the last query.",
		[]string{"cluster","hostname", "device"}, nil,
	)
	deviceFree = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "device_free"),
		"How many Devices were up at the last query.",
		[]string{"cluster","hostname", "device"}, nil,
	)
	deviceUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "device_used"),
		"How many Devices were up at the last query.",
		[]string{"cluster","hostname", "device"}, nil,
	)
	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_count"),
		"Number of bricks at last query.",
		[]string{"cluster","hostname", "device"}, nil,
	)
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	hostname string
}

// Describe all the metrics exported by Heketi exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up // done
	ch <- clusterCount // done
	ch <- volumesCount // done
	ch <- nodesCount // done
	ch <- deviceCount // done
	ch <- deviceSize // done
	ch <- deviceFree // done
	ch <- deviceUsed // done
	ch <- brickCount // done

}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Collect metrics from volume info
	topinfo, err := TopInfo()
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		log.Errorf("couldn't get topology info: %v", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		clusterCount, prometheus.GaugeValue, float64(len(topinfo.ClusterList)),
	)
	for _, cluster := range topinfo.ClusterList {
		log.Info("ClusterID: ", cluster.Id)
                ch <- prometheus.MustNewConstMetric(
			volumesCount, prometheus.GaugeValue, float64(len(cluster.Volumes)),cluster.Id,
		)
		ch <- prometheus.MustNewConstMetric(
		        nodesCount, prometheus.GaugeValue, float64(len(cluster.Nodes)),cluster.Id,
		)
	//	for _, volumes := range cluster.Volumes {
	//			// Not Using for now
	//	}
		for _, nodes := range cluster.Nodes {
                        log.Info("NodeHost: ",nodes.Hostnames.Manage[0])
			ch <- prometheus.MustNewConstMetric(
				deviceCount, prometheus.GaugeValue, float64(len(nodes.DevicesInfo)),cluster.Id,nodes.Hostnames.Manage[0],
			)
			for _, device := range nodes.DevicesInfo {
                                log.Info("Device: ", device.Name)
				ch <- prometheus.MustNewConstMetric(
					deviceSize, prometheus.GaugeValue, float64(device.Storage.Total),cluster.Id,nodes.Hostnames.Manage[0], device.Name,
				)
				ch <- prometheus.MustNewConstMetric(
					deviceFree, prometheus.GaugeValue, float64(device.Storage.Free),cluster.Id,nodes.Hostnames.Manage[0], device.Name,
				)
				ch <- prometheus.MustNewConstMetric(
					deviceUsed, prometheus.GaugeValue, float64(device.Storage.Used),cluster.Id,nodes.Hostnames.Manage[0], device.Name,
				)
				ch <- prometheus.MustNewConstMetric(
					brickCount, prometheus.GaugeValue, float64(len(device.Bricks)),cluster.Id,nodes.Hostnames.Manage[0], device.Name,
				)
			}
		}
	}
   log.Info("Finished collecting metrics")
}

// NewExporter initialises exporter
func NewExporter(hostname string) (*Exporter, error) {
	return &Exporter{
		hostname: hostname,
	}, nil
}

func versionInfo() {
	fmt.Println(version.Print("heketi_exporter"))
	os.Exit(0)
}

func init() {
	prometheus.MustRegister(version.NewCollector("heketi_exporter"))
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
			<head><title>heketiFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>heketiFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
