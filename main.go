package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	for {
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			return
		}
		if n == 0 {
			break
		}
		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
