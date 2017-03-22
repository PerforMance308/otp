package testgs

import (
	"fmt"
	"otp"
	"time"
)

type gsArgs struct{
}

type callMsg struct{
	arg string
}

type callReturnMsg struct{
	arg string
}

type castMsg struct{
	arg string
}
type sendMsg struct{
	arg string
}


var _gs *otp.GenServerStruct

func Start(ap *otp.Application){
	ap.GenServerStart(&gsArgs{})
	time.Sleep(time.Second *2)
	call()
	time.Sleep(time.Second *2)
	cast()
	time.Sleep(time.Second *2)
	send()
	time.Sleep(time.Second *2)
	_gs.Terminate("gen1")
}

func call(){
	returnMsg := _gs.Call(&callMsg{"call  fdsfd"})
	fmt.Println("call return msg:", returnMsg.(callReturnMsg).arg)
}

func cast(){
	_gs.Cast(&castMsg{"cast  fdsfd"})
}

func send(){
	_gs.Send(&sendMsg{"send  fdsfd"})
}

////////////////////////////////////////////////////////////////////////
//			Init
////////////////////////////////////////////////////////////////////////
func (gsa *gsArgs) Init (gs *otp.GenServerStruct) error {
	fmt.Println("gen server start")
	_gs = gs
	return nil
}

////////////////////////////////////////////////////////////////////////
//			Handler
////////////////////////////////////////////////////////////////////////
func (msg *callMsg)HandleCall(from chan interface{}){
	fmt.Println("recv call msg:",msg.arg)
	from <- callReturnMsg{"vvvvv"}
}

func (msg *castMsg)HandleCast(){
	fmt.Println("recv cast msg:",msg.arg)
}

func (msg *sendMsg)HandleInfo(){
	fmt.Println("recv send msg:",msg.arg)
}

////////////////////////////////////////////////////////////////////////
//			Exit
////////////////////////////////////////////////////////////////////////
func (gsa *gsArgs) Terminate () error {
	fmt.Println("gen server stop")
	return nil
}
