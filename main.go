package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	host string
)

func main() {
	//commands := []string{
	//	"ls",
	//	"put",
	//	"get",
	//	"cd",
	//}
	flag.StringVar(&host, "host", "127.0.0.1:4000", "Listening host")
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting goftp at", time.Now().Format("15:04:05\n"))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	_, err := io.WriteString(conn, "Connected to goftp\n")
	if err != nil {
		return
	}
	for {
		buffer := make([]byte, 512)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Printf("%s", buffer)
	}
}
