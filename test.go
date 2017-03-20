package main

import (
	"otp"
	"fmt"
)

type test struct{
	arg1 int
	arg2 string
}

type message struct{

}

type rpMessage struct{

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
}

func cast(){
	otpMgr.GenServerCast(mod, HandleCast, 33, "22")
}

func info(){
	otpMgr.GenServerInfo(mod, HandleInfo, "aa", 343, 555)
}

func (gs test)Init(){
	fmt.Println("arg1:", gs.arg1)
	fmt.Println("arg2:", gs.arg2)
}

func HandleInfo(arg1 int32, arg2 string){
	fmt.Println("handle info", arg1, "+", arg2)
}

func HandleCast(arg1 string, arg2 int32, arg3 int32){
	fmt.Println("handle cast", arg1, "+", arg2, "+", arg3)
}
