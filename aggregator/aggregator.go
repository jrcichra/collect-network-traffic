package aggregator

import (
	"time"

	"github.com/jrcichra/collect-network-traffic/network"

	"github.com/jrcichra/collect-network-traffic/mysql"

	"github.com/jrcichra/collect-network-traffic/packet"
)

//Aggregator - takes packets and collects them until they are flushed on an interval
type Aggregator struct {
	interval     time.Duration
	packetBuffer map[string]packet.Packet
	mysql        *mysql.MySQL
	networkUtils network.Network
}

func (g *Aggregator) inserter(request chan struct{}, response chan map[packet.Packet]int) {
	//Request the data
	request <- struct{}{}
	//Collect the response
	packets := <-response
	//Put it in the database
	for p, bytes := range packets {
		//p is the packet, bytes is the bytes aggregated for this connection/time
		//Fix up the bytes
		p.Bytes = bytes
		// sum += bytes
		//Divide the number of bytes by the number of seconds in this interval
		p.Bytes /= int(g.interval.Seconds())
		//Replace IPs with hostnames
		p.SrcName = g.networkUtils.GetHostname(p.SrcName)
		p.DstName = g.networkUtils.GetHostname(p.DstName)
		//Insert it into mysql
		g.mysql.Insert(&p, g.interval)
	}
}

//called as a goroutine that starts an insert
func (g *Aggregator) insertTimer(interval time.Duration, request chan struct{}, response chan map[packet.Packet]int) {
	g.interval = interval
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			//we should insert in the database
			go g.inserter(request, response)
		}

	}
}

//Start - takes an aggregation interval and a channel of packets to aggregate
func (g *Aggregator) Start(interval time.Duration, packetChan chan packet.Packet, m *mysql.MySQL) {
	//Get a network object for the aggregator (getHostname)
	g.networkUtils = network.Network{}
	g.networkUtils.Start()
	//give this object the mysql struct
	g.mysql = m
	//the timer will request the packet cache on a given interval with no meaningful data
	request := make(chan struct{})
	//the response will be a copy of the packet cache for this given interval
	response := make(chan map[packet.Packet]int)
	//Set a go ticker for the given interval. When this hits we should flush our data to the database
	go g.insertTimer(interval, request, response)
	//Make a cache of Packet structs with the value it points to as the true aggregate value (ignore the packet.Packet size, should be zero here)
	packetCache := make(map[packet.Packet]int)
	//Handle incoming packets
	for p := range packetChan {
		//Check if we have a request for the current cache from the inserter
		select {
		case <-request:
			//we got a request! send the packet cache over and flush our copy
			response <- packetCache
			packetCache = make(map[packet.Packet]int)
		default:
			//Didn't get a request, continue building the cache
		}

		//Hold onto the size
		bytes := p.Bytes
		//Set the size to zero to match up the rest of the keys
		p.Bytes = 0
		//Clear any other fields that don't want to be aggregated
		//TODO: bubble up this feature to the user
		//Determine if it has the same attributes as previous packets
		if _, ok := packetCache[p]; ok {
			//We found the key, update
			packetCache[p] += bytes
		} else {
			//It's not here, create
			packetCache[p] = bytes
		}

	}
}
