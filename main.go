package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jrcichra/collect-network-traffic/analyzer"
	"github.com/jrcichra/collect-network-traffic/mysql"
)

func main() {
	var a analyzer.Analyzer
	interval := flag.Int("interval", 60, "Interval you want to capture each interface with")
	sInterfaces := flag.String("interfaces", "", "comma separated list of interfaces to listen on")
	dsn := flag.String("dsn", "", "The connect string for your database - see https://github.com/go-sql-driver/mysql#dsn-data-source-name")

	flag.Parse()
	m := &mysql.MySQL{}
	err := m.ConnectToDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if *sInterfaces == "" {
		log.Println("No interfaces were specified! Please specify an interface")
		os.Exit(1)
	} else {
		interfaces := strings.Split(*sInterfaces, ",")
		a.Start(m, *interval, interfaces...)
	}
}
