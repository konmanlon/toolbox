package modem_test

import (
	"fmt"
	"testing"
	"toolbox/config"
	"toolbox/modem"
)

func TestAir72x(t *testing.T) {
	config.LoadConfig("/root/go/src/toolbox/config.yaml")

	dev := modem.DeviceAir72x{
		CommandPort: "/dev/ttyUSB0",
		NotifyPort:  "/dev/ttyUSB0",
	}

	if err := dev.InitDevice(); err != nil {
		fmt.Println(err)
		return
	}

	dev.Watch()

}
