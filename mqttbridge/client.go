package mqttbridge

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

func Connect(broker string) {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("crestbridge")

	client = mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		panic(err)
	}

	//fmt.Println("MQTT connected")
}

func PublishState(name, kind string, value any) {

	topic := fmt.Sprintf("%s/%s/value", name, kind)
	token := client.Publish(topic, 0, true, value)
	token.Wait()
}

func GetClient() mqtt.Client {
	return client
}
