package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	UPLOAD_PATH = "./uploads/"
	BUFF_SIZE   = 1024
)

var (
	connPool map[net.Conn]bool
)

func main() {
	var connHost = flag.String("host", "localhost", "Hostname or IP address. By default is localhost.")
	var connPort = flag.String("port", "9999", "Server`s port. By default is 9999.")
	var connType = flag.String("type", "tcp", "You can use TCP, UDP and IP networks. By default is TCP.")
	flag.Parse()
	connPool = make(map[net.Conn]bool)
	fmt.Println("Starting server...")

	ln, err := net.Listen(*connType, *connHost+":"+*connPort)
	if err != nil {
		fmt.Println("Seems like port :" + *connPort + " is already opened")
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Connection refused: " + conn.RemoteAddr().String())
			continue
		}
		// Store every new connection
		connPool[conn] = true
		// connPool = append(connPool, conn)
		fmt.Printf("Now is %d guests in room\n", len(connPool))
		// fmt.Print(connPool)

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, BUFF_SIZE)
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				delete(connPool, conn)
				fmt.Println("Somebody has been disconnected.")
				conn.Close()
				return
			}
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// switch 0x0 {
		switch buffer[0] {
		case 0x1:
			fnameEndPos := int(buffer[1]) + 2
			fname := string(buffer[2:fnameEndPos])
			fmt.Println("fname: " + fname)
			saveFile(conn, fname, fnameEndPos)
		default:
		case 0x0:
			getMsg(conn, buffer)
		}
	}
}

func saveFile(conn net.Conn, fname string, contentPos int) {
	var bytePos = int64(contentPos) + 1
	fileBuff := make([]byte, BUFF_SIZE)

	file, err := os.Create(UPLOAD_PATH + strings.TrimSpace(fname))
	defer file.Close()

	// fmt.Println(fname)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		_, err := conn.Read(fileBuff)
		if err != nil {
			fmt.Println("error occured:" + err.Error())
		}
		// cleanedFileBuffer := bytes.Trim(fileBuff, "\x00")
		wb, err := file.WriteAt(fileBuff, bytePos)
		bytePos += int64(wb)
		if err == io.EOF {
			break
		}
	}
}

func getMsg(conn net.Conn, buffer []byte) {
	// Print message on server
	fmt.Println(string(buffer))
	// Send message for every connection
	for c := range connPool {
		if c != conn {
			c.Write([]byte(buffer))
		}
	}

}
