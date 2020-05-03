package main

import "github.com/jrcichra/influx-network-traffic/analyzer"

func main() {
	var a analyzer.Analyzer
	a.Start("wlan0")
}
