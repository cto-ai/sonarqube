FROM golang:1.13.8 AS build-ss

ENV DEBIAN_FRONTEND=noninteractive
ENV SONARSCANNER_VERSION 4.2.0.1873-linux
RUN apt update && apt install -y binutils wget unzip

WORKDIR /ss
RUN wget --quiet https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-${SONARSCANNER_VERSION}.zip -O ss.zip && \
    unzip -qq ss.zip && \
    mv sonar-scanner-${SONARSCANNER_VERSION}/ /sonar-scanner

############################
# Build UI container
############################
FROM golang:1.13.8 AS build-ui

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y binutils

WORKDIR /ops

ADD . .
RUN go build -ldflags='-s -w' -o main && strip -S main && chmod 777 main
# for debugging
# RUN go build -o main

############################
# Final container
############################
FROM registry.cto.ai/official_images/base:latest

ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /ops/proj

RUN apt update && apt install -y --reinstall openssl libssl1.1 ca-certificates && rm -rf /var/lib/apt/lists && ln -s /bin/bin/sonar-scanner /bin/sonar-scanner
COPY --from=build-ss /sonar-scanner /bin/
COPY --from=build-ui /ops/main /ops/main
RUN chown -R 9999:9999 /home/ops && chown -R 9999:9999 /ops