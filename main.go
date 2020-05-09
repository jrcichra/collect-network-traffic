package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jrcichra/influx-network-traffic/influx"

	"github.com/jrcichra/influx-network-traffic/analyzer"
)

func main() {
	var a analyzer.Analyzer
	interval := flag.Int("interval", 60, "Interval you want to capture each interface with")
	sInterfaces := flag.String("interfaces", "", "comma separated list of interfaces to listen on")
	hostname := flag.String("hostname", "influxdb", "hostname/ip of the influxdb")
	db := flag.String("db", "netmetrics", "db in influxdb")
	username := flag.String("username", "", "influx username")
	password := flag.String("password", "", "influx password")
	port := flag.Int("port", 8086, "influx port")

	flag.Parse()

	//Make an influx connection struct
	influxConn := influx.Connection{*hostname, *db, *username, *password, *port}

	if *sInterfaces == "" {
		log.Println("No interfaces were specified! Please specify an interface")
		os.Exit(1)
	} else {
		interfaces := strings.Split(*sInterfaces, ",")
		a.Start(influxConn, *interval, interfaces...)
	}
}
