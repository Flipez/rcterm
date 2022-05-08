package ws

import (
	"encoding/json"
)

type SubEvent struct {
	Message    string `json:"msg"`
	Collection string `json:"collection"`
}

type SubMessage struct {
	Message string `json:"msg"`
	Id      string `json:"id"`
	Name    string `json:"name"`
}

type StreamRoomMessage struct {
	SubMessage
	Params []string `json:"params"`
}

type StreamRoomMessageEvent struct {
	SubEvent
	Fields StreamRoomMessageEventFields `json:"fields"`
}

type StreamRoomMessageEventFields struct {
	EventName string    `json:"eventName"`
	Args      []Message `json:"args"`
}

func (c *Connection) SubscribeOwn() error {
	id := c.NewUUID("stream-room-messages")
	msg, _ := json.Marshal(StreamRoomMessage{SubMessage: SubMessage{Message: "sub", Name: "stream-room-messages", Id: id.String()}, Params: []string{"__my_messages__"}})
	err := c.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) ParseStreamRoomMessage(message []byte) {
	var msg StreamRoomMessageEvent
	json.Unmarshal(message, &msg)
	c.MessageSubChannel <- msg.Fields.Args[0]
}
