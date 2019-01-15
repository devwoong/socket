package main

import (
	"bufio"
	"fmt"
	"log"
	"socket/client/message"
	"socket/client/protocol"
	"time"
)

var messageProcs []message.MessageProc

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println(os.Stderr, "Usage %s host:port", os.Args[0])
	// 	os.Exit(0)
	// }
	// service := os.Args[1]

	// service := "localhost:1201"
	err := protocol.GetInstance().ConnectServer()
	if err != nil {

	}
	fmt.Printf("------starting client------")
	defer protocol.GetInstance().Conn.Close()

	messageProcs = []message.MessageProc{
		&message.AliveCheck{ReConnectCnt: 0, IsAlivce: true},
	}

	go recvMessage()
	go recvMsgProc()
	go tick()
	go sendMessage()
	for {

	}
}
func recvMessage() {
	for {
		readPacket := protocol.Packet{}

		r := bufio.NewReader(protocol.GetInstance().Conn)
		readPacket.Data = make([]byte, 1024)
		n, err := r.Read(readPacket.Data)
		if err != nil {
			return
		}
		if n > 0 {
			protocol.GetInstance().RecvMessage <- readPacket.UnPack(n)
			// message := readPacket.UnPack(n)
		}
	}
}

func recvMsgProc() {
RECV_EXIT:
	for {
		select {
		case msg := <-protocol.GetInstance().RecvMessage:
			switch msg.Msg {
			case "&&EXIT&&":
				break RECV_EXIT
			default:
				log.Printf("Receive: %s \n", msg.Msg)
				for _, v := range messageProcs {
					isSend, sendMsg := v.RecvMessage(msg)
					if isSend == true {
						protocol.GetInstance().SendMessage <- sendMsg
					}
				}
			}
		}
	}
}

func tick() {
	for _, v := range messageProcs {
		go v.SendMessage(protocol.GetInstance().Conn)
	}
	for {
		message := protocol.Message{}
		message.Msg = "TICK TICK"
		protocol.GetInstance().SendMessage <- message
		time.Sleep(time.Second * 1)
	}
}

func sendMessage() {
SEND_EXIT:
	for {
		select {
		case msg := <-protocol.GetInstance().SendMessage:
			switch msg.Msg {
			case "&&EXIT&&":
				break SEND_EXIT
			default:
				protocol.GetInstance().Conn.Write(msg.Pack().Data[:])
			}
		}
	}
}
