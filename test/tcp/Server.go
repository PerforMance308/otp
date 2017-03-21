package tcp

import (
	"otp"
	"test/configs"
	"sync"
)

type serverArgs struct{
	name string
	host string
	port int
	wg *sync.WaitGroup
}

func StartServer(wg *sync.WaitGroup, name string){
	host := "127.0.0.1"
	port := 8080
	otp.NewGenServer(configs.TEST_APP_NAME, name, &serverArgs{name:name, host:host, port:port, wg:wg})
}

func (gs *serverArgs)Init() error{
	start(gs.wg, gs.host, gs.port)
	return nil
}

func (gs *serverArgs)Terminate() error{
	return nil
}