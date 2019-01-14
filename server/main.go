package main

import (
	"fmt"
	"net"
	"os"
	"socket/tcp/server/protocol"
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
	go recvMessageHandle()

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

		n, err := conn.Read(readPacket.Data)
		if err != nil {
			return
		}
		server.RecvMessage[id] <- readPacket.UnPack()
		fmt.Fprintf(os.Stdout, "read size : %d", n)

	}
}

func recvMessageHandle() {
EXITRECV:
	for {
		for i, v := range server.RecvMessage {
			select {
			case message := <-v:
				if message.Msg == "&&EXIT&&" {
					break EXITRECV
				} else {
					fmt.Printf("client %d message \n: %s", i, message.Msg)
				}
			}
		}
	}
}
func sendMessageHandle(idx int) bool {
	// EXIT:
	// 	for {
	// 		for _, v := range server.Clients {
	// 			select {
	// 			case message := <-v.RecvPacket:
	// 				if message.Msg == "&&EXIT&&" {
	// 					break EXIT
	// 				} else {
	// 					sendPack := message.Pack()
	// 					_, err := v.Conn.Write(sendPack.Data[:])
	// 					if err != nil {

	// 					}
	// 				}
	// 			}

	// 		}
	// 	}

EXIT:
	for {
		select {
		case message := <-server.Clients[idx].RecvPacket:
			if message.Msg == "&&EXIT&&" {
				break EXIT
			} else {
				// w := bufio.NewWriter(server.Clients[idx].Conn)
				// r := bufio.NewReader(conn)
				// n, err := r.Read(buf)
				sendPack := message.Pack()
				// w.Write(sendPack.Data)
				// w.Flush()
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

		// for _, v := range server.Clients {
		// 	Message := protocol.Message{}
		// 	Message.Msg = chat
		// 	v.RecvPacket <- Message
		// }
		Message := protocol.Message{}
		Message.Msg = chat
		server.Clients[idx].RecvPacket <- Message
		time.Sleep(time.Second * 10)
	}
}
