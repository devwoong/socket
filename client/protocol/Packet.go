package protocol

type Packet struct {
	Data   []byte
	Buffer [512]byte
}

type Message struct {
	Msg string
}

func (p *Packet) UnPack() Message {
	resultMessage := Message{}
	resultMessage.Msg = string(p.Data[:])
	return resultMessage
}

func (m *Message) Pack() Packet {
	pack := Packet{}
	pack.Data = []byte(m.Msg)
	return pack
}
