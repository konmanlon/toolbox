package common_test

import (
	"fmt"
	"testing"
	"toolbox/common"
	"toolbox/config"
)

func TestNotification(t *testing.T) {
	config.LoadConfig("/root/go/src/toolbox/config.yaml")

	text := fmt.Sprintf(
		"\\#CMTI\n```\nNumber: %s\n\nContent: %s\n\nTime: %s\n```",
		"10683487000000000558",
		"How are you doing, Today",
		"Tue, 09 Jul 2024 22:11:39 +0800",
	)

	err := common.DefaultNotify().SendToTelegram(common.TEXT, []byte(text))
	if err != nil {
		fmt.Println(err)
	}
}
