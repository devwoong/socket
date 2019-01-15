package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"socket/client/protocol"
	"time"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println(os.Stderr, "Usage %s host:port", os.Args[0])
	// 	os.Exit(0)
	// }
	// service := os.Args[1]
	service := "localhost:1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	fmt.Printf("start")
	// conn.Write([]byte("sdfdsf"))
	// conn.Write([]byte("\r\n\r\n"))
	go recvMessage(conn)
	go sendMessage(conn)
	for {

	}
}
func recvMessage(conn net.Conn) {
	for {
		buff := make([]byte, 1024)
		n, _ := conn.Read(buff)
		if n > 0 {
			log.Printf("Receive: %s", buff[:n])
		}
	}
}

func sendMessage(conn net.Conn) {

	for {
		message := protocol.Message{}
		message.Msg = "client Mesasge"
		fmt.Printf("write to server message is %s\n", message.Pack().Data[:])
		conn.Write(message.Pack().Data[:])
		time.Sleep(time.Second * 10)
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}
}
