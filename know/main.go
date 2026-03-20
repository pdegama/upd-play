package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	go listenKnowServer("udp4", ":5601")
	go listenKnowServer("udp6", ":5602")
	for {
	}
}

func listenKnowServer(network string, address string) {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		fmt.Printf("Error while generate address")
		panic(err)
	}
	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		fmt.Printf("Error while listen on address")
		panic(err)
	}

	rd := make([]byte, 512)
	for {
		rl, add, err := conn.ReadFromUDP(rd)
		if err != nil {
			fmt.Println("Error while read dgram.")
			continue
		}

		parts := strings.Split(string(rd[:rl]), "/")
		switch parts[0] {
		case "want":
			addStr := fmt.Sprintf("addr/%s/\n", addr.String())
			conn.WriteToUDP([]byte(addStr), add)
		default:
		}

		fmt.Println(add.String(), ":", string(rd[:rl]))
	}
}
