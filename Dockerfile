FROM golang:1.15.3-alpine3.12 as builder
WORKDIR /app
RUN apk add libpcap-dev
COPY . .
RUN go build
FROM alpine:3.12
WORKDIR     /app
RUN apk add libpcap
COPY --from=builder /app/influx-network-traffic .
CMD ./influx-network-traffic