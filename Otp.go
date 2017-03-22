package otp

import (
	"sync"
)

type OtpStructs struct{
	mu sync.RWMutex
}

func UseOtp() *OtpStructs {
	return &OtpStructs{}
}

func (otpMgr *OtpStructs) ApplicationStart(args App){
	application := &Application{otpMgr: otpMgr}
	application.args = args

	go application.Start()
}


