package protocol

import (
	"fmt"
	"net"
	"os"
)

type server struct {
	Service     string
	Conn        net.Conn
	RecvMessage chan Message
	SendMessage chan Message
}

func (s *server) ConnectServer() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.Service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	s.Conn = conn
	return err
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s\n", err.Error())
	}
}

var serverInstance *server

func GetInstance() *server {
	if serverInstance == nil {
		serverInstance = &server{"localhost:1201", nil, make(chan Message), make(chan Message)}
	}
	return serverInstance
}
