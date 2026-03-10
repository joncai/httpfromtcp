package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}
