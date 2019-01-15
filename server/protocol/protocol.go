package protocol

import "net"

type Client struct {
	Conn       net.Conn
	ClientID   int
	SendPacket chan Message
	RecvPacket chan Message
}

type Server struct {
	Exit    map[int]chan bool
	Clients map[int]*Client
}

func (s *Server) Initialize() {
	s.Clients = make(map[int]*Client)
}
func (s *Server) AddClient(client *Client) {
	s.Clients[client.ClientID] = client
}
func (s *Server) DeleteClient(id int) {
	close(s.Clients[id].RecvPacket)
	close(s.Exit[id])

	delete(s.Clients, id)
	delete(s.Exit, id)
}

func (c *Client) CreateClient(conn net.Conn, id int) {
	c.Conn = conn
	c.ClientID = id
	c.SendPacket = make(chan Message)
	c.RecvPacket = make(chan Message)
}
func (c *Client) GetSendPacket() chan Message {
	return c.SendPacket
}
func (c *Client) GetRecvPacket() chan Message {
	return c.RecvPacket
}
func (c *Client) GetConnection() net.Conn {
	return c.Conn
}
