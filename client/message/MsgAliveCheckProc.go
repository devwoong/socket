package message

import (
	"fmt"
	"net"
	"os"
	"socket/client/protocol"
	"time"
)

type AliveCheck struct {
	MessageProc
	ReConnectCnt int
	IsAlivce     bool
}

func (ac *AliveCheck) RecvMessage(msg protocol.Message) (bool, protocol.Message) {
	fmt.Printf("isAlive \n")
	resultPack := protocol.Message{}
	switch msg.Msg {
	case "&&ISALIVE%%":
		resultPack.Msg = "&&ALIVE&&"
	case "&&ALIVE&&":
		ac.IsAlivce = true
	}
	return false, resultPack
}

func (ac *AliveCheck) SendMessage(conn net.Conn) {
	for {
		sendMsg := protocol.Message{}
		sendMsg.Msg = "&&ISALIVE&&"
		conn.Write(sendMsg.Pack().Data[:])

		time.Sleep(time.Second * 10)
		if ac.IsAlivce == false {
			fmt.Printf("re connecting.... %d \n", ac.ReConnectCnt)
			protocol.GetInstance().ConnectServer()
			ac.ReConnectCnt++
			if ac.ReConnectCnt >= 3 {
				fmt.Printf("re connect fail.... exit...\n")
				os.Exit(0)
			}
		} else {
			ac.IsAlivce = false
			ac.ReConnectCnt = 0
		}
	}
}
