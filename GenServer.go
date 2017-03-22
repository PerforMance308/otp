package otp

import (
	"reflect"
	"sync"
)

type GenServerStruct struct {
	App       *Application
	castpid   chan interface{}
	sendpid   chan interface{}
	callpid   chan *callMsg
	genServer GenServer
	mu        sync.RWMutex
	flag      chan string
}

type GenServer interface {
	Init(gs *GenServerStruct) error
	Terminate() error
}

type callMsg struct {
	from chan interface{}
	msg  interface{}
}

func (app *Application) GenServerStart(gServer GenServer) *GenServerStruct {
	app.mu.Lock()
	defer app.mu.Unlock()

	gsStruct := &GenServerStruct{}
	gsStruct.genServer = gServer
	gsStruct.App = app
	gsStruct.flag = make(chan string, 10)

	gsStruct.start()

	return gsStruct
}

func (gs *GenServerStruct) Terminate() {
	if err1 := gs.genServer.Terminate(); err1 != nil {
		gs.flag <- "stop"
	}
}

func (gs *GenServerStruct) start() {
	cpid := make(chan interface{}, 10)
	ipid := make(chan interface{}, 10)
	capid := make(chan *callMsg, 10)
	gs.castpid = cpid
	gs.sendpid = ipid
	gs.callpid = capid
	go gs.gen_server()
}

func (gs *GenServerStruct) Call(args interface{}) interface{} {
	fr := make(chan interface{}, 10)
	defer close(fr)

	callMsg := &callMsg{fr, args}
	gs.callpid <- callMsg
	rMsg := <-fr
	return rMsg
}

func (gs *GenServerStruct) Cast(args interface{}) {
	gs.castpid <- args
}

func (gs *GenServerStruct) Send(args interface{}) {
	gs.sendpid <- args
}

func (gs *GenServerStruct) gen_server() {
	gs.genServer.Init(gs)

	for {
		select {
		case msg := <-gs.castpid:
			castFunc(gs, msg)
		case msg := <-gs.sendpid:
			sendFunc(gs, msg)
		case msg := <-gs.callpid:
			from := msg.from
			callFunc(gs, from, msg.msg)
		case <-gs.flag:
			break
		default:
			continue
		}
	}
}

func callFunc(gs *GenServerStruct, from chan interface{}, msg interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 2)
	in[0] = elem
	in[1] = reflect.ValueOf(from)
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func castFunc(gs *GenServerStruct, msg interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 1)
	in[0] = elem
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}

func sendFunc(gs *GenServerStruct, msg interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	elem := reflect.ValueOf(msg)
	in := make([]reflect.Value, 1)
	in[0] = elem
	reflect.TypeOf(msg).Method(0).Func.Call(in)
}
