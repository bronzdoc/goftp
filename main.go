package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
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
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	_, err := io.WriteString(conn, "Connected to goftp\n")
	if err != nil {
		return
	}

	for {
		buffer := make([]byte, 256)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Print(err)
			return
		}

		command, args := getCommandAndArgs(buffer)
		out, err := handleCommand(command, args)
		if err != nil {
			log.Fatal(err)
		}

		if out == "quit" {
			conn.Close()
		}

		io.WriteString(conn, out)
	}
}

func getCommandAndArgs(buffer []byte) (command string, args []string) {
	message := strings.Fields(string(buffer))
	command = strings.TrimSpace(message[0])
	args = message[1:len(message)]
	return command, args
}

func handleCommand(command string, args []string) (string, error) {
	var cmd *exec.Cmd

	switch command {
	case "ls":
		cmd = exec.Command(command, "-l")
	case "cd":
		cmd = exec.Command("echo", "changed dir")
		os.Chdir(args[0])
	case "pwd":
		cmd = exec.Command("bash", "-c", "pwd")
	case "quit", "bye":
		return "quit", nil
	default:
		return "", errors.New(fmt.Sprintf("Invalid command: %s", command))
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
