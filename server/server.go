package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	UploadPath = "../uploads/"
	BuffSize   = 1024
)

var (
	connPool map[net.Conn]bool
)

func main() {
	var connType = "tcp"
	var connHost = flag.String("host", "localhost", "Hostname or IP address. By default is localhost.")
	var connPort = flag.String("port", "9999", "Server`s port. By default is 9999.")

	flag.Parse()

	connPool = make(map[net.Conn]bool)

	ln, err := net.Listen(connType, *connHost+":"+*connPort)
	if err != nil {
		fmt.Println("Seems like port :" + *connPort + " is already opened")
		os.Exit(1)
	}
	fmt.Println("Server is running...")
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Connection refused: " + conn.RemoteAddr().String())
			continue
		}
		// Store every new connection
		connPool[conn] = true
		fmt.Printf("Now is %d guests in room\n", len(connPool))

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	for {
		cType := make([]byte, 1)
		_, err := conn.Read(cType)
		if err == io.EOF {
			delete(connPool, conn)
			fmt.Println("Somebody has been disconnected.")
			conn.Close()
			return
		}

		switch int(cType[0]) {
		case 1:
			fNameLen := make([]byte, 1)
			fSize := make([]byte, 4)
			conn.Read(fNameLen)
			conn.Read(fSize)

			fName := make([]byte, int(fNameLen[0]))
			conn.Read(fName)
			//fmt.Println(fName, fSize)
			saveFile(string(fName), binary.LittleEndian.Uint32(fSize), conn)
		default:
		case 0:
			getMsg(conn)
		}
	}
}

func saveFile(fName string, fSize uint32, conn net.Conn) {
	fmt.Println(fName, fSize)

	var bytePos uint32
	fileBuff := make([]byte, BuffSize)

	fPath := UploadPath + strings.TrimSpace(fName)

	file, err := os.Create(fPath)
	defer file.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		//fmt.Println("End1.")
		n, err := conn.Read(fileBuff)
		fmt.Printf("Read %d bytes\n", n)
		if err != nil && err != io.EOF {
			fmt.Println("error occured:" + err.Error())
		}
		//wb, err := file.WriteAt(fileBuff, int64(bytePos))
		wb, err := file.Write(fileBuff[:n])
		fmt.Printf("%d bytes are writed.\n", wb)
		bytePos += uint32(wb)
		if bytePos == fSize {
			fmt.Println("File with name '" + fName + "' successfully uploaded")
			conn.Write([]byte("Your file has been upload successfully"))
			break
		}
	}
}

func getMsg(conn net.Conn) {
	buffer := make([]byte, BuffSize)
	conn.Read(buffer)
	// Print message on server
	fmt.Println(string(buffer))
	// Send message for every connection
	for c := range connPool {
		if c != conn {
			c.Write(buffer)
		}
	}

}
