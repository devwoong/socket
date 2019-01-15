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
	go messageRead(id)
	go recvMessageHandle(id)

CLIENTEND:
	for {
		select {
		case <-server.Clients[id].Exit:
			fmt.Printf("client %d exit \n", id)
			server.DeleteClient(id)
			break CLIENTEND
		}
	}
}

func messageRead(id int) {
	for {
		readPacket := protocol.Packet{}

		r := bufio.NewReader(server.Clients[id].Conn)
		readPacket.Data = make([]byte, 1024)
		n, err := r.Read(readPacket.Data)
		//n, err := conn.Read(readPacket.Data)
		if err != nil {
			return
		}
		if n != 0 {
			message := readPacket.UnPack(n)
			server.Clients[id].SendPacket <- message
			fmt.Printf("read size : %d \n", n)
		}

	}
}

func recvMessageHandle(idx int) {
EXITRECV:
	for {
		select {
		case message := <-server.Clients[idx].SendPacket:
			switch message.Msg {
			case "&&EXIT&&":
				break EXITRECV
			case "&&ISALIVE&&":
				{
					Message := protocol.Message{}
					Message.Msg = "&&ALIVE&&"
					server.Clients[idx].RecvPacket <- Message
				}
			case "&&ALIVE&&":
				{
					server.Clients[idx].IsAlive = true
				}
			default:
				fmt.Printf("client %d message : %s \n", idx, message.Msg)
			}
		}
	}
}
func sendMessageHandle(idx int) bool {
EXIT:
	for {
		select {
		case message := <-server.Clients[idx].RecvPacket:
			switch message.Msg {
			case "&&EXIT&&":
				break EXIT
			default:
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
		// chat := "enter send client message  \n"
		// Message := protocol.Message{}
		// Message.Msg = chat
		chat := "&&ISALIVE&&"
		Message := protocol.Message{}
		Message.Msg = chat
		server.Clients[idx].RecvPacket <- Message
		time.Sleep(time.Second * 10)
		if server.Clients[idx].IsAlive == false {
			server.Clients[idx].Exit <- true
		} else {
			server.Clients[idx].IsAlive = false
		}
	}
}
