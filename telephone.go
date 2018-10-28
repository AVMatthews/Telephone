package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// #include <stddef.h>
// #include <stdint.h>
// uint16_t checksum(void *data, size_t size) {
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
	fmt.Println("|" + in + "|")
	var tmp C.size_t = C.size_t(len(in))
	if len(in) == 0 {
		return 0
	}
	cstring := C.CString(in)
	temp := C.checksum(unsafe.Pointer(cstring), tmp)
	return uint16(temp)
}

func readHeaders(in string) map[string]string {
	m := make(map[string]string)
	if in == "" {
		m["Hop"] = "0"
		m["MessageId"] = "0"
	} else {
		sin := bufio.NewScanner(strings.NewReader(in))
		sin.Split(bufio.ScanWords)
		for {
			if !sin.Scan() {
				return m
			}
			_, in := m["Hop"]
			if sin.Text() == "Hop:" && !in {
				sin.Scan()
				m["Hop"] = sin.Text()
				continue
			}
			_, in = m["MessageId"]
			if sin.Text() == "MessageId:" && !in {
				sin.Scan()
				m["MessageId"] = sin.Text()
				continue
			}
		}
	}

	return m
}

func readMesgChecksum(in string) uint16 {
	sin := bufio.NewScanner(strings.NewReader(in))
	sin.Split(bufio.ScanWords)
	for {
		sin.Scan()
		if sin.Text() == "MessageChecksum:" {
			sin.Scan()
			value, err := strconv.ParseUint(sin.Text(), 16, 16)
			if err != nil {
			}
			return uint16(value)
		}
	}

}

