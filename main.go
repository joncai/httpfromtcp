package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	currentLine := ""
	for {
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			return
		}
		if n == 0 {
			break
		}
		currentString := string(buffer[:n])
		if strings.Contains(currentString, "\n") {
			parts := strings.Split(currentString, "\n")
			currentLine += parts[0]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
			currentString = parts[1]
		}
		currentLine += currentString
	}
	if currentLine != "" {
		fmt.Printf("read: %s\n", currentLine)
	}
}
