package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	BUFF_SIZE = 1024
)

func main() {
	var connHost = flag.String("host", "localhost", "Hostname or IP address. By default is localhost.")
	var connPort = flag.String("port", "9999", "Server`s port. By default is 9999.")
	var connType = flag.String("type", "tcp", "You can use TCP, UDP and IP networks. By default is TCP.")
	flag.Parse()

	fmt.Println("Connecting to server ...")
	conn, err := net.Dial(*connType, *connHost+":"+*connPort)
	if err != nil {
		fmt.Println("Seems the server is down.")
		os.Exit(1)
	}

	defer conn.Close()
	go getMsgs(conn)
	handleMsgs(conn)
}

func getMsgs(conn net.Conn) {
	for {
		buffer := make([]byte, BUFF_SIZE)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(string(buffer))
	}
}

func handleMsgs(conn net.Conn) {
	fmt.Println("Lets start the chat!")
	fmt.Println("To upload the file just type 'file <filename>'")

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		arrayOfCommands := strings.Split(strings.TrimSpace(text), " ")
		if len(arrayOfCommands) == 2 && arrayOfCommands[0] == "file" {
			sendFile(conn, arrayOfCommands[1])
		} else {
			sendMsg(conn, text)
		}
	}
}

func sendMsg(conn net.Conn, text string) {
	conn.Write([]byte{0x0})
	conn.Write([]byte(text))
}

func sendFile(conn net.Conn, fname string) {
	var bytePos int64
	fileBuff := make([]byte, BUFF_SIZE)

	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("File `" + fname + "` not found.")
		return
	}

	conn.Write([]byte{byte(1)})
	conn.Write([]byte{byte(len(fname))})
	conn.Write([]byte(fname))
	for {
		rb, err := file.ReadAt(fileBuff, bytePos)
		bytePos += int64(rb)
		if err == io.EOF {
			fmt.Println("Your file has been upload successfully.")
			break
		}
		conn.Write(fileBuff)
	}

	file.Close()
}
