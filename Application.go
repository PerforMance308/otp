package otp

import (
	"sync"
)

type App interface {
	Start(ap *Application) error
}

type Application struct{
	otpMgr *OtpStructs
	args App
	mu sync.RWMutex
}

func (ap *Application)Start(){
	if err := ap.args.Start(ap); err != nil{
		panic("start app error")
	}
}


