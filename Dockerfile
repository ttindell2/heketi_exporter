FROM centos:7
MAINTAINER Robert Tindell <Robert@Tindell.info>

EXPOSE 9189

# Copy heketi_exporter
COPY heketi_exporter /usr/bin/heketi_exporter

CMD /usr/bin/heketi_exporter
