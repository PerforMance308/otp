package tcp

import (
	"net"
	"fmt"
	"os"
	"sync"
	"otp"
	"test/configs"
)

type Server struct{
	listener net.Listener
	running bool
	mu sync.RWMutex
	clients  map[uint64]*Client
	gcid     uint64
}

func start(wg *sync.WaitGroup, host string, port int){
	defer wg.Done()
	s := &Server{running:true}
	s.clients = make(map[uint64]*Client)
	hp := fmt.Sprintf("%s:%d", host, port)
	l, e := net.Listen("tcp", hp)
	if e != nil {
		os.Exit(-1)
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()
	fmt.Println("TCP server listen at ", host, ":", port)
	for s.isRunning(){
		conn, err := l.Accept()

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Println("Temporary Client accept Error:", err)
			}else if s.isRunning() {
				fmt.Println("Accept error:", err)
			}
			continue
		}

		s.handleConnect(conn)
	}

	otp.Terminate(configs.TEST_APP_NAME, configs.SERVER_NAME)
}

func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Server) removeClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c.id)
}

func (s *Server) handleConnect(conn net.Conn) {
	c := StartClient(conn, s)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[c.id] = c

}

func (s *Server) BroadCast(msg *[]byte){
	for _, c := range s.clients{
		c.WriteToNetConn(msg)
	}
}

