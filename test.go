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
	time.Sleep(time.Second * 10)
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

func (im *infoMessage)HandleInfo(){
	fmt.Println("handle info", im.arg1, "+", im.arg2, "+", im.arg3)
}

func (cm *castMessage)HandleCast(){
	fmt.Println("handle cast", cm.arg1, "+", cm.arg2)
}
