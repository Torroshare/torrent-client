package bigend

import (
	"encoding/binary"
	"io"
)

func parseMsg(msg string, buffer []byte) (int, []byte) {
	var index int
	var payload []byte

	offset := 4
	l := MsgByteFill(msg) / offset

	for i := offset * l; i > offset; i /= offset {
		begin := int(binary.BigEndian.Uint32(buffer[i-offset : i]))
		if begin > len(buffer) { //error handler
			continue
		}
		payload = buffer[i:] //now payload not null
	}
	index = int(binary.BigEndian.Uint32(buffer[0:4])) //index is common for both
	return index, payload                             //for have return index,null; for piece return index payload
}

func Read(r io.Reader) (*Message, error) {
	length, err := marginPointer(&r, 4)

	if err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	} // keep-alive message

	id, err := marginPointer(&r, int(length))

	if err != nil {
		return nil, err
	}

	payload := make([]byte, int(length))

	_, err = io.ReadFull(r, payload)

	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(id),
		Payload: payload,
	}

	return &m, nil
}

func marginPointer(r *io.Reader, margin int) (uint32, error) {
	partial := make([]byte, margin)
	_, err := io.ReadFull(*r, partial)
	return binary.BigEndian.Uint32(partial), err
}
