package otp

import (
	"fmt"
	"reflect"
)

type GenServerStruct struct{
	os *OtpStructs
	castpid chan interface{}
	infopid chan interface{}
	callpid chan *callMsg
	genServer GenServer
}

type GenServer interface {
	Init()
}

type callMsg struct {
	from chan interface{}
	msg interface{}
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
	capid := make(chan *callMsg, 10)
	gs.castpid = cpid
	gs.infopid = ipid
	gs.callpid = capid
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

func (otpMgr *OtpStructs) GenServerCall(mod string, args interface{}) interface{}{
	if gs, err := otpMgr.gshandlers[mod]; !err{
		fmt.Println("call error mod exist:", err)
		return nil
	}else{
		fr := make(chan interface{}, 10)
		defer close(fr)

		callMsg := &callMsg{fr, args}
		gs.callpid <- callMsg
		rMsg := <- fr
		return rMsg
	}
}


func (gs *GenServerStruct)gen_server(){
	gs.genServer.Init()
	gs.os.wg.Done()

	for{
		select{
		case msg := <- gs.castpid:
			castFunc(msg)
		case msg := <- gs.infopid:
			infoFunc(msg)
		case msg := <- gs.callpid:
			from := msg.from
			callFunc(from, msg.msg)
		default:
			continue
		}
	}
}

func castFunc(msg interface{}){
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 1)
	in[0] = elem
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func infoFunc(msg interface{}){
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 1)
	in[0] = elem
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func callFunc(from chan interface{}, msg interface{}){
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 2)
	in[0] = elem
	in[1] = reflect.ValueOf(from)
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}
