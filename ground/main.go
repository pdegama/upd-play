package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var sigAddr string = "46.62.235.81:9090"

func main() {
	room := os.Args[1]
	go findPeer("local/ip4", room+"_local_ip4")
	// go findPeer("local/ip6", room+"_local_ip6")
	for {
	}
}

func findPeer(peerType string, room string) {
	uuid := time.Now().Unix()

	// connect with sig server
	sig, err := net.Dial("tcp", sigAddr)
	if err != nil {
		fmt.Println("Error while connecting signal server!!!")
		panic(err)
	}

	peerAddr, peerNetwork := getLocalPeerAddr(peerType)
	fmt.Println(peerAddr)
	peer, err := net.ListenUDP(peerNetwork, peerAddr)
	if err != nil {
		fmt.Println("Error while create peer socket")
		return
	}

	reg := fmt.Sprintf("register/%d/%s/\n", uuid, room)
	_, err = sig.Write([]byte(reg))
	if err != nil {
		fmt.Println("Error while register ground")
		return
	}
	time.Sleep(5 * time.Second) // wait for other binding

	go findAndSendAddrToSig(sig, peer, peerType)

	var otherPeer *net.UDPAddr

	sigRead := bufio.NewReader(sig)
	for {
		sigLn, err := sigRead.ReadString('\n')
		if err != nil {
			fmt.Println("Error whilte read sig line")
			return
		}
		parts := strings.Split(sigLn, "/")
		if len(parts) < 1 {
			fmt.Println("Invalid command")
			continue
		}

	break_loop1:
		switch parts[0] {
		case "sig":
			if parts[1] == "start" {
				sig.Write([]byte("stop/sig/\n"))
				sig.Close()

				// punch
				punch(peer, otherPeer, room)

				break break_loop1
			} else {
				a, err := net.ResolveUDPAddr(peerNetwork, parts[1])
				if err != nil {
					fmt.Println("Invalid udp address", parts[1])
					continue
				}
				otherPeer = a

				sig.Write([]byte("akg/sig/\n"))
			}
		default:
			fmt.Println("Invalid command")
		}
	}
}

func findAndSendAddrToSig(sig net.Conn, peer *net.UDPConn, peerType string) {
	switch peerType {
	case "local/ip4":
		fmt.Println("peer loacl addr", peer.LocalAddr().String())
		sigStr := fmt.Sprintf("sig/%s/\n", peer.LocalAddr().String())
		sig.Write([]byte(sigStr))
	case "remote/ip4":
	case "remote/ip6":
	}
}

func getLocalPeerAddr(peerType string) (*net.UDPAddr, string) {
	switch peerType {
	case "local/ip4":
		ad := getLocalIp() + ":0"
		a, e := net.ResolveUDPAddr("udp4", ad)
		if e != nil {

			fmt.Println("Error while resole addr")
			panic(e)
		}
		return a, "udp4"
	case "remote/ip4":
		a, e := net.ResolveUDPAddr("udp4", ":0")
		if e != nil {

			fmt.Println("Error while resole addr")
			panic(e)
		}
		return a, "udp4"
	case "remote/ip6":
		a, e := net.ResolveUDPAddr("udp", ":0")
		if e != nil {
			fmt.Println("Error while resole addr")
			panic(e)
		}
		return a, "udp6"
	default:
		panic("invalid peer type")
	}
}

func punch(peer *net.UDPConn, otherPeerAddr *net.UDPAddr, room string) {
	fmt.Println(peer.LocalAddr(), otherPeerAddr)
	go func(peer *net.UDPConn) {
		for {
			peer.WriteToUDP([]byte("--PUNCH--: punch...ping..."), otherPeerAddr)
			time.Sleep(400 * time.Millisecond)
		}
	}(peer)

	r := make([]byte, 1024)

	go func() {
		//	punchd := false
		for {
			rl, add, err := peer.ReadFromUDP(r)
			if err != nil {
				continue
			}
			//	if otherPeerAddr.IP.Equal(add.IP) && otherPeerAddr.Port == add.Port && !punchd {
			//	fmt.Println(room, "Punch Successfull")
			//punchd = false
			//}
			//if !strings.HasPrefix(string(r[:rl]), "--PUNCH--:") {
			fmt.Printf("%s: read \"%s\" from: %s\n", room, r[:rl], add.String())
			//}
		}
	}()

	var input string
	for {
		// fmt.Printf("Write to %s:", otherPeerAddr.String())
		fmt.Scan(&input)
		peer.WriteToUDP([]byte(input), otherPeerAddr)
	}
}

func getLocalIp() string {
	// Connect to an external address (the connection is never actually made)
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	// Get the local address used for this simulated connection
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
