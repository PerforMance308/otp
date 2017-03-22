package main

import (
	"otp"
	"sync"
	"testgs"
)

type appArgs struct{
}

func main(){
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()

	otpMgr := otp.UseOtp()
	otpMgr.ApplicationStart(&appArgs{})

	for true{
	}

	wg.Wait()
}


func (t *appArgs) Start (ap *otp.Application) error {
	testgs.Start(ap)
	return nil
}