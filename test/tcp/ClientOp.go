package tcp

import (
	"time"
	"fmt"
	"io"
	"net"
	"otp"
	"test/configs"
)

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

		c.receivedBytes = append(c.receivedBytes, data...)
		fmt.Println("recv:", c.receivedBytes)
	}
}