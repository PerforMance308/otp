package main

import (
	"otp"
	"test/tcp"
	"test/configs"
	"sync"
)

type appArgs struct{
	name string
	wg *sync.WaitGroup
}

func main(){
	var wg sync.WaitGroup
	otp.UseOtp()
	wg.Add(1)
	otp.ApplicationStart(configs.TEST_APP_NAME, &appArgs{configs.TEST_APP_NAME, &wg})
	wg.Wait()
}

func (appargs *appArgs) StartApp() error{
	startServer(appargs.wg)
	appargs.wg.Done()
	return nil
}

func startServer(wg *sync.WaitGroup){
	wg.Add(1)
	tcp.StartServer(wg, configs.SERVER_NAME)
}
