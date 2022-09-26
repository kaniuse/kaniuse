# syntax=docker/dockerfile:experimental

FROM golang:1.19 AS build-env
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /
COPY . /data-scraper
RUN cd /data-scraper && go build -o /data-scraper/data-scraper .

FROM debian:stable-slim AS certs
RUN apt-get update && apt-get -uy upgrade
RUN apt-get -y install ca-certificates && update-ca-certificates

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/kaniuse/kaniuse
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=build-env /data-scraper/data-scraper /data-scraper
ENTRYPOINT ["/data-scraper"]
