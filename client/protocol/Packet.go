package protocol

type Packet struct {
	Data   []byte
	Buffer [512]byte
}

func (p *Packet) UnPack() Message {
	resultMessage := Message{}
	resultMessage.Msg = string(p.Data[:])
	return resultMessage
}
