package crestron

import (
	"crestbridge/mqttbridge"
	"strconv"
)

func Listen() {
	for {
		changed, idChanged := updateStates(false)
		if changed {
			//fmt.Println(idToState)
			mqttbridge.PublishState(
				idToName[idChanged],
				"analog",
				strconv.Itoa(GetState(idChanged)),
			)
			state := "ON"
			if GetState(idChanged) == 0 {
				state = "OFF"
			}
			mqttbridge.PublishState(
				idToName[idChanged],
				"digital",
				state,
			)

		}
	}
}
