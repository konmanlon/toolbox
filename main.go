package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toolbox/config"
	"toolbox/ddns"
	"toolbox/modem"
)

func main() {
	configPath := flag.String("f", "config.yaml", "config file path")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	c := config.Config

	if c.DDNS.Enable {
		go func() {
			ddns.IsInit()

			tasker := time.NewTicker(time.Second * time.Duration(c.DDNS.ScheduledTask))

			for {
				<-tasker.C
				msg, err := ddns.RunDDNS()
				if err != nil {
					log.Println(err)
				}
				if msg != nil {
					log.Println(*msg)
				}
			}
		}()
	}

	if c.Modem.Enable {
		go func() {
			port := c.Modem.SerialPort

			dev := modem.DeviceAir72x{
				CommandPort: port,
				NotifyPort:  port,
			}

			if err := dev.InitDevice(); err != nil {
				log.Println(err)
				return
			}

			dev.Watch()
		}()
	}

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT)

	<-chSignal
}
