package main

import (
	"fmt"
	"net"
	"os"
)

func client(input chan string, ipaddr string) {
	conn, err := net.Dial("tcp", ipaddr)
	if err != nil {
		panic(err)
	}

}

func server() {

}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: " + os.Args[0] + " <originator> <source> <dest>")
		return
	}
	source := os.Args[2]
	dest := os.Args[3]

	comm := make(chan string)

	if os.Args[1] == "1" {
		go client(comm, dest)
		server()
	} else {
		go server()
		client(comm, dest)
	}
}
