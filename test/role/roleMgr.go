package role

import (
	"otp"
	"fmt"
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

func StartRole(appName string, name string){
	otp.NewGenServer(appName, name, args{name})
}

func (gs args)Init() error{
	fmt.Println("gen server init name:", gs.name)
	return nil
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
