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
	otpMgr.GenServerCast(mod, 11, "22")
}

func info(){
	otpMgr.GenServerInfo(mod, "aa", 343)
}

func (gs test)Init(){
	fmt.Println("arg1:", gs.arg1)
	fmt.Println("arg2:", gs.arg2)
}

func (gs test)HandleInfo(a ...interface{}){
	fmt.Println("handle info")
}

func (gs test)HandleCast(a ...interface{}){
	fmt.Println("handle cast")
}
