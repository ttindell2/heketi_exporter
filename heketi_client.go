package main

import (
	"github.com/heketi/heketi/client/api/go-client"
	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/prometheus/common/log"
	"os"
)

func TopInfo() (*api.TopologyInfoResponse, error) {
	log.Info("Heketi CLI Server: ", os.Getenv("HEKETI_CLI_SERVER"))
	log.Info("Heketi User: ", os.Getenv("HEKETI_CLI_USER"))
	log.Info("Heketi Key: <REDACTED>")

	// Create a client to talk to Heketi
	heketi := client.NewClient(os.Getenv("HEKETI_CLI_SERVER"), os.Getenv("HEKETI_CLI_USER"), os.Getenv("HEKETI_CLI_KEY"))
	// Create Topology
	topoinfo, err := heketi.TopologyInfo()
	if err != nil {
		log.Errorf("Topology Fail %v",err)
		return nil, err
	}

	return topoinfo, nil
}
