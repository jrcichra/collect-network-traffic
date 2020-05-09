# influx-network-traffic
Store network traffic information in influxdb

![Preview](preview.png "Preview")

# Usage
```bash
Usage of ./influx-network-traffic:
  -db string
        db in influxdb (default "netmetrics")
  -hostname string
        hostname/ip of the influxdb (default "influxdb")
  -interfaces string
        comma separated list of interfaces to listen on
  -interval int
        Interval you want to capture each interface with (default 60)
  -password string
        influx password
  -port int
        influx port (default 8086)
  -username string
        influx username
```

# Getting started

1. Go to Releases and get the latest copy of influx-network-traffic for your architecture
2. Install libpcap for your distro: (i.e) `sudo apt install libpcap`
3. Install influxDB & Chronograf (or use my docker-compose for something quick) `docker-compose up -d`
4. Goto port `8888` on the system you installed influx/chronograf on
5. Set it up to your influxdb credentials (docker-compose is admin/admin, hostname influxdb)
6. Skip importing charts and kapacitor
7. Goto the admin panel and set your retention policy for netmetrics
8. Goto your grafana install. Make a new datasource for influx
9. Import `grafana.json` and add a new datasource for your influxdb database. Make sure the datasource name matches with the JSON
10. Run influx-network-traffic similar to this: `sudo ./influx-network-traffic -interfaces wlan0 -interval 10 -hostname localhost -username admin -password admin -db netmetrics`
11. You should have a screen similar to the preview. Play with the query to make awesome charts and make PRs!