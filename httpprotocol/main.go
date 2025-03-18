package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("./message.txt")
	if err != nil {
		panic("error opening the message file")
	}
	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChannel := make(chan string)
	go func() {
		data := make([]byte, 8)
		line := ""
		for {
			n, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic("an error occurred while reading the file")
			}
			parts := strings.Split(string(data[:n]), "\n")
			if len(parts) == 1 {
				line += parts[0]
				continue
			}
			for i := range len(parts) - 1 {
				line += parts[i]
			}
			lineChannel <- line
			line = parts[len(parts)-1]
		}
		if len(line) > 0 {
			lineChannel <- line
		}
		f.Close()
		close(lineChannel)
	}()
	return lineChannel
}
