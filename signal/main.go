// sig server

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listenAdd := os.Args[1]

	s, err := net.Listen("tcp", listenAdd)
	if err != nil {
		fmt.Println("error while creating socket")
		panic(err)
	}

	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println("Error while accept connection.")
		}

		go func(conn net.Conn) {
			handleSignals(conn)
		}(c)
	}
}

type connInfo struct {
	room string
	uuid string
	conn net.Conn
	wkn  bool
}

var connsInfo []connInfo

func handleSignals(conn net.Conn) {
	room := ""
	uuid := ""
	reader := bufio.NewReader(conn)

stop:
	for {
		rl, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while read connection.")
			panic(err)
		}

		parts := strings.Split(rl[:len(rl)-2], "/")
		fmt.Println(parts)
		if len(parts) < 1 {
			fmt.Println("Invalid command")
			continue
		}

		switch parts[0] {
		case "register":
			// register/uuid/room/
			uuid = parts[1]
			room = parts[2]

			ci := connInfo{
				room: room,
				uuid: uuid,
				conn: conn,
				wkn:  false,
			}
			connsInfo = append(connsInfo, ci)
		case "sig":
			// sig/addr/

			// brod cast thire address to other
			for _, c := range connsInfo {
				if c.room == room && c.uuid != uuid {
					c.conn.Write([]byte(rl))
				}
			}
		case "akg":
			// akg/sig/

			// set first of all true
			for i, c := range connsInfo {
				if c.room == room && c.uuid == uuid {
					fmt.Println(c)
					connsInfo[i].wkn = true
				}
			}

			// search all have true
			all := true
			for _, c := range connsInfo {
				if c.room == room {
					if !c.wkn {
						all = false
						break
					}
				}
			}

			// broadcast wkg to all
			if all {
				for _, ci := range connsInfo {
					ci.conn.Write([]byte("sig/start\n"))
				}
			}

		case "stop":
			// stop/sig/
			conn.Close()
			break stop

		default:
			fmt.Println("Invalid command options")
		}
	}
}
