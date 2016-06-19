package main

import (
	"flag"
	"io"
	"log"
	"net"
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
	flag.StringVar(&host, "host", "127.0.0.1:20", "Listening host")
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

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
	_, err := io.WriteString(conn, "Connected to goftp")
	if err != nil {
		return
	}
}
