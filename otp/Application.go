package otp

import "sync"

type Application interface {
	StartApp() error
}

type application struct{
	otpMgr *OtpStructs
	app Application
	gshandlers map[string]*GenServerStruct
	mu sync.RWMutex
}

func ApplicationStart(appName string, app Application) {
	application := &application{}
	application.gshandlers = make(map[string]*GenServerStruct)
	application.app = app
	if _,exist := OtpMgr.applications[appName]; exist{
		panic("application already start")
	}

	go application.Start()

	OtpMgr.applications[appName] = application
}

func (ap *application)Start(){
	if err := ap.app.StartApp();err !=nil {
		panic("application start failed")
	}
}


