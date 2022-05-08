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

type Connection struct {
	config            config.RctermConfig
	con               *websocket.Conn
	pendingActions    map[uuid.UUID]string
	RoomChannel       chan ChatRoom
	MessagesChannel   chan []Message
	MessageSubChannel chan Message
}

func (c *Connection) Connect() {
	c.config = config.ReadConfig()
	c.pendingActions = make(map[uuid.UUID]string)
	var err error

	u := url.URL{Scheme: "wss", Host: c.config.URL, Path: "/websocket"}

	c.con, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	msg, _ := json.Marshal(ConnectMessage{Message: "connect", Version: "1", Support: []string{"1"}})
	c.Send(msg)

	c.Listen()
	c.SubscribeOwn()
}

func (c *Connection) Send(message []byte) error {
	log.Println(string(message))
	//c.LogsChannel <- string(message)
	err := c.con.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write:", err)
		return err
	}

	return nil
}

func (c *Connection) Listen() {
	go func() {

		f, err := os.OpenFile("messages.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)

		for {
			_, message, err := c.con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Println(string(message))
			//c.LogsChannel <- string(message)
			var incomingMessage InMessage
			json.Unmarshal(message, &incomingMessage)
			switch incomingMessage.Message {
			case "connected":
				msg, _ := json.Marshal(LoginMessage{Message: "method", Method: "login", Id: uuid.New().String(), Params: []ResumeParam{{Resume: c.config.Token}}})
				c.Send(msg)

				time.Sleep(2 * time.Second)
				err := c.GetRooms()
				if err != nil {
					panic(err)
				}
			case "ping":
				msg, _ := json.Marshal(InMessage{Message: "pong"})
				c.Send(msg)
			case "result":
				actionId := uuid.Must(uuid.Parse(incomingMessage.Id))
				resultType := c.pendingActions[actionId]

				switch resultType {
				case "rooms/get":
					c.ParseGetRooms(message)
				case "loadHistory":
					c.ParseGetHistory(message)
				}
			case "changed":
				var subMessage SubEvent
				json.Unmarshal(message, &subMessage)

				if subMessage.Collection == "stream-room-messages" {
					c.ParseStreamRoomMessage(message)
				}

			}
		}
	}()
}

func (c *Connection) NewUUID(t string) uuid.UUID {
	id := uuid.New()
	c.pendingActions[id] = t

	return id
}
