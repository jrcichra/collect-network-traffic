FROM golang:1.15.3-alpine3.11 as builder
WORKDIR /app
RUN apk add libpcap-dev build-base
COPY . .
RUN go build
FROM alpine:3.11
WORKDIR     /app
RUN apk add libpcap
COPY --from=builder /app/collect-network-traffic .
CMD ./collect-network-traffic