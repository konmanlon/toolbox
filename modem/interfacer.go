package modem

type Modem interface {
	InitDevice() error
	// 以下所有方法必须先调用 InitDevice()
	Watch()
	Close()
	SendSMS() error
}
