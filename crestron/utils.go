package crestron

import (
	"bytes"
	"encoding/xml"
	"log"
	"net"
	"time"
)

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findStart(b []byte) int {
	return bytes.Index(b, []byte("<cresnet>"))
}

func findEnd(b []byte) int {
	return bytes.Index(b, []byte("</cresnet>"))
}

func handleChunk(chunk []byte) []Cresnet {
	buffer = append(buffer, chunk...)
	var cresnet []Cresnet
	for {
		start := bytes.Index(buffer, []byte("<cresnet>"))
		end := bytes.Index(buffer, []byte("</cresnet>"))

		if start == -1 || end == -1 || end <= start {
			return cresnet
		}

		end += len("</cresnet>")

		msg := buffer[start:end]
		buffer = buffer[end:]

		var c Cresnet
		err := xml.Unmarshal(msg, &c)
		if err != nil {
			continue
		}
		// if err := xml.Unmarshal(msg, &c); err == nil {
		// 	fmt.Println("parsed cresnet message")
		// }
		cresnet = append(cresnet, c)
	}
}

func updateStates(first bool) (bool, int) {
	changed := false
	idChanged := 0
	buf := make([]byte, 4096)

	for {
		// 1. set deadline BEFORE read
		states.SetReadDeadline(time.Now().Add(50 * time.Millisecond))

		n, err := states.Read(buf)
		if err != nil && !first {
			// timeout = end of batch (only if not first sync)
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				states.SetReadDeadline(time.Time{})
				return changed, idChanged
			}

			log.Fatal(err)
		}
		cresnet := handleChunk(buf[:n])

		for _, msg := range cresnet {
			//fmt.Println("s")
			if msg.Data != nil &&
				msg.Data.UpdateCommand != nil &&
				msg.Data.UpdateCommand.EndOfUpdate != nil {

				states.SetReadDeadline(time.Time{})
				return changed, idChanged
			}
			if msg.Data != nil && msg.Data.I32 != nil {
				i32 := *msg.Data.I32
				stateMu.RLock()
				old, exists := idToState[i32.ID]
				stateMu.RUnlock()
				if exists {
					//idToState[i32.ID] = i32.Value
					SetState(i32.ID, i32.Value)
					if old != i32.Value {
						changed = true
						idChanged = i32.ID
					}
				}
			}
			if msg.Data != nil && msg.Data.Bool != nil {
				state := msg.Data.Bool.Value
				stateMu.RLock()
				old, exists := idToState[msg.Data.Bool.ID]
				stateMu.RUnlock()
				if exists && idToCheckBool[msg.Data.Bool.ID] {
					if state {
						SetState(msg.Data.Bool.ID, 1)
					} else {
						SetState(msg.Data.Bool.ID, 0)
					}
					if GetState(msg.Data.Bool.ID) != old {
						changed = true
						idChanged = msg.Data.Bool.ID
					}
				}

			}
			if msg.Data != nil && msg.Data.String != nil {
				str := msg.Data.String
				stateMu.RLock()
				old, exists := idToState[str.ID]
				stateMu.RUnlock()
				if exists {
					if str.Value == idToStateString[str.ID] {
						//idToState[str.ID] = 1
						SetState(str.ID, 1)
					} else {
						// idToState[str.ID] = 0
						SetState(str.ID, 0)
					}

					if GetState(str.ID) != old {
						changed = true
						idChanged = str.ID
					}
				}
			}
		}
	}
}

func SetState(id int, value int) {
	stateMu.Lock()
	defer stateMu.Unlock()
	idToState[id] = value
}

func GetState(id int) int {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return idToState[id]
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
