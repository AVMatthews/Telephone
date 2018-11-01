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
// #include <string.h>
// unsigned short checksum(void *data, unsigned long size) {
// unsigned long sum = 0, i = 0;
// unsigned char *udata = data;
// while(i < size) {
// sum += udata[i] << ( (i&1) ? 0 : 8 );
// i++;
// }
// while(sum >> 16) sum = (sum & 0xFFFF) + (sum >> 16);
// return ~((unsigned short) sum);
// }
import "C"

const VERSION string = "1.7.1"

func checksum(in string) uint16 {
	fmt.Println("|" + in + "|")
	// var size C.size_t = C.size_t(len(in))
	if len(in) == 0 {
		return 0
	}
	cstring := C.CString(in)
	temp := C.checksum(unsafe.Pointer(cstring), C.strlen(cstring))
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
				break
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
		sin = bufio.NewScanner(strings.NewReader(in))
		sin.Split(bufio.ScanLines)

		Systems := make([]string, 0)
		Programs := make([]string, 0)
		Authors := make([]string, 0)
		Warnings := make([]string, 0)
		Transforms := make([]string, 0)
		FromHosts := make([]string, 0)
		ToHosts := make([]string, 0)
		SendingTimestamps := make([]string, 0)

		for {
			//fmt.Println("second for")
			if !sin.Scan() {
				_, in := m["System"]

				//Make Slices into Strings to print
				i := 0
				next := 0
				systemString := "The Unique Systems that touched this message are"
				for i < len(Systems) {
					j := 0
					for j < i {
						if Systems[i] == Systems[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						systemString = systemString + " (" + strconv.Itoa(next) + ") " + Systems[i]
					}
					i++
				}
				if !in {
					m["System"] = systemString
					//fmt.Println(systemString)
				}

				_, in = m["Program"]
				i = 0
				next = 0
				programString := "The Unique Programs that touched this message are"
				for i < len(Programs) {
					j := 0
					for j < i {
						if Programs[i] == Programs[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						programString = programString + " (" + strconv.Itoa(next) + ") " + Programs[i]
					}
					i++
				}
				if !in {
					m["Program"] = programString
					//fmt.Println(systemString)
				}

				_, in = m["Authors"]
				i = 0
				next = 0
				authorString := "The Unique Authors that touched this message are"
				for i < len(Authors) {
					j := 0
					for j < i {
						if Authors[i] == Authors[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						authorString = authorString + " (" + strconv.Itoa(next) + ") " + Authors[i]
					}
					i++
				}
				if !in {
					m["Author"] = authorString
					//fmt.Println(systemString)
				}

				_, in = m["Warning"]
				i = 0
				next = 0
				warningString := "The Unique Warnings that touched this message are"
				for i < len(Warnings) {
					j := 0
					for j < i {
						if Warnings[i] == Warnings[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						warningString = warningString + " (" + strconv.Itoa(next) + ") " + Warnings[i]
					}
					i++
				}
				if !in {
					m["Warning"] = warningString
					//fmt.Println(systemString)
				}

				_, in = m["Transform"]
				i = 0
				next = 0
				transformString := "The Unique Transforms that touched this message are"
				for i < len(Transforms) {
					j := 0
					for j < i {
						if Transforms[i] == Transforms[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						transformString = transformString + " (" + strconv.Itoa(next) + ") " + Transforms[i]
					}
					i++
				}
				if !in {
					m["Transform"] = transformString
					//fmt.Println(systemString)
				}

				_, in = m["FromHost"]
				i = 0
				next = 0
				fromHostString := "The Unique FromHosts that touched this message are"
				for i < len(FromHosts) {
					j := 0
					for j < i {
						if FromHosts[i] == FromHosts[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						fromHostString = fromHostString + " (" + strconv.Itoa(next) + ") " + FromHosts[i]
					}
					i++
				}
				if !in {
					m["FromHost"] = fromHostString
					//fmt.Println(systemString)
				}

				_, in = m["ToHost"]
				i = 0
				next = 0
				toHostString := "The Unique ToHosts that touched this message are"
				for i < len(ToHosts) {
					j := 0
					for j < i {
						if ToHosts[i] == ToHosts[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						toHostString = toHostString + " (" + strconv.Itoa(next) + ") " + ToHosts[i]
					}
					i++
				}
				if !in {
					m["ToHost"] = toHostString
					//fmt.Println(systemString)
				}

				_, in = m["SendingTimestamp"]
				i = 0
				next = 0
				sendingTimestampString := "The Unique SendingTimestamps that touched this message are"
				for i < len(SendingTimestamps) {
					j := 0
					for j < i {
						if SendingTimestamps[i] == SendingTimestamps[j] {
							break
						}
						j++
					}
					if i == j {
						next++
						sendingTimestampString = sendingTimestampString + " (" + strconv.Itoa(next) + ") " + SendingTimestamps[i]
					}
					i++
				}
				if !in {
					m["SendingTimestamp"] = sendingTimestampString
					//fmt.Println(systemString)
				}

				return m
			}

			if strings.HasPrefix(sin.Text(), "System:") {
				newSystem := strings.SplitAfter(sin.Text(), ": ")[1]
				Systems = append(Systems, newSystem)
			}
			if strings.HasPrefix(sin.Text(), "Program:") {
				newProgram := strings.SplitAfter(sin.Text(), ": ")[1]
				Programs = append(Programs, newProgram)
			}
			if strings.HasPrefix(sin.Text(), "Author:") {
				newAuthor := strings.SplitAfter(sin.Text(), ": ")[1]
				Authors = append(Authors, newAuthor)
			}
			if strings.HasPrefix(sin.Text(), "Warning:") {
				newWarning := strings.SplitAfter(sin.Text(), ": ")[1]
				Warnings = append(Warnings, newWarning)
			}
			if strings.HasPrefix(sin.Text(), "Transform:") {
				newTransform := strings.SplitAfter(sin.Text(), ": ")[1]
				Transforms = append(Transforms, newTransform)
			}
			if strings.HasPrefix(sin.Text(), "FromHost:") {
				newFromHost := strings.SplitAfter(sin.Text(), ": ")[1]
				FromHosts = append(FromHosts, newFromHost)
			}
			if sin.Text() == "ToHost:" {
				newToHost := strings.SplitAfter(sin.Text(), ": ")[1]
				ToHosts = append(ToHosts, newToHost)
			}
			if sin.Text() == "SendingTimestamp:" {
				newSendingTimestamp := strings.SplitAfter(sin.Text(), ": ")[1]
				SendingTimestamps = append(SendingTimestamps, newSendingTimestamp)
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
	out += "Program: Golang/Go1.9.4\r\n"
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

func server(output chan string, source string, isOriginator string) {
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
	if nin.Text() != "HELLO 1.7.1" { //change to check global version number
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
		nChecksum := checksum(extractMess(mesg))
		fmt.Println("Made checksum")
		oChecksum := readMesgChecksum(mesg)
		fmt.Println("read Checksum")

		if isOriginator == "1" {
			if nChecksum == oChecksum {
				fmt.Println("Checksums are VALID")
				headers := readHeaders(mesg)
				fmt.Println("Number of Hops: " + headers["Hop"])
				fmt.Println("MessageId: " + headers["MessageId"])
				if headers["System"] != "The Unique Systems that touched this message are" {
					fmt.Println(headers["System"])
				}
				if headers["Program"] != "The Unique Programs that touched this message are" {
					fmt.Println(headers["Program"])
				}
				if headers["Author"] != "The Unique Authors that touched this message are" {
					fmt.Println(headers["Author"])
				}
				if headers["FromHost"] != "The Unique FromHosts that touched this message are" {
					fmt.Println(headers["FromHost"])
				}
				if headers["ToHost"] != "The Unique ToHosts that touched this message are" {
					fmt.Println(headers["ToHost"])
				}
				if headers["Warning"] != "The Unique Warnings that touched this message are" {
					fmt.Println(headers["Warning"])
				}
				if headers["Transform"] != "The Unique Transforms that touched this message are" {
					fmt.Println(headers["Transform"])
				}
				if headers["SendingTimestamp"] != "The Unique SendingTimestamps that touched this message are" {
					fmt.Println(headers["SendingTimestamp"])
				}
			} else {
				fmt.Println("Checksums are INVALID")
				headers := readHeaders(mesg)
				fmt.Println("Number of Hops: " + headers["Hop"])
				fmt.Println("MessageId: " + headers["MessageId"])
				if headers["System"] != "The Unique Systems that touched this message are" {
					fmt.Println(headers["System"])
				}
				if headers["Program"] != "The Unique Programs that touched this message are" {
					fmt.Println(headers["Program"])
				}
				if headers["Author"] != "The Unique Authors that touched this message are" {
					fmt.Println(headers["Author"])
				}
				if headers["FromHost"] != "The Unique FromHosts that touched this message are" {
					fmt.Println(headers["FromHost"])
				}
				if headers["ToHost"] != "The Unique ToHosts that touched this message are" {
					fmt.Println(headers["ToHost"])
				}
				if headers["Warning"] != "The Unique Warnings that touched this message are" {
					fmt.Println(headers["Warning"])
				}
				if headers["Transform"] != "The Unique Transforms that touched this message are" {
					fmt.Println(headers["Transform"])
				}
				if headers["SendingTimestamp"] != "The Unique SendingTimestamps that touched this message are" {
					fmt.Println(headers["SendingTimestamp"])
				}
			}

		} else {
			if nChecksum == oChecksum {
				fmt.Println("Checksums are VALID")
			} else {
				fmt.Println("Checksums are INVALID")
			}
		}

		if nChecksum == oChecksum {
			conn.Write([]byte("SUCCESS\r\n"))
			output <- mesg
		} else {
			conn.Write([]byte("WARN\r\n"))
			output <- mesg
		}
	}
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
		server(comm, source, os.Args[1])
	} else {
		go server(comm, source, os.Args[1])
		client(comm, source, dest)
	}
}
