package otp

import (
	"fmt"
	"reflect"
)

type GenServerStruct struct{
	os *OtpStructs
	castpid chan interface{}
	infopid chan interface{}
	genServer GenServer
}

type GenServer interface {
	Init()
}

func (otpMgr *OtpStructs)NewGenServer(mod string, gServer GenServer){
	otpMgr.mu.Lock()
	defer otpMgr.mu.Unlock()

	gsStruct := &GenServerStruct{}
	gsStruct.genServer = gServer
	gsStruct.os = otpMgr

	otpMgr.gshandlers[mod] = gsStruct

	gsStruct.start()

	otpMgr.wg.Wait()
}

func (gs *GenServerStruct) start() {
	gs.os.wg.Add(1)
	cpid := make(chan interface{}, 10)
	ipid := make(chan interface{}, 10)
	gs.castpid = cpid
	gs.infopid = ipid
	go gs.gen_server()
}


func (otpMgr *OtpStructs) GenServerCast(mod string, args interface{}) {
	if gs, err := otpMgr.gshandlers[mod]; !err{
		fmt.Println("cast error mod exist:", err)
	}else{
		gs.castpid <- args
	}
}

func (otpMgr *OtpStructs) GenServerInfo(mod string, args interface{}) {
	if gs, err := otpMgr.gshandlers[mod]; !err{
		fmt.Println("info error mod exist:", err)
	}else{
		gs.infopid <- args
	}
}

func (gs *GenServerStruct)gen_server(){
	gs.genServer.Init()
	gs.os.wg.Done()

	for{
		select{
		case msg := <- gs.castpid:
			callFunc(msg)
		case msg := <- gs.infopid:
			elem := reflect.ValueOf(msg)
			in := make([]reflect.Value, 1)
			in[0] = elem
			reflect.TypeOf(msg).Method(0).Func.Call(in)
		default:
			continue
		}
	}
}

func callFunc(msg interface{}){
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 1)
	in[0] = elem
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}
