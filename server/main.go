package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"socket/server/protocol"
	"time"
)

var server protocol.Server

func main() {
	service := "0.0.0.0:1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal errors  : %s", err.Error())
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, "host Ip : %s \t host port : %s", tcpAddr.IP, tcpAddr.Port)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal errors : %s", err.Error())
		os.Exit(1)
	}

	server = protocol.Server{}
	server.Initialize()
	connectLoop(listener)
}
func connectLoop(listener *(net.TCPListener)) {
	var i int = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go clientHandle(conn, i)
		i++
	}
}
func clientHandle(conn net.Conn, id int) {
	defer conn.Close()
	client := protocol.Client{}
	client.CreateClient(conn, id)
	server.AddClient(&client)

	go messageChat(id)
	go sendMessageHandle(id)
	go messageRead(conn, id)
	go recvMessageHandle(id)

CLIENTEND:
	for {
		select {
		case <-server.Exit[id]:
			fmt.Printf("client %d exit \n", id)
			server.DeleteClient(id)
			break CLIENTEND
		}
	}
}

func messageRead(conn net.Conn, id int) {
	for {
		readPacket := protocol.Packet{}

		r := bufio.NewReader(conn)
		readPacket.Data = make([]byte, 1024)
		n, err := r.Read(readPacket.Data)
		//n, err := conn.Read(readPacket.Data)
		if err != nil {
			return
		}
		if n != 0 {
			message := readPacket.UnPack(n)
			server.Clients[id].SendPacket <- message
			fmt.Fprintf(os.Stdout, "read size : %d", n)
		}

	}
}

func recvMessageHandle(idx int) {
EXITRECV:
	for {
		select {
		case message := <-server.Clients[idx].SendPacket:
			if message.Msg == "&&EXIT&&" {
				break EXITRECV
			} else {
				fmt.Printf("client %d message \n: %s", idx, message.Msg)
			}
		}
	}
}
func sendMessageHandle(idx int) bool {
EXIT:
	for {
		select {
		case message := <-server.Clients[idx].RecvPacket:
			if message.Msg == "&&EXIT&&" {
				break EXIT
			} else {
				sendPack := message.Pack()
				_, err := server.Clients[idx].Conn.Write(sendPack.Data[:])
				if err != nil {

				}
			}
		}

	}
	return true
}

func messageChat(idx int) {
	for {
		chat := "enter send client message  \n"
		Message := protocol.Message{}
		Message.Msg = chat
		server.Clients[idx].RecvPacket <- Message
		time.Sleep(time.Second * 10)
	}
}
