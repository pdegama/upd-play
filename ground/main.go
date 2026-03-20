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
	findPeer("local/ip4", room+"_local_ip4")
}

func findPeer(peerType string, room string) {
	uuid := time.Now().Unix()

	// connect with sig server
	sig, err := net.Dial("tcp", sigAddr)
	if err != nil {
		fmt.Println("Error while connecting signal server!!!")
		panic(err)
	}

	peerAddr := getLocalPeerAddr(peerType)
	peer, err := net.ListenUDP("udp", peerAddr)
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
				punch(peer, otherPeer)

				break break_loop1
			} else {
				a, err := net.ResolveUDPAddr("udp", parts[1])
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
	case "local/ip6":
		fmt.Println("peer loacl addr", peer.LocalAddr().String())
		sigStr := fmt.Sprintf("sig/%s/\n", peer.LocalAddr().String())
		sig.Write([]byte(sigStr))
	case "remote/ip4":
	case "remote/ip6":
	}
}

func getLocalPeerAddr(peerType string) *net.UDPAddr {
	switch peerType {
	case "local/ip4", "remote/ip4":
		a, e := net.ResolveUDPAddr("udp", "0.0.0.0:0")
		if e != nil {
			fmt.Println("Error while resole addr")
			panic(e)
		}
		return a
	case "local/ip6", "remote/ip6":
		a, e := net.ResolveUDPAddr("udp", "[::]:0")
		if e != nil {
			fmt.Println("Error while resole addr")
			panic(e)
		}
		return a
	default:
		panic("invalid peer type")
	}
}

func punch(peer *net.UDPConn, otherPeerAddr *net.UDPAddr) {
	fmt.Println("Punch starting...")
}
