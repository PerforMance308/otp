package otp

import (
	"fmt"
	"reflect"
	"sync"
)

type GenServerStruct struct{
	Name string
	app *application
	castpid chan interface{}
	infopid chan interface{}
	callpid chan *callMsg
	genServer GenServer
	mu sync.RWMutex
	flag chan string
}

type GenServer interface {
	Init() error
	Terminate() error
}

type callMsg struct {
	from chan interface{}
	msg interface{}
}

func NewGenServer(appName string, name string, gServer GenServer){
	if app, exist := OtpMgr.applications[appName]; !exist{
		panic("application not start")
	}else{
		app.mu.Lock()
		defer app.mu.Unlock()

		gsStruct := &GenServerStruct{}
		gsStruct.Name = name
		gsStruct.genServer = gServer
		gsStruct.app = app
		gsStruct.flag = make(chan string, 10)

		app.gshandlers[name] = gsStruct

		gsStruct.start()
	}
}

func Terminate(appName string, name string){
	if app, exist := OtpMgr.applications[appName]; !exist{
		panic("application not start")
	}else{
		app.mu.Lock()
		defer app.mu.Unlock()

		if gs, err := app.gshandlers[name]; err{
			delete(app.gshandlers, name)
			if err1 := gs.genServer.Terminate(); err1 != nil{
				gs.flag <- "stop"
			}
		}
	}
}

func (gs *GenServerStruct) start() {
	cpid := make(chan interface{}, 10)
	ipid := make(chan interface{}, 10)
	capid := make(chan *callMsg, 10)
	gs.castpid = cpid
	gs.infopid = ipid
	gs.callpid = capid
	go gs.gen_server()
}


func (app *application) GenServerCast(mod string, args interface{}) {
	if gs, err := app.gshandlers[mod]; !err{
		fmt.Println("cast error mod exist:", err)
	}else{
		gs.castpid <- args
	}
}

func (app *application) GenServerInfo(mod string, args interface{}) {
	if gs, err := app.gshandlers[mod]; !err{
		fmt.Println("info error mod exist:", err)
	}else{
		gs.infopid <- args
	}
}

func (app *application) GenServerCall(mod string, args interface{}) interface{}{
	if gs, err := app.gshandlers[mod]; !err{
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

	for{
		select{
		case msg := <- gs.castpid:
			castFunc(gs, msg)
		case msg := <- gs.infopid:
			infoFunc(gs, msg)
		case msg := <- gs.callpid:
			from := msg.from
			callFunc(gs, from, msg.msg)
		case <- gs.flag:
			break;
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
