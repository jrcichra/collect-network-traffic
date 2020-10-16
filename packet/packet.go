package packet

//Packet - valuable bits from a packet
type Packet struct {
	Interface string `json:"interface"`
	Bytes     int    `json:"bytes,string,omitempty"`
	SrcName   string `json:"src_name"`
	DstName   string `json:"dst_name"`
	Hostname  string `json:"hostname"`
	Proto     string `json:"proto"`
	SrcPort   int    `json:"scr_port,string,omitempty"`
	DstPort   int    `json:"dst_port,string,omitempty"`
}
