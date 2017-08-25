package worlds

import (
	"github.com/Hucaru/Valhalla/common/constants"
	"github.com/Hucaru/Valhalla/common/packet"
)

// Message system between world manager and anything else
type Message struct {
	Opcode  int
	Message interface{}
}

var activeSessions map[string]chan Message

// MainLoop - Main loop of world manager
func MainLoop() {
	activeSessions = make(map[string]chan Message)

	for {
		for sessionID, client := range activeSessions {
			message := <-client

			switch message.Opcode {
			case CLIENT_NOT_ACTIVE:
				delete(activeSessions, sessionID)
			case WORLD_LIST:
				generateWorldList(message.Message.(chan [][]byte))
			}
		}
	}
}

// NewClient - Add a new client to the active sessions
func NewClient(client chan Message, sessionID string) {
	activeSessions[sessionID] = client
}

func generateWorldList(result chan [][]byte) {
	worlds := make([][]byte, 0)
	for j := 0; j < 10; j++ {
		pac := packet.NewPacket()
		pac.WriteByte(constants.LOGIN_SEND_WORLD_LIST)
		pac.WriteByte(byte(j))
		pac.WriteString("scania") // World name -
		pac.WriteByte(byte(2))    // Ribbon on world - 0 = normal, 1 = event, 2 = new, 3 = hot
		pac.WriteString("event description")
		pac.WriteByte(byte(20)) // number of channels
		pac.WriteByte(3)

		for i := 0; i < 20; i++ {
			pac.WriteString("scania-" + string(i+1)) // channel name
			pac.WriteInt32(9001)
			pac.WriteByte(byte(j))
			pac.WriteByte(byte(i))
			pac.WriteByte(0)
		}
		worlds = append(worlds, pac)
	}
	result <- worlds
}