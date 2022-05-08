package ws

import "encoding/json"

type InMessage struct {
	Message string `json:"msg"`
	Id      string `json:"id"`
}

type ConnectMessage struct {
	Message string   `json:"msg"`
	Version string   `json:"version"`
	Support []string `json:"support"`
}

type ResumeParam struct {
	Resume string `json:"resume"`
}

type LoginMessage struct {
	Message string        `json:"msg"`
	Method  string        `json:"method"`
	Id      string        `json:"id"`
	Params  []ResumeParam `json:"params"`
}

type User struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
}

type ChatRoom struct {
	Id      string   `json:"_id"`
	Type    string   `json:"t"`
	Name    string   `json:"name"`
	Creator User     `json:"u"`
	Topic   string   `json:"topic"`
	Muted   []string `json:"muted"`
}

type Message struct {
	Id        string `json:"id"`
	Rid       string `json:"rid"`
	Message   string `json:"msg"`
	Timestamp string `json:"ts"`
	Sender    User   `json:"u"`
}

type Date int

func (d Date) MarshalJSON() ([]byte, error) {
	result := make(map[string]int)

	result["$date"] = int(d)

	return json.Marshal(result)
}
