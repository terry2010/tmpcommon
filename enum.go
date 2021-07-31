package Common

type serverRole struct {
	value string
}

var ServerRole = struct {
	Unknown serverRole
	Master  serverRole
	Slave   serverRole
	Api     serverRole
}{
	Unknown: serverRole{"Unknown"},
	Master:  serverRole{"Master"},
	Slave:   serverRole{"Slave"},
	Api:     serverRole{"Api"},
}

func (s serverRole) String() string {
	return s.value
}

type taskOpenType struct {
	value string
}

var TaskOpenType = struct {
	Proxy  taskOpenType
	Direct taskOpenType
}{
	Proxy:  taskOpenType{"Proxy"},
	Direct: taskOpenType{"Direct"},
}

func (t taskOpenType) String() string {
	return t.value
}



type clientOpration struct {
	value string
}

var ClientOpration = struct {
	Navigate clientOpration
	Client  clientOpration
	Sleep clientOpration
	FilterURL clientOpration
	ScreenShot clientOpration
	GetHTML clientOpration
	UploadSAE clientOpration
	Callback clientOpration
	AutoResize clientOpration
}{
	Navigate: clientOpration{"navigate"},
	Client:  clientOpration{"client"},
	Sleep: clientOpration{"sleep"},
	FilterURL:clientOpration{"filter"},
	ScreenShot:clientOpration{"screenshot"},
	GetHTML:clientOpration{"getHTML"},
	UploadSAE:clientOpration{"uploadSAE"},
	Callback:clientOpration{"callback"},
	AutoResize:clientOpration{"autoResize"},
}

func (c clientOpration) String() string {
	return c.value
}



type clientOprationParam struct {
	value string
}

var ClientOprationParam = struct {
	Time clientOprationParam
	Client  clientOprationParam
	URL clientOprationParam
	Finish clientOprationParam
	Running clientOprationParam
	Success clientOprationParam
	Fail clientOprationParam
	Status clientOprationParam
	On clientOprationParam
	Off clientOprationParam
	Kill clientOprationParam
}{
	Time: clientOprationParam{"time"},
	Client:  clientOprationParam{"client"},
	URL: clientOprationParam{"url"},
	Finish: clientOprationParam{"Finish"},
	Running: clientOprationParam{"running"},
	Success: clientOprationParam{"success"},
	Fail: clientOprationParam{"fail"},
	Status: clientOprationParam{"status"},
	On: clientOprationParam{"on"},
	Off: clientOprationParam{"off"},
	Kill: clientOprationParam{"kill"},
}

func (c clientOprationParam) String() string {
	return c.value
}
