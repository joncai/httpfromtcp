package main

import (
	"bufio"
	"fmt"
	"net"
)

func getLinesChannel(c net.Conn) <-chan string {
	ch := make(chan string)
	go func() {
		defer c.Close()
		defer close(ch)
		scanner := bufio.NewScanner(c)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection accepted")

	for line := range getLinesChannel(conn) {
		fmt.Println(line)
	}
}
