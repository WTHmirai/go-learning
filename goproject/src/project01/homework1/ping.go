package main

import (
	"fmt"
	"flag"
	"os"
	"net"
	"time"
	"bytes"
	"encoding/binary"
)

var (
	timeout int  //设置超时时间
	size int //设置报文数据大小
	count int //要发送的回显请求数
	typ uint8 = 8
	code uint8 = 8
)

type ICMP struct {
	Typ uint8
	Code uint8
	CheckSum uint16
	Id uint16
	SequenceNum uint16
}

func GetInput() {
	flag.IntVar(&timeout,"w",1000,"超时时间")
	flag.IntVar(&size,"l",32,"缓冲区大小")
	flag.IntVar(&count,"n",4,"发送请求数")
	flag.Parse()
}

func checkSum(data []byte) uint16 {
	var sum uint32 = 0
	var res uint16
	length := len(data)
	for i := 1 ; i < length ; i+=2 {
		sum += uint32(data[i-1])<<8 + uint32(data[i])
	}
	if length%2 != 0 {
		sum += uint32(data[length-1])
	}
	for ;(sum >> 16) != 0; {
		sum = sum >> 16 + uint32(uint16(sum))
	}
	res = uint16(^sum)
	return res
}

func main() {
	GetInput()

	destIp := os.Args[len(os.Args)-1]   //最后一个参数是ip地址
	
	fmt.Println(destIp)
	Conn,err := net.DialTimeout("ip:icmp",destIp,time.Duration(timeout)*time.Millisecond)

	if err != nil {
		fmt.Println(err)
		return 
	}
	defer Conn.Close()

	fmt.Println(timeout)
	icmp := &ICMP {
		Typ : typ,
		Code : code,
		CheckSum : 0,
		Id : 1,
		SequenceNum : 1,
	}

	data := make([]byte,size)
	var buffer bytes.Buffer  //buffer是icmp报文的二进制格式
	binary.Write(&buffer,binary.BigEndian,icmp) //头部
	buffer.Write(data) //数据部分
	data = buffer.Bytes() //将data换为整个icmp报文，以便后续计算校验和
	checksum := checkSum(data) //校验和
	data[2] = byte(checksum >> 8) //uint16转8取低位，所以右移，大端存储第二位为高位
	data[3] = byte(checksum)

	fmt.Println(data)
	Conn.SetDeadline(time.Now().Add(time.Duration(timeout)*time.Millisecond)) //设置的是绝对期限
	_,err = Conn.Write(data) //向连接中写入数据
	if err != nil {
		fmt.Println(err)
		return
	}
	
	buf := make([]byte,65535) //读取数据的缓冲区
	n,err2 := Conn.Read(buf)
	if err2 != nil {
		fmt.Println(err)
		return 
	}
	fmt.Printf("ip地址:%d.%d.%d.%d,字节:%d,时间:%ds,TTL:%d",buf[12],buf[13],buf[14],buf[15],n-28,0,buf[8])
}