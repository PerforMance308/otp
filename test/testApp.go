package main

import (
	"otp"
	"test/role"
	"time"
)

type appArgs struct{
	name string
}

const TEST_APP_NAME = "testapp"

func main(){
	otp.UseOtp()
	otp.ApplicationStart(TEST_APP_NAME, &appArgs{TEST_APP_NAME})
	time.Sleep(time.Second * 5)
}

func (appargs *appArgs) StartApp() error{
	startRole1(appargs.name)
	startRole2(appargs.name)
	return nil
}

func startRole1(appName string){
	role.StartRole(appName, "role1")
}

func startRole2(appName string){
	role.StartRole(appName, "role2")
}