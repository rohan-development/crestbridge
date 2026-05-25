package crestron

import (
	"encoding/xml"
	"net"
)

func Initialize() {
	mapDevices()
	control = auth()
	//states = auth()
	states = control
	startHeartbeat(control)
	startHeartbeat(states)
	triggerStates()
	updateStates(true)
	// Hacky fix. Could cause problems.
	for id, kind := range idToKind {
		if kind == "analog" && GetState(id) == 1 {
			SetState(id, 0)
		}
	}
	//fmt.Println(GetState(54))
	//return conn
}

func auth() net.Conn {
	conn, err := net.Dial("tcp", address)
	check_err(err)
	xml, err := xml.Marshal(login)
	check_err(err)
	conn.Write(xml)
	return conn
}

func triggerStates() {
	msg := Cresnet{
		Data: &Data{
			UpdateRequest: &UpdateRequest{},
		},
	}
	xml, err := xml.Marshal(msg)
	check_err(err)
	states.Write(xml)
}

func mapDevices() { //id to state, name to id, id to name
	//idToState = make(map[int]int)
	nameToID = make(map[string]int)
	idToName = make(map[int]string)
	idToRoom = make(map[int]string)
	idToStateString = make(map[int]string)
	idToDown = make(map[int]int)
	idToUp = make(map[int]int)
	idToCtrlID = make(map[int][]int)
	idToCtrlIDOff = make(map[int][]int)
	idToKind = make(map[int]string)
	idToCheckBool = make(map[int]bool)
	// Generic Devices
	for i := 0; i < len(dev); i++ {
		// stateMu.Lock()
		// idToState[dev[i].ID] = 0
		// stateMu.Unlock()
		SetState(dev[i].ID, 0)
		MQTTSubscribe(dev[i].Room+"/"+dev[i].Name, dev[i].Type)
		nameToID[dev[i].Room+"/"+dev[i].Name] = dev[i].ID
		idToName[dev[i].ID] = dev[i].Room + "/" + dev[i].Name
		idToRoom[dev[i].ID] = dev[i].Room
		idToDown[dev[i].ID] = dev[i].Down
		idToCheckBool[dev[i].ID] = dev[i].CheckBool
		idToUp[dev[i].ID] = dev[i].Up
		idToCtrlID[dev[i].ID] = dev[i].CTRLON
		idToCtrlIDOff[dev[i].ID] = dev[i].CTRLOFF
		idToKind[dev[i].ID] = dev[i].Type
		if dev[i].StateString != "" {
			idToStateString[dev[i].ID] = dev[i].StateString
		}
	}
}
