package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"socket/server/protocol"
)

var server protocol.Server

func main() {
	service := "0.0.0.0:1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal errors  : %s", err.Error())
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, "host Ip : %s \t host port : %s\n", tcpAddr.IP, tcpAddr.Port)
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
		//go isAlivceClient()
		i++
	}
}
func clientHandle(conn net.Conn, id int) {
	defer conn.Close()
	client := protocol.Client{}
	client.CreateClient(conn, id)
	server.AddClient(&client)

	fmt.Printf("connect client id : %d, ip : %s \n", id, server.Clients[id].Conn.RemoteAddr())
	fmt.Printf("current client connet num : %d \n", len(server.Clients))

	go sendMessageHandle(id)
	go recvMessageHandle(id)
	messageRead(id)
}

func messageRead(id int) {
	for {
		readPacket := protocol.Packet{}

		r := bufio.NewReader(server.Clients[id].Conn)
		readPacket.Data = make([]byte, 1024)
		n, err := r.Read(readPacket.Data)
		//n, err := conn.Read(readPacket.Data)
		if err == io.EOF {
			exit := "&&EXIT&&"
			exitMessage := protocol.Message{}
			exitMessage.Msg = exit
			server.Clients[id].SendPacket <- exitMessage
			server.Clients[id].RecvPacket <- exitMessage
			fmt.Printf("client exit id : %d, ip : %s \n", id, server.Clients[id].Conn.RemoteAddr())
			server.DeleteClient(id)
			fmt.Printf("current client connet num : %d \n", len(server.Clients))
			return
		}
		if n != 0 {
			message := readPacket.UnPack(n)
			server.Clients[id].SendPacket <- message
			fmt.Printf("read size : %d \n", n)
		}

	}
}

func recvMessageHandle(id int) {
EXITRECV:
	for {
		select {
		case message := <-server.Clients[id].SendPacket:
			switch message.Msg {
			case "&&EXIT&&":
				break EXITRECV
			case "&&ISALIVE&&":
				{
					Message := protocol.Message{}
					Message.Msg = "&&ALIVE&&"
					server.Clients[id].RecvPacket <- Message
				}
			// case "&&ALIVE&&":
			// 	{
			// 		server.Clients[id].IsAlive = true
			// 	}
			default:
				fmt.Printf("client %d message : %s \n", id, message.Msg)
			}
		}
	}
}
func sendMessageHandle(id int) bool {
EXIT:
	for {
		select {
		case message := <-server.Clients[id].RecvPacket:
			switch message.Msg {
			case "&&EXIT&&":
				break EXIT
			default:
				sendPack := message.Pack()
				_, err := server.Clients[id].Conn.Write(sendPack.Data[:])
				if err != nil {

				}
			}
		}
	}
	return true
}

// func isAlivceClient() {
// 	live := "&&ISALIVE&&"
// 	Message := protocol.Message{}
// 	Message.Msg = live
// 	for {
// 		for _, client := range server.Clients {
// 			if client == nil {
// 				continue
// 			}
// 			go func(client *protocol.Client) {
// 				if client.IsAlive == false {
// 					fmt.Printf("client %d exit \n", client.ClientID)
// 					server.DeleteClient(client.ClientID)
// 				} else {
// 					client.RecvPacket <- Message
// 					client.IsAlive = false
// 				}
// 				time.Sleep(time.Second * 10)
// 			}(client)
// 		}
// 		time.Sleep(time.Second * 10)
// 	}
// }
