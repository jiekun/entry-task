// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// Call receive RPC service's name and a function pointer,
// bind RPC calling to this function.
// See: https://golang.org/pkg/reflect/#example_MakeFunc
func (c *Client) Call(serviceName string, funcPtr interface{}) {
	// Get the Value of funcPtr
	// See: https://golang.org/pkg/reflect/#Value
	v := reflect.ValueOf(funcPtr).Elem()
	f := func(args []reflect.Value) []reflect.Value {
		clientTrans := NewTransport(c.conn)
		numOut := v.Type().NumField()

		// Output length is specified. Build an output
		// when error happened with zero value.
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, numOut)
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(v.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		sendArgs := make([]interface{}, 0, len(args))
		for i := range args {
			// Transform arguments into interface{} type
			// so that it can match Data.Args
			sendArgs = append(sendArgs, args[i].Interface())
		}
		err := clientTrans.Send(Data{
			Name: serviceName,
			Args: sendArgs,
		})
		if err != nil {
			return errorHandler(err)
		}

		respData, err := clientTrans.Receive()
		if err != nil {
			return errorHandler(err)
		}
		if respData.Err != "" {
			return errorHandler(errors.New(respData.Err))
		}

		// No err && no reply data
		if len(respData.Args) == 0 {
			respData.Args = make([]interface{}, numOut)
		}

		replyArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				if respData.Args[i] == nil {
					replyArgs[i] = reflect.Zero(v.Type().Out(i))
				} else {
					replyArgs[i] = reflect.ValueOf(respData.Args[i])
				}
			} else {
				replyArgs[i] = reflect.Zero(v.Type().Out(i))
			}
		}
		return replyArgs
	}
	// Replace funcPtr's Value with f
	v.Set(reflect.MakeFunc(v.Type(), f))
}
