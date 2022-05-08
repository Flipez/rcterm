package ws

import (
	"encoding/json"
)

type GetMessage struct {
	Message string `json:"msg"`
	Method  string `json:"method"`
	Id      string `json:"id"`
}

type RoomsGetMessage struct {
	GetMessage
}

type RoomGetResultMessage struct {
	InMessage
	Rooms []ChatRoom `json:"result"`
}

func (c *Connection) GetRooms() error {
	id := c.NewUUID("rooms/get")
	msg, _ := json.Marshal(RoomsGetMessage{GetMessage: GetMessage{Message: "method", Method: "rooms/get", Id: id.String()}})
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
