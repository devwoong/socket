package protocol

import "net"

type Client struct {
	Conn       net.Conn
	ClientID   int
	SendPacket Message
	RecvPacket chan Message
}

type Server struct {
	SendMessage map[int]chan Message
	RecvMessage map[int]chan Message
	Exit        map[int]chan bool
	Clients     map[int]*Client
}

func (s *Server) Initialize() {
	s.Clients = make(map[int]*Client)
	s.SendMessage = make(map[int]chan Message)
	s.RecvMessage = make(map[int]chan Message)
}
func (s *Server) AddClient(client *Client) {
	s.Clients[client.ClientID] = client
}
func (s *Server) DeleteClient(id int) {
	close(s.Clients[id].RecvPacket)
	close(s.Exit[id])
	close(s.SendMessage[id])
	close(s.RecvMessage[id])

	delete(s.Clients, id)
	delete(s.Exit, id)
	delete(s.SendMessage, id)
	delete(s.RecvMessage, id)
}
func (s *Server) GetSendMessage() map[int]chan Message {
	return s.SendMessage
}
func (s *Server) GetRecvMessage() map[int]chan Message {
	return s.RecvMessage
}

func (c *Client) CreateClient(conn net.Conn, id int) {
	c.Conn = conn
	c.ClientID = id
	c.SendPacket = Message{}
	c.RecvPacket = make(chan Message)
}
func (c *Client) GetSendPacket() Message {
	return c.SendPacket
}
func (c *Client) GetRecvPacket() chan Message {
	return c.RecvPacket
}
func (c *Client) GetConnection() net.Conn {
	return c.Conn
}
