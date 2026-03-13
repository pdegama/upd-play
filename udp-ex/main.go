// UDP exchange

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	arg := os.Args[1]
	listenAddr, err := net.ResolveUDPAddr("udp", arg)
	if err != nil {
		fmt.Println("Invalid Address:", arg)
		panic(err)
	}

	exServer, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("Error to listen UDP server to: ", arg)
		panic(err)
	}

	go readFromUDP(exServer)
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
		parseMsg(string(readData[:readLen]), clientAddr)
		pingMsgFromEx := "PING FROM EX SERVER"
		conn.WriteToUDP([]byte(pingMsgFromEx), clientAddr)
	}
}

func parseMsg(msg string, addr *net.UDPAddr) {
	fmt.Println("Msg resive: ", msg, "From:", addr)
}

type Connections struct {
	userType string
	addr     *net.UDPAddr
	room     string
	active   bool
}

type ConnectionsHandlers struct {
	conns []Connections
}

func (connHandler *ConnectionsHandlers) add(c Connections) {
	connHandler.conns = append(connHandler.conns, c)
}

func (connHandlers *ConnectionsHandlers) get(room string) []Connections {
	var filterConn []Connections
	for i := range connHandlers.conns {
		if connHandlers.conns[i].room == room {
			filterConn = append(filterConn, connHandlers.conns[i])
		}
	}
	return filterConn
}
