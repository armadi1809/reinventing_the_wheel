package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic("error setting up the listener")
	}
	defer listener.Close()
	// treat connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("error accepting the connection")
		}
		fmt.Println("Connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("closing closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChannel := make(chan string)
	go func() {
		defer f.Close()
		defer close(lineChannel)
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
	}()
	return lineChannel
}
