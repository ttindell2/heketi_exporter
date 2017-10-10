# heketi_exporter
Heketi exporter for Prometheus

## Installation

```
go get github.com/ttindell2/heketi_exporter
./heketi_exporter
```

## Usage of `heketi_exporter`
Help is displayed with `-h`.

| Option                   | Default             | Description
| ------------------------ | ------------------- | -----------------
| -help                    | -                   | Displays usage.
| -listen-address          | `:9189`             | The address to listen on for HTTP requests.
| -log.format              | `logger:stderr`     | Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
| -log.level               | `info`              | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
| -metrics-path            | `/metrics`          | URL Endpoint for metrics
| -version                 | -                   | Prints version information


## Env Variables to Set:

HEKETI_CLI_SERVER
HEKETI_CLI_USER
HEKETI_CLI_KEY

### Command: `heketi topology info`

| Name                                               | type     | impl. state |
| -------------------------------------------------- | -------- | ------------|

#TODO

