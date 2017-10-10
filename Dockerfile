FROM alpine:latest
MAINTAINER Robert Tindell <Robert@Tindell.info>

EXPOSE 9189

# Copy heketi_exporter
COPY heketi_exporter /usr/bin/heketi_exporter

ENTRYPOINT /usr/bin/heketi_exporter
