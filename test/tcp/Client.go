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
	"time"
	"io"
)

type Client struct {
	clientStr string
	id uint64
	uuid bson.ObjectId
	mu sync.RWMutex
	nc net.Conn
	srv *Server
	bw             *bufio.Writer
	ReceivedBytes  []byte
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
	c.ReceivedBytes = []byte{}
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

		c.ReceivedBytes = nil
		close(c.concurrenceChan)
		c.concurrenceChan = nil
	}
	return nil
}



////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func startClient(c *Client){
	c.readLoop()
}

func (c *Client) netConn() net.Conn {
	c.muNc.RLock()

	defer c.muNc.RUnlock()
	return c.nc
}

func (c *Client) Closed() bool {
	c.muClosed.RLock()
	defer c.muClosed.RUnlock()
	return c.closed
}

func (c *Client) SetClosed() bool {
	c.muClosed.Lock()
	defer c.muClosed.Unlock()
	if c.closed == true {
		return false
	}
	c.closed = true
	return true
}

func (c *Client) closeBW() {
	c.bwmu.Lock()
	defer c.bwmu.Unlock()
	c.bw.Flush()
	c.bw = nil
}

func (c *Client) closeNetConn() {
	c.muNc.Lock()
	defer c.muNc.Unlock()
	c.nc = nil
}

func (c *Client) CloseConnection(after time.Duration, wait bool) {
	if after > 0 {
		time.Sleep(after)
	}

	if c.Closed() {
		return
	}

	otp.Terminate(configs.TEST_APP_NAME, c.clientStr)
}

func (c *Client)readLoop(){

	if c.netConn() == nil {
		return
	}

	bytes := make([]byte, 32768)
	for c.srv.isRunning() && !c.Closed() {
		c.netConn().SetReadDeadline(time.Now().Add(5 * time.Minute))
		i, err := c.netConn().Read(bytes)
		if err != nil {
			if err != io.EOF {
				c.CloseConnection(0*time.Second, false)
				return
			}
		}
		if i == 0 {
			continue
		}
		data := bytes[:i]

		fmt.Println("recv:", string(data))
		c.srv.BroadCast(&data)
	}
}

func (c *Client) WriteToNetConn(data *[]byte) {
	if c.netConn() == nil || c.Closed() {
		return
	}
	c.netConn().Write(*data)
	c.bwmu.Lock()

	defer c.bwmu.Unlock()
}