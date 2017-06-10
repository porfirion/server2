package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	// "reflect"
	"time"
	"github.com/porfirion/server2/network"
)

type WebsocketMessageWrapper struct {
	MessageType int    `json:"type"`
	Data        string `json:"data"`
}

func (wrapper *WebsocketMessageWrapper) GetMessage() (msg interface{}, err error) {
	res := network.GetValueByTypeId(wrapper.MessageType)

	// fmt.Printf("Unmarshalling into %#v, (type %v)\n", res, reflect.TypeOf(res))
	// fmt.Printf("Message body: %#v\n", wrapper.Data)

	err = json.Unmarshal([]byte(wrapper.Data), res)

	//fmt.Printf("unmarshalled: %#v error: %#v\n", res, err)

	return res, err
}

type WebsocketConnection struct {
	*network.BasicConnection
	ws *websocket.Conn
}

func (connection *WebsocketConnection) ParseMessage(data []byte) (*WebsocketMessageWrapper, error) {
	wrapper := new(WebsocketMessageWrapper)

	err := json.Unmarshal(data, wrapper)

	if err != nil {
		log.Println("error: ", err)
		return nil, errors.New("Can't parse message")
	}

	return wrapper, nil
}

func (connection *WebsocketConnection) ReadMessage() (interface{}, error) {
	msgType, data, err := connection.ws.ReadMessage()

	if err != nil {
		if err == io.EOF {
			log.Println("WSCon: EOF")
		}
		log.Println("Error reading websocket", err)
		return nil, io.EOF
	} else if msgType == websocket.CloseMessage {
		log.Println("Close message received")
		return nil, io.EOF
	} else if msgType == websocket.PingMessage {
		log.Println("Ping message received")
		return nil, nil
	} else if wrapper, err := connection.ParseMessage(data); err != nil {
		log.Println("Error parsing message", err)
		return nil, err
	} else {
		return wrapper.GetMessage()
	}
}

func (connection *WebsocketConnection) StartReading(ch network.UserMessagesChannel) {
	go func() {
		defer connection.Close(0, "unimplemented")

		for {
			msg, err := connection.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err)
				break
			}
			if msg != nil {
				switch msg.(type) {
				case *network.SyncTimeMessage:
					// log.Println("Sync time message", time.Now().UnixNano()/int64(time.Millisecond), int64(time.Now().UnixNano()/int64(time.Millisecond)))
					connection.ResponseChannel <- network.SyncTimeMessage{Time: int64(time.Now().UnixNano() / int64(time.Millisecond))}
				default:
					ch <- network.UserMessage{connection.Id, msg}
				}
			}
		}

		log.Println("Reading finished for", connection.Id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		//log.Println("Writing started ", connection.id)

		for message := range connection.GetResponseChannel() {
			//log.Println(fmt.Sprintf("WsCon. Sending message %T for %d", message, connection.id))
			bytes, err := json.Marshal(message)
			if err == nil {
				connection.ws.WriteJSON(WebsocketMessageWrapper{network.GetMessageTypeId(message), string(bytes)})
			} else {
				log.Println("Error serializing message", message, err)
			}

		}

		log.Println("Writing finished for", connection.Id)
	}()
}

func (connection *WebsocketConnection) Close(code int, message string) {
	log.Println("WSCon: Closing websocket connection")
	//	if (len(message) > 0) {
	//		connection.responseChannel <-
	//	}

	close(connection.ResponseChannel)
	connection.ws.Close()
	connection.ClosingChannel <- connection.Id
}

func (connection *WebsocketConnection) GetAuth() (*network.AuthMessage, error) {
	msg, err := connection.ReadMessage()
	if err == nil {
		if auth, ok := msg.(*network.AuthMessage); ok {
			return auth, nil
		} else {
			fmt.Printf("Converted message %#v\n", msg)
			return nil, errors.New("Wrong message type")
		}
	} else {
		return nil, err
	}
}

func NewWebsocketConnection(ws *websocket.Conn) network.Connection {
	connection := &WebsocketConnection{ws: ws, BasicConnection: &network.BasicConnection{}}
	connection.StartWriting()
	log.Println("Created new ws connection", connection)
	return connection
}
