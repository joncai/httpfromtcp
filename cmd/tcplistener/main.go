package main

import (
	"fmt"
	"net"

	"httpfromtcp/internal/request"
)

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

	r, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Request line:")
	fmt.Println("- Method:", r.RequestLine.Method)
	fmt.Println("- Target:", r.RequestLine.RequestTarget)
	fmt.Println("- Version:", r.RequestLine.HttpVersion)
}
