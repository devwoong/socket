package protocol

type Packet struct {
	Data   []byte
	Buffer [512]byte
}

func (p *Packet) UnPack(size int) Message {
	resultMessage := Message{}
	resultMessage.Msg = string(p.Data[:size])
	return resultMessage
}
