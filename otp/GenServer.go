package otp

import (
	"fmt"
	"reflect"
	"sync"
)

type GenServerStruct struct{
	Name string
	os *OtpStructs
	castpid chan interface{}
	infopid chan interface{}
	callpid chan *callMsg
	genServer GenServer
	mu sync.RWMutex
}

type GenServer interface {
	Init()
}

type callMsg struct {
	from chan interface{}
	msg interface{}
}

func (otpMgr *OtpStructs)NewGenServer(name string, gServer GenServer){
	otpMgr.mu.Lock()
	defer otpMgr.mu.Unlock()

	gsStruct := &GenServerStruct{}
	gsStruct.Name = name
	gsStruct.genServer = gServer
	gsStruct.os = otpMgr

	otpMgr.gshandlers[name] = gsStruct

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
			castFunc(gs, msg)
		case msg := <- gs.infopid:
			infoFunc(gs, msg)
		case msg := <- gs.callpid:
			from := msg.from
			callFunc(gs, from, msg.msg)
		default:
			continue
		}
	}
}

func castFunc(gs *GenServerStruct, msg interface{}){
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 2)
	in[0] = elem
	in[1] = reflect.ValueOf(gs)
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func infoFunc(gs *GenServerStruct, msg interface{}){
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 2)
	in[0] = elem
	in[1] = reflect.ValueOf(gs)
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func callFunc(gs *GenServerStruct, from chan interface{}, msg interface{}){
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 3)
	in[0] = elem
	in[1] = reflect.ValueOf(from)
	in[2] = reflect.ValueOf(gs)
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}
