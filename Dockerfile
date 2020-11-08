FROM golang:1.15.4-buster as builder
WORKDIR /app
RUN apt-get update && apt-get install -y libpcap-dev gcc && rm -rf /var/lib/apt/lists/*
COPY . .
RUN go build
FROM debian:buster
WORKDIR     /app
RUN apt-get update && apt-get install -y libpcap-dev && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/collect-network-traffic .
CMD ./collect-network-traffic