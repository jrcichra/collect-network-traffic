package influx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/influxdata/influxdb/client"
	"github.com/joncrlsn/dque"
	"github.com/jrcichra/influx-network-traffic/packet"
)

//Connection - details to connect to influx
type Connection struct {
	Hostname string
	Db       string
	Username string
	Password string
	Port     int
}

//Influx - wrapper to influx client with state
type Influx struct {
	client     *client.Client
	Connection Connection
}

//PointBuilder - for enqueuing batch points
func PointBuilder() interface{} {
	return &client.BatchPoints{}
}

func makeQueue() *dque.DQue {
	qName := "gps_queue"
	qDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println("qDir=", qDir)
	segmentSize := 50
	if err != nil {
		panic(err)
	}
	q, err := dque.NewOrOpen(qName, qDir, segmentSize, dbRecordBuilder)
	if err != nil {
		panic(err)
	}
	return q
}

//Connect - connect to the influx database
func (f *Influx) Connect(influxConn Connection) {
	f.Connection = influxConn
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", influxConn.Hostname, influxConn.Port))
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

	f.client.SetAuth(influxConn.Username, influxConn.Password)

	log.Println("Connected to influxdb:", influxConn.Hostname, influxConn.Port, "as", influxConn.Username)
}

//Enqueue - queues the influx write
func (f *Influx) Enqueue(measurement string, packet packet.Packet, interval time.Duration, t time.Time) (*client.Response, error) {

	//Convert struct to JSON so we get the keys we want
	jsonB, err := json.Marshal(packet)
	if err != nil {
		panic(err)
	}

	//Convert it into a map[string][string] for the tags
	var m map[string]string
	err = json.Unmarshal(jsonB, &m)
	if err != nil {
		panic(err)
	}

	//Remove bytes
	delete(m, "bytes")
	//Extract out things that shouldn't be tags, but should be fields
	dstPort := m["dst_port"]
	delete(m, "dst_port")
	srcPort := m["src_port"]
	delete(m, "src_port")
	proto := m["proto"]
	delete(m, "proto")

	pt := client.Point{
		Measurement: "throughput",
		Tags:        m,
		Fields: map[string]interface{}{
			"throughput": packet.Bytes,
			"interval":   int(interval.Seconds()),
			"dst_port":   dstPort,
			"src_port":   srcPort,
			"proto":      proto,
		},
		Time: t,
	}

	pts := make([]client.Point, 1)
	pts = append(pts, pt)

	bps := client.BatchPoints{
		Points:          pts,
		Database:        f.Connection.Db,
		RetentionPolicy: "autogen",
	}
	//enqueue

	// Return back the write error
	return f.client.Write(bps)
}
