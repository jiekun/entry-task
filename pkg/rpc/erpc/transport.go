// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc

import (
	"bufio"
	"encoding/binary"
	"net"
)

// Transport stored a TCP connection
type Transport struct {
	conn net.Conn
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// Send prepare a data struct with header recording body length.
// Each package sent out will have a fixed length for header
// and variable length for body.
// See TLV: https://en.wikipedia.org/wiki/Type-length-value
func (t *Transport) Send(req Data) error {
	// Generate byte body
	b, err := encode(req)
	if err != nil {
		return err
	}

	// Generate binary request: 4 byte header + body
	buf := make([]byte, 4+len(b))

	// Generate a header carrying body length info
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b))) // First 4 byte as header

	// Merge header and body
	copy(buf[4:], b)

	// Write binary request to connection
	_, err = t.conn.Write(buf)
	return err
}

// Receive read byte data from connection and transform
// them into Data struct.
func (t *Transport) Receive() (Data, error) {
	// Read the header from a connection
	header := make([]byte, 4)
	reader := bufio.NewReader(t.conn)
	_, err := reader.Read(header)
	//_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return Data{}, err
	}

	// Read the body length from header
	bodyLen := binary.BigEndian.Uint32(header)

	// Read bodyLen size data from connection
	// and decode to Data struct
	byteData := make([]byte, bodyLen)
	_, err = reader.Read(byteData)
	//_, err = io.ReadFull(t.conn, byteData)
	if err != nil {
		return Data{}, nil
	}
	data, err := decode(byteData)
	return data, err
}
