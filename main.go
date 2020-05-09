package main

import (
	"os"
	"strconv"

	"github.com/jrcichra/influx-network-traffic/analyzer"
)

func main() {
	var a analyzer.Analyzer
	interval, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	a.Start(interval, "wlan0")
}
