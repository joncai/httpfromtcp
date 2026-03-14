// UDP sender that sends a message to a UDP server
// 1. The script will need to be able to resolve the address localhost:42069
// 2. prepare a UDP connection, as well as defer closing it
// 3. print ">" to the console when the script is ready to send a message
// 4. read a line from the standard input, and log error if any
// 5. write the message to the UDP connection, and log error if any
// 6. repeat from step 3
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
			}
			return
		}
		msg = strings.TrimSpace(msg)
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Printf("write error: %v", err)
			continue
		}
	}
}