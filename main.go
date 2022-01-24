package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

type packet struct {
	Settings       uint8
	Stratum        uint8
	Poll           uint8
	Precision      uint8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimeSec     uint32
	RefTimeFrac    uint32
	OrigTimeSec    uint32
	OrigTimeFrac   uint32
	RxTimeSec      uint32
	RxTimeFrac     uint32
	TxTimeSec      uint32
	TxTimeFrac     uint32
}

var (
	server = flag.String("server", "0.us.pool.ntp.org:123", "The NTP server")
)

func main() {
	flag.Parse()

	conn, err := net.Dial("udp4", *server)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(time.Second * 15)); err != nil {
		panic(err.Error())
	}

	req := &packet{Settings: 0x1B}
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		panic(err.Error())
	}

	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		panic(err.Error())
	}

	// NTP date start from 1900, UNIX start from 1970, so we need to remove 70 years from the NTP packet
	const ntpEpochOffset uint32 = 2208988800
	secs := rsp.TxTimeSec - ntpEpochOffset

	fmt.Println(time.Unix(int64(secs), 0))
}
