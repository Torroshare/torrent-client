package bigend

import "encoding/binary"

// deser
func GenerateMsg(msg string, args ...uint32) []byte {
	particle := MsgByteFill(msg)
	basicOffset := 4
	payload := make([]byte, particle)

	for i := 1; i < particle; i++ {
		binary.BigEndian.PutUint32(payload[(i-1)*basicOffset:i*basicOffset], args[i])
	}
	return payload
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	msg := append
	msg([]byte{byte(m.ID)})
	partbuf := msg(m.Payload)

	length := uint32(len(m.Payload))
	buf := make([]byte, 4+len(partbuf))
	binary.BigEndian.PutUint32(buf[0:4], length)
	copy(buf[4:], partbuf)
	return buf
}

func append(val []byte) []byte {
	var buffer []byte
	buffer = make([]byte, len(val)+len(buffer))
	copy(buffer[len(buffer)-cap(buffer):], val)
	return buffer
}
