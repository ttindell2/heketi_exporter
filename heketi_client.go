package main

import (
	"fmt"
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
	// Get the cluster list and iterate over
	for i, _ := range topoinfo.ClusterList {
		fmt.Fprintf(os.Stdout, "\nCluster Id: %v\n", topoinfo.ClusterList[i].Id)
		fmt.Fprintf(os.Stdout, "\n    File:  %v\n", topoinfo.ClusterList[i].File)
		fmt.Fprintf(os.Stdout, "    Block: %v\n", topoinfo.ClusterList[i].Block)
		fmt.Fprintf(os.Stdout, "\n    %s\n", "Volumes:")
		for k, _ := range topoinfo.ClusterList[i].Volumes {

			// Format and print volumeinfo  on this cluster
			v := topoinfo.ClusterList[i].Volumes[k]
			s := fmt.Sprintf("\n\tName: %v\n"+
				"\tSize: %v\n"+
				"\tId: %v\n"+
				"\tCluster Id: %v\n"+
				"\tMount: %v\n"+
				"\tMount Options: backup-volfile-servers=%v\n"+
				"\tDurability Type: %v\n",
				v.Name,
				v.Size,
				v.Id,
				v.Cluster,
				v.Mount.GlusterFS.MountPoint,
				v.Mount.GlusterFS.Options["backup-volfile-servers"],
				v.Durability.Type)

			switch v.Durability.Type {
			case api.DurabilityEC:
				s += fmt.Sprintf("\tDisperse Data: %v\n"+
					"\tDisperse Redundancy: %v\n",
					v.Durability.Disperse.Data,
					v.Durability.Disperse.Redundancy)
			case api.DurabilityReplicate:
				s += fmt.Sprintf("\tReplica: %v\n",
					v.Durability.Replicate.Replica)
			}
			if v.Snapshot.Enable {
				s += fmt.Sprintf("\tSnapshot: Enabled\n"+
					"\tSnapshot Factor: %.2f\n",
					v.Snapshot.Factor)
			} else {
				s += "\tSnapshot: Disabled\n"
			}
			s += "\n\t\tBricks:\n"
			for _, b := range v.Bricks {
				s += fmt.Sprintf("\t\t\tId: %v\n"+
					"\t\t\tPath: %v\n"+
					"\t\t\tSize (GiB): %v\n"+
					"\t\t\tNode: %v\n"+
					"\t\t\tDevice: %v\n\n",
					b.Id,
					b.Path,
					b.Size/(1024*1024),
					b.NodeId,
					b.DeviceId)
			}
			fmt.Fprintf(os.Stdout, "%s", s)
		}

		// format and print each Node information on this cluster
		fmt.Fprintf(os.Stdout, "\n    %s\n", "Nodes:")
		for j, _ := range topoinfo.ClusterList[i].Nodes {
			info := topoinfo.ClusterList[i].Nodes[j]
			fmt.Fprintf(os.Stdout, "\n\tNode Id: %v\n"+
				"\tState: %v\n"+
				"\tCluster Id: %v\n"+
				"\tZone: %v\n"+
				"\tManagement Hostname: %v\n"+
				"\tStorage Hostname: %v\n",
				info.Id,
				info.State,
				info.ClusterId,
				info.Zone,
				info.Hostnames.Manage[0],
				info.Hostnames.Storage[0])
			fmt.Fprintf(os.Stdout, "\tDevices:\n")

			// format and print the device info
			for j, d := range info.DevicesInfo {
				fmt.Fprintf(os.Stdout, "\t\tId:%-35v"+
					"Name:%-20v"+
					"State:%-10v"+
					"Size (GiB):%-8v"+
					"Used (GiB):%-8v"+
					"Free (GiB):%-8v\n",
					d.Id,
					d.Name,
					d.State,
					d.Storage.Total/(1024*1024),
					d.Storage.Used/(1024*1024),
					d.Storage.Free/(1024*1024))

				// format and print the brick information
				fmt.Fprintf(os.Stdout, "\t\t\tBricks:\n")
				for _, d := range info.DevicesInfo[j].Bricks {
					fmt.Fprintf(os.Stdout, "\t\t\t\tId:%-35v"+
						"Size (GiB):%-8v"+
						"Path: %v\n",
						d.Id,
						d.Size/(1024*1024),
						d.Path)
				}
			}
		}
	}

	return topoinfo, nil
}
