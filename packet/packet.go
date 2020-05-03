package packet

//Packet - valuable bits from a packet
type Packet struct {
	Interface string `json:"interface"`
	Bytes     int    `json:"bytes"`
	SrcIP     string `json:"src_ip"`
	DstIP     string `json:"dst_ip"`
	Proto     string `json:"proto"`
	SrcPort   int    `json:"scr_port"`
	DstPort   int    `json:"dst_port"`
}
