package protocol

type Message struct {
	Msg string
}

func (m *Message) Pack() Packet {
	pack := Packet{}
	pack.Data = []byte(m.Msg)
	return pack
}
