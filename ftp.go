package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type FTPServer struct {
	Commands map[string]bool
	cmd      *exec.Cmd
}

func (f *FTPServer) list(filename string) {
	f.cmd = exec.Command("ls", "-l")
}

func (f *FTPServer) cwd(dirname string) {
	f.cmd = exec.Command("echo", "changed dir")
	os.Chdir(dirname)
}

func (f *FTPServer) pwd() {
	f.cmd = exec.Command("bash", "-c", "pwd")
}

func (f *FTPServer) quit() string {
	return "quit"
}

func (f *FTPServer) mkd(dirname string) {
	f.cmd = exec.Command("mkdir", dirname)
}

func (f *FTPServer) dele(filename string) {
	f.cmd = exec.Command("rm", filename)
}

func (f *FTPServer) cdup() {
	f.cmd = exec.Command("echo", "changed dir")
	os.Chdir("..")
}

func (f *FTPServer) HandleConn(conn net.Conn) {
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
		out, err := f.HandleCommand(command, args)
		if err != nil {
			log.Fatal(err)
		}

		if out == "quit" {
			conn.Close()
		}

		io.WriteString(conn, out)
	}
}

func (f *FTPServer) HandleCommand(command string, args []string) (string, error) {

	switch command {
	case "LIST":
		f.list(args[0])
	case "CWD":
		f.cwd(args[0])
	case "PWD":
		f.pwd()
	case "QUIT":
		return f.quit(), nil
	case "MKD":
		f.mkd(args[0])
	case "DELE":
		f.dele(args[0])
	case "CDUP":
		f.cdup()
	default:
		return fmt.Sprintf("Invalid command: %s\n", command), nil
	}

	out, err := f.cmd.Output()

	return string(out), err
}

func getCommandAndArgs(buffer []byte) (command string, args []string) {
	message := strings.Fields(string(buffer))
	command = strings.TrimSpace(message[0])
	args = message[1:len(message)]
	return command, args
}
