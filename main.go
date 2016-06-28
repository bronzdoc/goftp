package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	host string
	ftp  FTPServer
)

func main() {
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
		go ftp.HandleConn(conn)
	}
}
