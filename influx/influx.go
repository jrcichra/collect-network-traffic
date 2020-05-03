package influx

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"

	"github.com/jrcichra/influx-network-traffic/packet"

	"github.com/influxdata/influxdb/client"
)

//Influx - wrapper to influx client with state
type Influx struct {
	client *client.Client
}

//Connect - connect to the influx database
func (f *Influx) Connect(host, port, username, password string) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}

	f.client, err = client.NewClient(client.Config{URL: *u})
	if err != nil {
		log.Fatal(err)
	}

	if _, _, err := f.client.Ping(); err != nil {
		log.Fatal(err)
	}

	f.client.SetAuth(username, password)
}

//Write - writes just as the script does
func (f *Influx) Write(measurement string, packet packet.Packet) (*client.Response, error) {
	fmt.Println("Inserting data...")

	//Convert struct to map
	m := structs.Map(packet)
	//make everything in here a string
	m2 := make(map[string]string)
	for k, v := range m {
		switch temp := v.(type) {
		case string:
			m2[strings.ToLower(k)] = strings.ToLower(temp)
		case int:
			m2[strings.ToLower(k)] = strconv.Itoa(temp)
		}
	}

	//Remove bytes as it's not a tag but a field (the only one)
	delete(m, "Bytes")

	pt := client.Point{
		Measurement: "throughput",
		Tags:        m2,
		Fields: map[string]interface{}{
			"throughput": packet.Bytes,
		},
		Time: time.Now(),
	}

	pts := make([]client.Point, 1)
	pts = append(pts, pt)

	bps := client.BatchPoints{
		Points:          pts,
		Database:        "netmetrics",
		RetentionPolicy: "autogen",
	}
	fmt.Println(f.client.Write(bps))
	return nil, nil
}
