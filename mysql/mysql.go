package mysql

import (
	"log"
	"time"

	"github.com/jrcichra/collect-network-traffic/packet"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Connection - details to connect
type Connection struct {
	Hostname string
	Db       string
	Username string
	Password string
	Port     int
}

//MySQL -
type MySQL struct {
	client *sql.DB
	stmt   *sql.Stmt
}

//ConnectToDB  -
func (m *MySQL) ConnectToDB(dsn *string) error {
	// connect to the database
	var err error
	m.client, err = sql.Open("mysql", *dsn)
	if err == nil {
		err = m.client.Ping()
	}
	// prepare the sole query
	m.stmt, err = m.client.Prepare("insert throughput set interface = ?, bytes = ?, src_name = ?, dst_name = ?, hostname = ?, proto = ?, src_port = ?, dst_port = ?, `interval` = ?")
	if err != nil {
		log.Println(err)
		//reconnect to the db
		m.client.Close()
		m.ReconnectToDB(dsn)
	}
	return err
}

//ReconnectToDB -
func (m *MySQL) ReconnectToDB(dsn *string) {
	var err error
	err = nil
	for err == nil {
		err = m.ConnectToDB(dsn)
		// prepare the sole query
		m.stmt, err = m.client.Prepare("insert throughput set interface = ?, bytes = ?, src_name = ?, dst_name = ?, hostname = ?, proto = ?, src_port = ?, dst_port = ?, `interval` = ?")
		if err != nil {
			log.Println(err)
			m.client.Close()
			time.Sleep(time.Duration(1) * time.Second)
			err = nil
		}
	}
}

//Insert - inserts row into database
func (m *MySQL) Insert(p *packet.Packet, interval time.Duration) {
	_, err := m.stmt.Exec(p.Interface, p.Bytes, p.SrcName, p.DstName, p.Hostname, p.Proto, p.SrcPort, p.DstPort, int(interval.Seconds()))
	if err != nil {
		panic(err)
	}
}
