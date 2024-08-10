package modem

import (
	"fmt"
	"log"
	"strings"
	"time"
	"toolbox/common"

	"github.com/xlab/at"
	"github.com/xlab/at/calls"
	"github.com/xlab/at/sms"
)

type DeviceAir72x struct {
	CommandPort    string
	NotifyPort     string
	initialization bool
	dev            *at.Device
	notifycation   common.Notifyer
	closed         chan struct{}
}

// 必须先调用此方法初始化设备
func (d *DeviceAir72x) InitDevice() (err error) {
	dev := &at.Device{
		CommandPort: d.CommandPort,
		Name:        d.NotifyPort,
	}

	if err = dev.Open(); err != nil {
		return
	}

	if err = dev.Init(at.DeviceAir72x()); err != nil {
		return
	}

	d.dev = dev
	d.initialization = true
	d.notifycation = common.DefaultNotify()
	d.closed = make(chan struct{})

	return
}

func (d *DeviceAir72x) SendSMS() error {
	return nil
}

func (d *DeviceAir72x) sendNotifications(msg any) error {
	switch m := msg.(type) {
	case *calls.CallerID:
		text := fmt.Sprintf("\\#CLIP\n```\nCallerID: %s\n```", m.CallerID)
		return d.notifycation.SendToTelegram(common.TEXT, []byte(text))
	case *sms.Message:
		// 短信发送时间，将类型转换成 time.Time
		date := time.Time(m.ServiceCenterTime)
		text := fmt.Sprintf(
			"\\#CMTI\n```\nNumber: %s\n\nContent: %s\n\nTime: %s\n```",
			m.Address, m.Text, date.Format(time.RFC1123),
		)
		return d.notifycation.SendToTelegram(common.TEXT, []byte(text))
	default:
		return nil
	}
}

// 阻塞函数
func (d *DeviceAir72x) Watch() {
	go func() {
		if err := d.dev.Watch(); err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		if err := d.dev.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case <-d.closed:
			return
		case caller := <-d.dev.IncomingCallerID():
			if err := d.sendNotifications(caller); err != nil {
				log.Println(err)
			}
		case sms := <-d.dev.IncomingSms():
			if err := d.sendNotifications(sms); err != nil {
				log.Println(err)
			}
		}
	}
}

func (d *DeviceAir72x) Close() {
	d.closed <- struct{}{}
}

// 通话录音
func (d *DeviceAir72x) CAUDREC() (str string, err error) {
	str, err = d.dev.Send(`AT+CAUDREC=1,"10001.wav",2,0`)
	return
}

// 获取文件系统剩余空间大小 bytes
func (d *DeviceAir72x) FSMEM() (size string, err error) {
	str, err := d.dev.Send(`AT+FSMEM`)
	size = strings.TrimPrefix(str, `+FSMEM: `)
	return
}

// 获取文件大小 bytes
func (d *DeviceAir72x) FSFLSIZE(filename string) (size string, err error) {
	str, err := d.dev.Send(`AT+FSFLSIZE=` + filename)
	size = strings.TrimPrefix(str, `+FSFLSIZE: `)
	return
}

// 读取文件
// 每次读取不能大于10240字节
func (d *DeviceAir72x) FSREAD(filename string) (str string, err error) {
	str, err = d.dev.Send(`AT+FSREAD=<filename>,<mode>,<filesize>,<position>` + filename)
	return
}

// 删除文件
func (d *DeviceAir72x) FSDEL(filename string) (str string, err error) {
	str, err = d.dev.Send(`AT+FSDEL=` + filename)
	return
}
