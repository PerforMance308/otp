package otp

import "sync"

type OtpStructs struct{
	mu sync.RWMutex
	gshandlers map[string]*GenServerStruct
	fsmhandlers map[string]*FsmStruct
	wg *sync.WaitGroup
}

type FsmStruct struct{
	name string
	pid chan []interface{}
}

func UseOtp() *OtpStructs{
	var waitgroup sync.WaitGroup
	gsHandlers := make(map[string]*GenServerStruct)
	fsmHandlers := make(map[string]*FsmStruct)
	return &OtpStructs{gshandlers: gsHandlers, fsmhandlers: fsmHandlers, wg: &waitgroup}
}



