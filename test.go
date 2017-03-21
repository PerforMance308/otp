package main

import (
	"otp"
	"fmt"
	"time"
)

type test struct{
	arg1 int
	arg2 string
}

type callMessage struct{
	arg1 int32
	arg2 string
}

type castMessage struct{
	arg1 int32
	arg2 string
}

type infoMessage struct{
	arg1 int32
	arg2 string
	arg3 string
}

const mod = "test"

var otpMgr *otp.OtpStructs

func main(){
	otpMgr = otp.UseOtp()
	MyTest()
}

func MyTest(){
	otpMgr.NewGenServer(mod, test{arg1: 1, arg2: "aa"})
	cast()
	info()
	call()
	time.Sleep(time.Second * 10)
}

func call(){
	returnMsg := otpMgr.GenServerCall(mod, &callMessage{11, "aa"})
	fmt.Println("recv call msg:", returnMsg.(callMessage).arg1, "+", returnMsg.(callMessage).arg2)
}

func cast(){
	otpMgr.GenServerCast(mod, &castMessage{11, "aa"})
}

func info(){
	otpMgr.GenServerInfo(mod, &infoMessage{111, "aaa", "bbb"})
}

func (gs test)Init(){
	fmt.Println("gen server init arg1:", gs.arg1,"arg2:", gs.arg2)
}

func (cm *callMessage)HandleCall(from chan interface{}){
	fmt.Println("handle call", cm.arg1, "+", cm.arg2)
	from <- callMessage{cm.arg1, cm.arg2}
}

func (im *infoMessage)HandleInfo(){
	fmt.Println("handle info", im.arg1, "+", im.arg2, "+", im.arg3)
}

func (cm *castMessage)HandleCast(){
	fmt.Println("handle cast", cm.arg1, "+", cm.arg2)
}
