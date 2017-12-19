# Heketi Metrics Exporter

**NOTE: Directly based on https://github.com/ttindell2/heketi_exporter**

Heketi exporter for Prometheus. Currently only authenticated Heketi servers are supported.

The purpose of this exporter was to get the metrics of a Heketi Controlled Gluster Cluster for Openshift Container Native Storage.

## Usage of `heketi-metrics-exporter`

Help is displayed with `-h`.
| Option                   | Default             | Description
| ------------------------ | ------------------- | -----------------
| -help                    | -                   | Displays usage.
| -listen-address          | `:9189`             | The address to listen on for HTTP requests.
| -log.format              | `logger:stderr`     | Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
| -log.level               | `info`              | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
| -metrics-path            | `/metrics`          | URL Endpoint for metrics
| -version                 | -                   | Prints version information


## Env Variables

### Mandatory
Using the Heketi Go Client to call topology info on authenticated heketi will require these environment variables to be set:

* HEKETI_CLI_SERVER: The Heketi Server (eg. http://localhost:8080)

* HEKETI_CLI_USER: The User should always be 'admin'

* HEKETI_CLI_KEY: The Secret key of the Heketi Admin user.

### Optional environment

* LISTEN_ADDRESS: The address to listen on for HTTP requests


### Metrics in prometheus

| Name                                               | Description                             |
| -------------------------------------------------- | ----------------------------------------|
| up                                                 | Was the last query of Heketi successful |
| cluster_count                                      | Amount of gluster Clusters heketi is controlling |
| volume_count                                       | Amount of volumes in each cluster       |
| node_count                                         | Amount of nodes in each cluster         |
| device_count                                       | Amount of devices mounted to each node in each cluster |
| device_size                                        | Size of each device mounted to each node in each cluster |
| device_free                                        | Free space of each device mounted to each node in each cluster |
| device_used                                        | Used space of each device mounted to each node in each cluster |
| brick_count                                        | Amount of bricks on each device mounted to each node in each cluster |


### Building Locally

Make sure you have make installed and docker installed.

```
make docker
```