func extractMess(in string) string {
	for i := 0; i+4 < len(in); i++ {
		if in[i] == '\r' && in[i+1] == '\n' && in[i+2] == '\r' && in[i+3] == '\n' {
			for j := i; j+4 < len(in); j++ {
				if in[j] == '\r' && in[j+1] == '\n' && in[j+2] == '.' && in[j+3] == '\r' && in[j+4] == '\n' {
					return in[i+4 : j-1]
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
	out += "Program: Golang/Go1.9,4\r\n"
	out += "Author: Hunter Bashaw/Abigail Matthews\r\n"
	out += "SendingTimestamp: " + strings.Replace(t.Format("15:04:05.000"), ".", ":", -1) + "\r\n"
	out += "MessageChecksum: " + chs + "\r\n"
	//Header checksum
	//Warning
	//Transform
	out += input
	return out
}

func client(input chan string, source string, ipaddr string) {
	//Open connection to server
	conn, err := net.Dial("tcp", ipaddr)
	if err != nil {
		panic(err)
	}
	//Close connection when function ends
	defer conn.Close()
	nin := bufio.NewScanner(bufio.NewReader(conn))
	nin.Split(bufio.ScanWords)
	nin.Scan()
	out := ""
	//Check for HELLO and do handshake
	if nin.Text() != "HELLO" {
		fmt.Println("Telephone handshake failed")
		fmt.Fprintf(conn, "QUIT\r\n")
		nin.Scan()
		if nin.Text() == "GOODBYE" {
			return
		}
		return
	}
	out = out + nin.Text()
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
	} else {
		out = out + " " + nin.Text() + "\r\n"
		fmt.Fprintf(conn, out)
	}

	//Loop forever
	for {
		//Get data to send from the channel connecting the client and server
		data := <-input
		fmt.Println("Client: " + data)
		//If data is SIGTERM then close the connection and end program
		if data == "SIGTERM" {
			fmt.Fprintf(conn, "QUIT\r\n")
			nin.Scan()
			if nin.Text() == "GOODBYE" {
				return
			}
			return
		}
		//Send data we got from the channel
		fmt.Fprintf(conn, "DATA\r\n")
		final := addHeaders(data, source, ipaddr)
		fmt.Fprintf(conn, final)

		//Get either SUCCESS or WARN on SUCCESS do nothing on WARN print error and move on
		nin.Scan()
		if nin.Text() == "WARN" {
			fmt.Println("Server warned of checksum failure")
		}
	}

}

func server(output chan string, source string) {
	ip_port := strings.SplitAfter(source, ":")
	//Start listening on Port
	c, err := net.Listen("tcp", "localhost:"+ip_port[1])
	if err != nil {
		panic(err)
	}
	//Close connection at end of function
	defer c.Close()
	//Accept Conenction
	conn, err := c.Accept()
	//Send HELLO <version#>
	conn.Write([]byte("HELLO " + VERSION + "\r\n"))
	//get response from client
	nin := bufio.NewScanner(bufio.NewReader(conn))
	//nin.Split(bufio.ScanWords)
	nin.Scan()
	if nin.Text() != "HELLO 1.7" { //change to check global version number
		fmt.Println("Warning incorrect handshake response from Client\n")
		fmt.Println(nin.Text())
	}
	for {
		//get DATA
		nin = bufio.NewScanner(bufio.NewReader(conn))
		nin.Split(bufio.ScanLines)
		nin.Scan()
		mesg := ""
		nextLine := ""
		if nin.Text() == "DATA" {
			for {
				nin.Scan()
				nextLine = nin.Text()
				mesg += (nextLine + "\r\n")
				if nextLine == "." {
					break
				}
			}
		} else if nin.Text() == "QUIT" {
			conn.Write([]byte("GOODBYE\r\n"))
			return
		}
		// fmt.Println("mesg: " + mesg)
		//fmt.Println("scanned")
		//mesg := extractMess(nin.Text())
		nChecksum := checksum(extractMess(mesg))
		fmt.Println("Made checksum")
		oChecksum := readMesgChecksum(mesg)
		fmt.Println("read Checksum")

		if nChecksum == oChecksum {
			conn.Write([]byte("SUCCESS\r\n"))
			output <- mesg
		} else {
			conn.Write([]byte("WARN\r\n"))
			output <- mesg
		}
	}
	/*nin.SplitAfter("\n")
	mesg = ""
	nextLine = ""
	for {
		nextLine = in.Scan()
		if nextLine == "\r\n"{
			break
		}
	}
	for {
		nextLine = nin.ScanLines();
		if mesg == "."{

		} else {
			mesg =
		}
	}*/
	//for { //Maybe make this loop surround aroudn the handshking too?
	//message, err := bufio.NewReader(conn).ReadString("\r\n.\r\n")
	//}
	//This line is for testing the client
	//output <- "Hop: 1\r\nMessageId: 3456\r\nFromHost: 192.168.0.12:9879\r\nToHost: 192.168.0.4:8888\r\nSystem: WINDOWS/XP\r\nProgram: JAVA/JAVAC\r\nAuthor: Frodo Baggins\r\nSendingTimestamp: 17:00:00:000\r\nMessageChecksum: 432F\r\nHeadersChecksum: A350\r\nHop: 0\r\nMessageId: 3456\r\nFromHost: 192.168.0.1:34953\r\nToHost: 192.168.0.12:8888\r\nSystem: LINIX/DEBIAN/R3.0\r\nProgram: C++/GCC\r\nAuthor: Alex, J./Jacky Elton/David Wang\r\nSendingTimestamp: 16:59:59:009\r\nMessageChecksum: 423F\r\nHeadersChecksum: 6F38\r\n\r\nHi how are you? I'm good.\r\n.\r\n"
	// output <- "SIGTERM"
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: " + os.Args[0] + " <originator> <source> <dest>")
		return
	}
	source := os.Args[2]
	dest := os.Args[3]

	comm := make(chan string)
	//Handle CTRL-C to gracefully close the clients connection to server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			comm <- "SIGTERM"
		}
	}()

	if os.Args[1] == "1" {
		go client(comm, source, dest)
		server(comm, source)
	} else {
		go server(comm, source)
		client(comm, source, dest)
	}
}
