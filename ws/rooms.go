package ws

import (
	"encoding/json"
)

type MethodMessage struct {
	Message string `json:"msg"`
	Method  string `json:"method"`
	Id      string `json:"id"`
}

type RoomsGetMessage struct {
	MethodMessage
}

type RoomGetResultMessage struct {
	InMessage
	Rooms []ChatRoom `json:"result"`
}

type RoomOpenMessage struct {
	MethodMessage
	Params []string `json:"params"`
}

func (c *Connection) GetRooms() error {
	id := c.NewUUID("rooms/get")
	msg, _ := json.Marshal(RoomsGetMessage{MethodMessage: MethodMessage{Message: "method", Method: "rooms/get", Id: id.String()}})
	err := c.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) ParseGetRooms(message []byte) {
	var result RoomGetResultMessage
	json.Unmarshal(message, &result)
	for _, room := range result.Rooms {
		//if room.Type == "c" {
		c.RoomChannel <- room
		//}
	}
}

func (c *Connection) OpenRoom(roomId string) error {
	id := c.NewUUID("openRoom")
	msg, _ := json.Marshal(RoomOpenMessage{MethodMessage: MethodMessage{Message: "method", Method: "openRoom", Id: id.String()}, Params: []string{roomId}})
	err := c.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
