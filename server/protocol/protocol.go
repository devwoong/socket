package protocol

import "net"

type Client struct {
	Conn       net.Conn
	ClientID   int
	SendPacket chan Message
	RecvPacket chan Message
	IsAlive    bool
	Exit       chan bool
}

type Server struct {
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
	close(s.Clients[id].SendPacket)
	s.Clients[id].Conn.Close()

	delete(s.Clients, id)
}

func (c *Client) CreateClient(conn net.Conn, id int) {
	c.Conn = conn
	c.ClientID = id
	c.IsAlive = true
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
