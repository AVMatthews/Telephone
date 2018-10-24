package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//#include <stddef.h>
//#include <stdint.h>
//uint16_t checksum(void *data, size_t size) {
// uint32_t sum = 0;
// uint16_t *data16 = data;
// while(size > 0) {
// sum += *data16++;
// size -= 2;
// }
// if(size > 0) sum += *((uint8_t *) data16);
// while(sum >> 16) sum = (sum & 0xFFFF) + (sum >> 16);
// return ~sum;
// }
import "C"

const VERSION string = "1.7"

func checksum(in string) uint16 {
	var tmp C.size_t = C.size_t(len(in) + 1)
	return uint16(C.checksum(unsafe.Pointer(C.CString(in)), tmp))
}

func readHeaders(in string) map[string]string {
	m := make(map[string]string)
	if in == "" {
		m["Hop"] = "0"
		m["MessageId"] = "0"
	} else {
		sin := bufio.NewScanner(strings.NewReader(in))
		sin.Split(bufio.ScanWords)
		sin.Scan()
		sin.Scan()
		m["Hop"] = sin.Text()
		sin.Scan()
		sin.Scan()
		m["MessageId"] = sin.Text()
	}

	return m
}

func extractMess(in string) string {
	for i := 0; i+4 < len(in); i++ {
		if in[i] == '\r' && in[i+1] == '\n' && in[i+2] == '\r' && in[i+3] == '\n' {
			for j := i; j+4 < len(in); j++ {
				if in[j] == '\r' && in[j+1] == '\n' && in[j+2] == '.' && in[j+3] == '\r' && in[j+4] == '\n' {
					return in[i+4 : j]
				}
			}
		}
	}
	return ""
}

func addHeaders(input string, srcIp string, destIp string) string {
	var out string
	m := readHeaders(input)
	tmp, _ := strconv.Atoi(m["Hop"])
	hopNum := strconv.Itoa(tmp + 1)
	t := time.Now()
	var chs string
	chs = fmt.Sprintf("%.4X", checksum(extractMess(input)))
	out += "Hop: " + hopNum + "\r\n"
	out += "MessageId: " + m["MessageId"] + "\r\n"
	out += "FromHost: " + srcIp + "\r\n"
	out += "ToHost: " + destIp + "\r\n"
	out += "System: " + runtime.GOOS + "\r\n"
	out += "Program: Golang/Go\r\n"
	out += "Author: Hunter Bashaw/Abigail Matthews\r\n"
	out += "SendingTimestamp: " + strings.Replace(t.Format("15:04:05.000"), ".", ":", -1) + "\r\n"
	out += "MessageChecksum: " + chs + "\r\n"
	out += "HeaderChecksum: " + "\r\n"
	//Warning
	//Transform
	out += input
	return out
}

func client(input chan string, source string, ipaddr string) {
	conn, err := net.Dial("tcp", ipaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	nin := bufio.NewScanner(bufio.NewReader(conn))
	nin.Split(bufio.ScanWords)
	nin.Scan()
	//Check for HELLO
	if nin.Text() != "HELLO" {
		fmt.Println("Telephone handshake failed")
		fmt.Fprintf(conn, "QUIT\r\n")
		nin.Scan()
		if nin.Text() == "GOODBYE" {
			return
		}
		return
	}
	nin.Scan()
	//Check version number
	if nin.Text() != VERSION {
		fmt.Println("Client version incompatable: " + VERSION + " server version: " + nin.Text())
		fmt.Fprintf(conn, "QUIT\r\n")
		nin.Scan()
		if nin.Text() == "GOODBYE" {
			return
		}
		return
	}

	fmt.Fprintf(conn, "DATA\r\n")
	final := addHeaders(<-input, source, ipaddr)
	fmt.Fprintf(conn, final)

	//Handle connection forever
	for {

	}

}

func server(output chan string) {
	output <- "Hop: 1\r\nMessageId: 3456\r\nFromHost: 192.168.0.12:9879\r\nToHost: 192.168.0.4:8888\r\nSystem: WINDOWS/XP\r\nProgram: JAVA/JAVAC\r\nAuthor: Frodo Baggins\r\nSendingTimestamp: 17:00:00:000\r\nMessageChecksum: 432F\r\nHeadersChecksum: A350\r\n\r\nHi how are you? I'm good.\r\n.\r\n"

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
		go client(comm, source, dest)
		server(comm)
	} else {
		go server(comm)
		client(comm, source, dest)
	}
}
