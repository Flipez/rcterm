package ws

import (
	"encoding/json"
	"time"
)

type LoadHistory struct {
	GetMessage
	Params LoadHistoryParams `json:"params"`
}

type LoadHistoryResultMessage struct {
	InMessage
	Result ResultMessages `json:"result"`
}

type ResultMessages struct {
	Messages []Message `json:"messages"`
}

type LoadHistoryParams struct {
	RoomId     string
	BeforeDate *Date
	Length     int
	CacheDate  Date
}

func (l LoadHistoryParams) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, 4)

	result[0] = l.RoomId
	result[1] = l.BeforeDate
	result[2] = l.Length
	result[3] = l.CacheDate

	return json.Marshal(result)
}

func (c *Connection) GetHistory(roomId string) error {
	id := c.NewUUID("loadHistory")

	params := LoadHistoryParams{
		RoomId:    roomId,
		Length:    50,
		CacheDate: Date(time.Now().Unix()),
	}

	msg, _ := json.Marshal(LoadHistory{GetMessage: GetMessage{Message: "method", Method: "loadHistory", Id: id.String()}, Params: params})
	err := c.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) ParseGetHistory(message []byte) {
	var result LoadHistoryResultMessage
	json.Unmarshal(message, &result)

	c.MessagesChannel <- result.Result.Messages

}
