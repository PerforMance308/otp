package tcp

import (
	"otp"
	"sync/atomic"
	"test/configs"
	"fmt"
	"net"
	"gopkg.in/mgo.v2/bson"
	"bufio"
	"sync"
)

type Client struct {
	clientStr string
	id uint64
	uuid bson.ObjectId
	mu sync.RWMutex
	nc net.Conn
	srv *Server
	bw             *bufio.Writer
	receivedBytes  []byte
	muNc           sync.RWMutex
	muClosed       sync.RWMutex
	bwmu       	   sync.RWMutex
	closed         bool
	concurrenceChan  chan int
}

func StartClient(conn net.Conn, sev *Server) *Client {
	id := atomic.AddUint64(&sev.gcid, 1)
	clientStr := fmt.Sprintf("r_%d", id)

	c := &Client{id: id, nc: conn, srv: sev, closed: false}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.uuid = bson.NewObjectId()
	c.bw = bufio.NewWriterSize(c.nc, 1400)
	c.receivedBytes = []byte{}
	c.clientStr = clientStr
	c.concurrenceChan = make(chan int, 3)

	otp.NewGenServer(configs.TEST_APP_NAME, clientStr, c)
	return c
}

func (c *Client)Init() error{
	startClient(c)
	return nil
}

func (c *Client)Terminate() error{
	if c.SetClosed() {
		if c.srv != nil {
			c.srv.removeClient(c)

		}

		if c.netConn() != nil {
			c.closeBW()
			c.closeNetConn()

		}
		c.mu.Lock()
		defer c.mu.Unlock()

		c.receivedBytes = nil
		close(c.concurrenceChan)
		c.concurrenceChan = nil
	}
	return nil
}