package message

import (
	"net"
	"socket/client/protocol"
)

type MessageProc interface {
	RecvMessage(msg protocol.Message) (bool, protocol.Message)
	SendMessage(conn net.Conn)
}
