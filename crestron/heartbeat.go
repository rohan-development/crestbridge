package crestron

import (
	"crestbridge/mqttbridge"
	"encoding/xml"
	"net"
	"strconv"
	"time"
)

func startHeartbeat(conn net.Conn) {
	go func() {
		for {
			heartbeat := Cresnet{
				Control: &Control{
					Comm: Comm{
						HeartbeatRequest: &HeartbeatRequest{},
					},
				},
			}
			xml, err := xml.Marshal(heartbeat)
			check_err(err)
			conn.Write(xml)
			time.Sleep(5 * time.Second)
			for id, state := range idToState {
				//time.Sleep(100 * time.Millisecond)
				if idToKind[id] == "analog" {
					mqttbridge.PublishState(
						idToName[id],
						"analog",
						strconv.Itoa(state),
					)
				} //else {
				stateNew := "ON"
				if state == 0 {
					stateNew = "OFF"
				}
				mqttbridge.PublishState(
					idToName[id],
					"digital",
					stateNew,
				)
				//}
			}
			//fmt.Print(idToState)
		}
	}()
}
