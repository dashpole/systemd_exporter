FROM alpine:3.12
MAINTAINER dashpole@google.com

COPY systemd_exporter /usr/bin/systemd_exporter

EXPOSE 8080

ENTRYPOINT ["/usr/bin/systemd_exporter", "-logtostderr"]