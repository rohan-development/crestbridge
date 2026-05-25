package crestron

import (
	"crestbridge/config"
	"encoding/xml"
	"net"
	"sync"
)

var address string
var login Cresnet
var dev []config.GenericDevice

var states net.Conn
var control net.Conn

// var idToState map[int]int
var (
	idToState = make(map[int]int)
	stateMu   sync.RWMutex
)
var idToName map[int]string
var nameToID map[string]int
var idToKind map[int]string
var idToRoom map[int]string
var idToStateString map[int]string
var idToDown map[int]int
var idToUp map[int]int
var idToCtrlID map[int][]int
var idToCtrlIDOff map[int][]int
var idToCheckBool map[int]bool
var buffer []byte

type Cresnet struct {
	XMLName xml.Name `xml:"cresnet"`

	Control *Control `xml:"control,omitempty"`
	Data    *Data    `xml:"data,omitempty"`
}

type Control struct {
	Comm Comm `xml:"comm"`
}

type Comm struct {
	ConnectRequest   *ConnectRequest   `xml:"connectRequest,omitempty"`
	HeartbeatRequest *HeartbeatRequest `xml:"heartbeatRequest,omitempty"`
}

type HeartbeatRequest struct{}

type ConnectRequest struct {
	Passcode string `xml:"passcode"`
	Mode     Mode   `xml:"mode"`
}

type Mode struct {
	IsUnicodeSupported bool `xml:"isUnicodeSupported,attr"`
}

type Data struct {
	Handle        int            `xml:"handle,attr"`
	Slot          int            `xml:"slot,attr"`
	I32           *I32           `xml:"i32,omitempty"`
	Bool          *Bool          `xml:"bool,omitempty"`
	UpdateRequest *UpdateRequest `xml:"updateRequest,omitempty"`
	UpdateCommand *UpdateCommand `xml:"updateCommand,omitempty"`
	String        *String        `xml:"string,omitempty"`
}

type String struct {
	ID    int    `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type UpdateCommand struct {
	EndOfUpdate *EndOfUpdate `xml:"endOfUpdate"`
}

type EndOfUpdate struct{}

type I32 struct {
	ID    int `xml:"id,attr"`
	Value int `xml:",chardata"`
}

type Bool struct {
	ID        int  `xml:"id,attr"`
	Value     bool `xml:"value,attr"`
	Repeating bool `xml:"repeating,attr"`
}

type UpdateRequest struct{}

func Loadconfig(conf *config.Config) {
	address = conf.CrestronIP + ":" + conf.Port
	login = Cresnet{
		Control: &Control{
			Comm: Comm{
				ConnectRequest: &ConnectRequest{
					Passcode: conf.Password,
					Mode: Mode{
						IsUnicodeSupported: true,
					},
				},
			},
		},
	}
	dev = conf.Devices
}
