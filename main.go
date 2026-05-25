package main

import (
	"crestbridge/config"
	"crestbridge/crestron"
	"crestbridge/mqttbridge"
)

func main() {
	config := config.Load()
	crestron.Loadconfig(config)
	mqttbridge.Connect("tcp://localhost:1883")
	crestron.Initialize()
	go crestron.Listen()

	select {}

}
