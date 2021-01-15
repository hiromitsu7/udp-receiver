package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 514,
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// RFC 3164 BSD syslogでは1024バイトが上限
	buf := make([]byte, 1024)
	fmt.Println("starting UDP server...")

	for {
		n, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		payload := string(buf[0:n])
		indexPriStart := strings.Index(payload, "<")
		indexPriEnd := strings.Index(payload, ">")

		pri := payload[indexPriStart+1 : indexPriEnd]
		timestamp := payload[indexPriEnd+1 : indexPriEnd+16]
		// ホスト名、IPアドレスはないものとする

		message := payload[indexPriEnd+17:]
		indexTagSep := strings.Index(message, ":")
		tag := message[:indexTagSep]
		content := message[indexTagSep+2:]

		if strings.Contains(content, "58:6b:14") || strings.Contains(content, "7c:04:d0") || strings.Contains(content, "5c:f9:38") {
			content += "(Apple)"
		} else if strings.Contains(content, "b4:cd:27") {
			content += "(Huawei)"
		}

		fmt.Printf("%s, addr: %s, pri: %s, tag: %-10s, content: %s\n", timestamp, addr.String(), pri, tag, content)
	}
}
