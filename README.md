# erlang otp writen by go
# useage
1.Application
```
func main(){
  otp.UseOtp()
  //start a new application
  otp.ApplicationStart(appName, &appArgs{})
}
//You need to implement the following several methods
func (args *appArgs)StartApp() error {
  return nil
}
```
2.gen_server
```
func start(){
  otp.NewGenServer(appName, genServerName, &args{})
}

//You need to implement the following several methods
func (gs *args)Init() error{
	return nil
}

func (gs *args)Terminate(){
}

func call(){
  returnMsg := otp.GenServerCall(appName, genServerName, callArgs{})
}

func cast(){
  otp.GenServerCast(appName, genServerName, castArgs{})
}

func info(){
  otp.GenServerInfo(appName, genServerName, infoArgs{})
}

//Handlers
func (args *callArgs) HandleCall(from chan interface{}){
  from <- result
}

func (args *castArgs) HandleCast(){
}

func (args *infoArgs) HandleInfo(){
}
```
