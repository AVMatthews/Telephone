package main

import (
	"fmt"
	"os"
)

func client() {

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

	if os.Args[1] == "1" {
		go client()
		server()
	} else {
		go server()
		client()
	}
}
