package main

import(
	"fmt"
	"encoding/binary"
	"bytes"
	"strings"
	"net"
)

var (
	DnsServer string = "192.168.1.1:53"
	Domain string = "www.baidu.com"
)

type DNS struct {
	ID uint16
	Flags uint16
	Questions uint16
	AnRRs uint16
	AuRRs uint16
	AdRRs uint16
} 

func SetFlag(qr uint16,Opcode uint16,AA uint16,TC uint16,RD uint16,RA uint16,rcode uint16) uint16 {
	var flag uint16
	flag = qr << 15 + Opcode << 11 + AA << 10 + TC << 9 + RD << 8 + RA << 7 + rcode
	return flag
}

func TransformDomain(domain string) []byte {
	var(
		query bytes.Buffer
		seg = strings.Split(domain,".")
	)
	for i := 0; i < len(seg) ; i++ {
		binary.Write(&query,binary.BigEndian,byte(len(seg[i])))
		binary.Write(&query,binary.BigEndian,[]byte(seg[i]))
	}
	binary.Write(&query,binary.BigEndian,[]byte{0,0,1,0,1})
	return query.Bytes()
}

func main() {
	header := DNS {
		ID : 0xFFFF,
		Flags : SetFlag(0,0,0,0,1,0,0),
		Questions : 1,
		AnRRs : 0,
		AuRRs : 0,
		AdRRs : 0,
	}

	var(
		conn net.Conn
		err error
		n int
	)

	var buffer bytes.Buffer
	binary.Write(&buffer,binary.BigEndian,header)
	binary.Write(&buffer,binary.BigEndian,TransformDomain(Domain))
	if conn,err = net.Dial("udp",DnsServer) ; err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	if _,err = conn.Write(buffer.Bytes()) ; err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buffer.Bytes())
	rebuffer := make([]byte,1024)
	if n,err = conn.Read(rebuffer) ; err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(n,rebuffer)
}
