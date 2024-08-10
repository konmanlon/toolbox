package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config *FullConfig

type Authorize struct {
	ApiKey string `yaml:"apiKey"`
	Email  string `yaml:"email"`
}

type DNS struct {
	Domain        string `yaml:"domain"`
	RecordName    string `yaml:"recordName"`
	Endpoint      string `yaml:"endpoint"`
	ScheduledTask int64  `yaml:"scheduledTask"`
}

type Cloudflare struct {
	Authorize `yaml:"authorize"`
	DNS       `yaml:"dns"`
}

type DDNS struct {
	Enable     bool `yaml:"enable"`
	Cloudflare `yaml:"cloudflare"`
}

type Modem struct {
	Enable     bool   `yaml:"enable"`
	SerialPort string `yaml:"serialPort"`
}

type Mail struct {
	Smtp     string `yaml:"smtp"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Cc       string `yaml:"cc"`
}

type Telegram struct {
	ChatID    int64  `yaml:"chatID"`
	BotToken  string `yaml:"botToken"`
	ParseMode string `yaml:"parseMode"`
}

type Wechat struct {
	Webhook string `yaml:"webhook"`
}

type Notifications struct {
	Mail     `yaml:"mail"`
	Telegram `yaml:"telegram"`
	Wechat   `yaml:"wechat"`
}

type FullConfig struct {
	Notifications `yaml:"notifications"`
	DDNS          `yaml:"ddns"`
	Modem         `yaml:"modem"`
}

func LoadConfig(filePath string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(f, &Config)
}
