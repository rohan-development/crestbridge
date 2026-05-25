package crestron

import (
	"crestbridge/mqttbridge"
	"encoding/xml"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const THRESHOLD = 5000
const OVERSHOOT = 7000

func pressJoin(join, handle int, value, repeating bool) {
	msg := Cresnet{
		Data: &Data{
			Handle: handle,
			Slot:   0,
			Bool: &Bool{
				ID:        join,
				Value:     value,
				Repeating: repeating,
			},
		},
	}
	xml, err := xml.Marshal(msg)
	check_err(err)
	states.Write(xml)
}

func stepUp(id int) {
	pressJoin(idToUp[id], 1, true, true)
}

func analogOn(id int) {
	pressJoin(idToUp[id], 1, true, false)
	time.Sleep(3 * time.Second)
	pressJoin(idToUp[id], 1, false, false)
}

func analogOff(id int) {
	pressJoin(idToDown[id], 1, true, false)
	time.Sleep(3 * time.Second)
	pressJoin(idToDown[id], 1, false, false)
}

func stepDown(id int) {
	pressJoin(idToDown[id], 1, true, true)
}

func stepAnalog(id, value int) {
	if value > GetState(id)+THRESHOLD {
		stepUp(id)
	} else if value < GetState(id)-THRESHOLD {
		stepDown(id)
	}
}

func MQTTSubscribe(name, kind string) {
	topic := name + "/set"
	//fmt.Println(room + "/" + entity + "/set")
	mqttbridge.GetClient().Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {

		payload := string(m.Payload())
		id := nameToID[name]
		if kind == "analog" {
			switch payload {
			case "ON":
				if GetState(id) < 1000 {
					analogOn(id)
				}
			case "OFF":
				analogOff(id)
			default:
				value, err := strconv.Atoi(payload)
				if err == nil {
					stepAnalog(id, value)
				}

			}
		} else if kind == "digital" {
			if len(idToCtrlID[id]) == 0 {
				switch payload {
				case "ON":
					pressJoin(id, 1, true, false)
					time.Sleep(time.Second)
					pressJoin(id, 1, false, false)
				case "OFF":
					pressJoin(id, 1, true, false)
					time.Sleep(time.Second)
					pressJoin(id, 1, false, false)

				}
			} else {
				switch payload {
				case "ON":
					for i := 0; i < len(idToCtrlID[id]); i++ {
						pressJoin(idToCtrlID[id][i], 3, true, false)
						time.Sleep(300 * time.Millisecond)
						pressJoin(idToCtrlID[id][i], 3, false, false)
					}
				case "OFF":
					for i := 0; i < len(idToCtrlIDOff[id]); i++ {
						pressJoin(idToCtrlIDOff[id][i], 3, true, false)
						time.Sleep(300 * time.Millisecond)
						pressJoin(idToCtrlIDOff[id][i], 3, false, false)
					}

				}
			}

		}

	})
}
