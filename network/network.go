package network

import (
	"net"

	"github.com/jrcichra/influx-network-traffic/dnscache"
)

//Network -
type Network struct {
	cache dnscache.DNSCache
}

//Start - builds a DNS cache
func (n *Network) Start() {
	n.cache = dnscache.DNSCache{}
	n.cache.Start()
}

//GetHostname - returns a hostname when given an IP
func (n *Network) GetHostname(ip string) string {
	c, exists := n.cache.Get(ip)
	var host string
	if exists {
		host = c.(string)
	} else {
		hosts, err := net.LookupAddr(ip)
		if len(hosts) < 1 || err != nil {
			host = ip
		} else {
			host = hosts[0]
		}
		n.cache.Set(ip, host)
	}
	return host
}
