package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	BuffSize = 1024
)

func main() {
	var connType = "tcp"
	var connHost = flag.String("host", "localhost", "Hostname or IP address. By default is localhost.")
	var connPort = flag.String("port", "9999", "Server`s port. By default is 9999.")
	flag.Parse()

	conn, err := net.Dial(connType, *connHost+":"+*connPort)
	if err != nil {
		fmt.Println("Seems the server is down.")
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	defer conn.Close()

	go getMsgs(conn)
	handleMsgs(conn)
}

func getMsgs(conn net.Conn) {
	for {
		buffer := make([]byte, BuffSize)
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
	conn.Write([]byte{0})
	conn.Write([]byte(text))
}

func sendFile(conn net.Conn, fName string) {
	cType := []byte{1}
	fNameLen := []byte{byte(len(fName))}
	var bytePos int64
	fileBuff := make([]byte, BuffSize)

	file, err := os.Open(fName)
	if err != nil {
		fmt.Println("File `" + fName + "` not found.")
		return
	}
	defer file.Close()

	conn.Write(cType)
	conn.Write(fNameLen)

	fInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("File meta error.\n%s\n", err.Error())
		return
	}

	fSize := fInfo.Size()
	err = binary.Write(conn, binary.LittleEndian, int32(fSize))
	if err != nil {
		fmt.Printf("Error occured:\n%s\n", err.Error())
		return
	}

	conn.Write([]byte(fName))
	//fmt.Println(cType, fNameLen, fSize, fName)
	for {
		rb, err := file.ReadAt(fileBuff, bytePos)
		//fmt.Println(rb, bytePos)
		bytePos += int64(rb)

		conn.Write(fileBuff[:rb])

		if err == io.EOF {
			fmt.Println("Your file has been sent.")
			break
		}
	}

}
