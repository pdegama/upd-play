package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var needSendExPing bool = true

func main() {
	exAddrString := os.Args[1]
	roomKey := os.Args[2]

	exAddr, err := net.ResolveUDPAddr("udp", exAddrString)
	if err != nil {
		fmt.Println("Error to Resolve Ex Address:", exAddrString)
		panic(err)
	}
	myAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		fmt.Println("Error to Resolve My Address.")
		panic(err)
	}

	fmt.Println("Key:", roomKey)
	fmt.Println("ExAddress:", exAddr.String())
	fmt.Println("MyAddress:", myAddr.String())

	myConn, err := net.ListenUDP("udp", myAddr)
	if err != nil {
		fmt.Println("Error while listen targat on:", myAddr)
		panic(err)
	}

	go readFromUDP(myConn)
	go sendExPing(myConn, exAddr)
	for {
	}
}

func readFromUDP(conn *net.UDPConn) {
	readData := make([]byte, 512)
	for {
		readLen, clientAddr, err := conn.ReadFromUDP(readData)
		if err != nil || readLen == 0 {
			fmt.Println("Error while read")
			panic(err)
		}
		fmt.Println("Data Recive:", string(readData[:readLen]), "From:", clientAddr.String())
	}
}

func sendExPing(conn *net.UDPConn, exAddr *net.UDPAddr) {
	for needSendExPing {
		conn.WriteToUDP([]byte("Want Punch\n"), exAddr)
		time.Sleep(200 * time.Millisecond)
	}
}
