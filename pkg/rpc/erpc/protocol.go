// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc

import (
	"bytes"
	"encoding/gob"
)

// Data will be encoded to byte stream and being
// transported between client and server
type Data struct {
	Name string        // Service name
	Args []interface{} // Request's or response's params except error
	Err  string        // Error response by server side
}

// encode marshal Data struct to byte
func encode(data Data) ([]byte, error) {
	// Use encoding/gob as encoder which manager binary value
	// exchanged between client and server.
	// https://golang.org/pkg/encoding/gob/
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decode unmarshal byte to Data struct
func decode(b []byte) (Data, error) {
	var data Data
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&data)
	if err != nil {
		return Data{}, err
	}
	return data, nil
}
