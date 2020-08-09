FROM alpine:3.12 AS build

RUN apk --no-cache add build-base go

ADD . /go/src/github.com/google/systemd_exporter
WORKDIR /go/src/github.com/google/systemd_exporter

ENV GOROOT /usr/lib/go
ENV GOPATH /go
RUN make build

FROM alpine:3.12
MAINTAINER dashpole@google.com

RUN apk --no-cache add libc6-compat

ADD systemd_exporter /usr/bin/systemd_exporter

EXPOSE 8080

ENTRYPOINT ["/usr/bin/systemd_exporter", "--logtostderr"]