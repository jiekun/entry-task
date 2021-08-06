// @Author: 2014BDuck
// @Date: 2021/8/6

package erpc

import (
	"net"
	"sync"
)

type ConnectionPool struct {
	connections []*net.Conn
	locks       []*sync.Mutex
	size        int
	lastIdx     int
	lock        *sync.Mutex
}

func NewConnectionPool(addr string, connNum int) (*ConnectionPool, error) {
	connections := make([]*net.Conn, connNum)
	locks := make([]*sync.Mutex, connNum)
	for i := 0; i < connNum; i++ {
		connection, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		lock := sync.Mutex{}
		connections[i] = &connection
		locks[i] = &lock
	}
	return &ConnectionPool{connections, locks, connNum, 0, &sync.Mutex{}}, nil
}

func (cp *ConnectionPool) Get() (*net.Conn, *sync.Mutex, error) {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	// Round Robin
	if cp.lastIdx < cp.size-1 {
		cp.lastIdx++
	} else {
		cp.lastIdx = 0
	}
	return cp.connections[cp.lastIdx], cp.locks[cp.lastIdx], nil
}
