package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const VERSION string = "1.7"

func client(input chan string, ipaddr string) {
	conn, err := net.Dial("tcp", ipaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	nin := bufio.NewScanner(bufio.NewReader(conn))
	nin.Split(bufio.ScanWords)
	nin.Scan()
	if nin.Text() != "HELLO" {
		fmt.Println("Telephone handshake failed")
		fmt.Fprintf(conn, "QUIT\n")
	}
	nin.Scan()

	if nin.Text() != VERSION {
		fmt.Println("Client version incompatable: " + VERSION + " server version: " + nin.Text())
		fmt.Fprintf(conn, "QUIT\n")
	}

}

func server() {

}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: " + os.Args[0] + " <originator> <source> <dest>")
		return
	}
	// source := os.Args[2]
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
