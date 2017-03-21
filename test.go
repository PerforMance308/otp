package main

import (
	"otp"
	"fmt"
	"time"
)

type args struct{
	name string
}

type callMessage struct{
	from string
}

type castMessage struct{
	from string
}

type infoMessage struct{
	from string
}

var otpMgr *otp.OtpStructs

func main(){
	otpMgr = otp.UseOtp()
	myTest1()
	myTest2()
	time.Sleep(time.Second * 10)
	do1()
	time.Sleep(time.Second * 10)
	do2()
	time.Sleep(time.Second * 10)
}

func myTest1(){
	otpMgr.NewGenServer("test1", args{"test1"})
}

func myTest2(){
	otpMgr.NewGenServer("test2", args{"test2"})
}

func do1(){
	cast1()
	info1()
	call1()
}

func do2(){
	cast2()
	info2()
	call2()
}

func call1(){
	returnMsg := otpMgr.GenServerCall("test1", &callMessage{"test2"})
	fmt.Println("recv call msg from:", returnMsg.(callMessage).from)
}

func cast1(){
	otpMgr.GenServerCast("test1", &castMessage{"test2"})
}

func info1(){
	otpMgr.GenServerInfo("test1", &infoMessage{"test2"})
}

func call2(){
	returnMsg := otpMgr.GenServerCall("test2", &callMessage{"test1"})
	fmt.Println("recv call msg from:", returnMsg.(callMessage).from)
}

func cast2(){
	otpMgr.GenServerCast("test2", &castMessage{"test1"})
}

func info2(){
	otpMgr.GenServerInfo("test2", &infoMessage{"test1"})
}

func (gs args)Init(){
	fmt.Println("gen server init my name:", gs.name)
}

func (cm *callMessage)HandleCall(from chan interface{}, gs *otp.GenServerStruct){
	fmt.Println(gs.Name, ":handle call from:", cm.from)
	from <- callMessage{cm.from}
}

func (im *infoMessage)HandleInfo(gs *otp.GenServerStruct){
	fmt.Println(gs.Name, ":handle info from", im.from)
}

func (cm *castMessage)HandleCast(gs *otp.GenServerStruct){
	fmt.Println(gs.Name, ":handle cast from", cm.from)
}
