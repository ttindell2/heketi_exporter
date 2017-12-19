FROM centos:7
LABEL authors="Robert Tindell <Robert@Tindell.info>, CSC Rahti Team <rahti-team@postit.csc.fi>"

EXPOSE 9189

# Copy heketi_exporter
COPY heketi-metrics-exporter /usr/bin/heketi_exporter

CMD /usr/bin/heketi_exporter
