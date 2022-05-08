package ws

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/flipez/rcterm/config"
)

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

type RoomsGetMessage struct {
	Message string `json:"msg"`
	Method  string `json:"method"`
	Id      string `json:"id"`
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

type RoomGetResultMessage struct {
	InMessage
	Rooms []ChatRoom `json:"result"`
}

type Message struct {
	Id        string `json:"id"`
	Rid       string `json:"rid"`
	Message   string `json:"msg"`
	Timestamp string `json:"ts"`
	Sender    User   `json:"u"`
}

type LoadHistoryResultMessage struct {
	InMessage
	Messages []Message `json:"result"`
}

func GetRooms(c *websocket.Conn) (uuid.UUID, error) {
	actionId := uuid.New()
	msg, _ := json.Marshal(RoomsGetMessage{Message: "method", Method: "rooms/get", Id: actionId.String()})
	err := c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return actionId, err
	}

	return actionId, nil
}

func GetHistory(c *websocket.Conn, roomId string) (uuid.UUID, error) {
	actionId := uuid.New()
	msg, _ := json.Marshal(RoomsGetMessage{Message: "method", Method: "loadHistory", Id: actionId.String()})
	err := c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return actionId, err
	}

	return actionId, nil
}

func ConnectServer(roomChannel chan ChatRoom) *websocket.Conn {
	rctermConfig := config.ReadConfig()

	pendingActions := map[uuid.UUID]string{}

	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	u := url.URL{Scheme: "wss", Host: rctermConfig.URL, Path: "/websocket"}
	//log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			//log.Printf("recv: %s", message)
			var incomingMessage InMessage
			json.Unmarshal(message, &incomingMessage)
			switch incomingMessage.Message {
			case "connected":
				msg, _ := json.Marshal(LoginMessage{Message: "method", Method: "login", Id: uuid.New().String(), Params: []ResumeParam{{Resume: rctermConfig.Token}}})
				messageOut <- string(msg)

				time.Sleep(2 * time.Second)
				actionId, err := GetRooms(c)
				if err != nil {
					panic(err)
				}
				pendingActions[actionId] = "rooms/get"
			case "ping":
				msg, _ := json.Marshal(InMessage{Message: "pong"})
				messageOut <- string(msg)
			case "result":
				actionId := uuid.Must(uuid.Parse(incomingMessage.Id))
				resultType := pendingActions[actionId]

				switch resultType {
				case "rooms/get":
					var result RoomGetResultMessage
					json.Unmarshal(message, &result)

					for _, room := range result.Rooms {
						if room.Type == "c" {
							roomChannel <- room
						}
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case m := <-messageOut:
				//log.Printf("Send Message %s", m)
				err := c.WriteMessage(websocket.TextMessage, []byte(m))
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-interrupt:
				//log.Println("interrupt")
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					//log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()

	//time.Sleep(3 * time.Second)
	msg, _ := json.Marshal(ConnectMessage{Message: "connect", Version: "1", Support: []string{"1"}})
	messageOut <- string(msg)

	return c
}
