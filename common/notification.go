package common

import (
	"fmt"
	"toolbox/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageType string

const (
	TEXT  MessageType = "text"
	AUDIO MessageType = "audio"
)

type telegramMessage struct {
	MsgType MessageType
	MsgData []byte
}

type Notifyer interface {
	SendToMail(subject, text string) error
	SendToTelegram(msgType MessageType, data []byte) error
	SendToWechat(text string) error
}

type Notifications struct {
	config.Mail
	config.Telegram
	config.Wechat
}

func DefaultNotify() *Notifications {
	c := config.Config.Notifications
	return &Notifications{
		c.Mail,
		c.Telegram,
		c.Wechat,
	}
}

func (n *Notifications) SendToMail(subject, text string) error {
	return nil
}

func (n *Notifications) sendTelegram(msg *telegramMessage) error {
	bot, err := tgbotapi.NewBotAPI(n.Telegram.BotToken)
	if err != nil {
		return fmt.Errorf("bot token %s", err)
	}

	switch msg.MsgType {
	case "text":
		msg := tgbotapi.NewMessage(n.Telegram.ChatID, string(msg.MsgData))
		msg.ParseMode = n.Telegram.ParseMode
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	case "audio":
		file := tgbotapi.FileBytes{
			Name:  "VoiceRecord.wav",
			Bytes: msg.MsgData,
		}
		voice := tgbotapi.NewAudio(n.Telegram.ChatID, file)
		_, err = bot.Send(voice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Notifications) SendToTelegram(t MessageType, d []byte) error {
	msg := &telegramMessage{
		MsgType: t,
		MsgData: d,
	}
	return n.sendTelegram(msg)
}

func (n *Notifications) SendToWechat(text string) error {
	return nil
}
