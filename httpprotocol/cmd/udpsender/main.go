package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error while resolving udp addr, %v", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error while establishing udp connection, %v", err)
	}
	defer conn.Close()

	stdoutReader := bufio.NewReader(os.Stdout)
	for {
		fmt.Print("> ")
		message, err := stdoutReader.ReadString('\n')
		if err != nil {
			log.Printf("error ocurred reading the message %v\n", err)
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("error ocurred sending the message %v\n", err)
		}
	}

}
