package otp

import "sync"

type OtpStructs struct{
	mu sync.RWMutex
	applications map[string]*application
}

var OtpMgr *OtpStructs

func UseOtp(){
	app := make(map[string]*application)
	OtpMgr = &OtpStructs{applications:app}
}



