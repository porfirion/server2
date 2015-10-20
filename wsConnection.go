package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	// "reflect"
	"time"
)

type WebsocketMessageWrapper struct {
	MessageType int    `json:"type"`
	Data        string `json:"data"`
}

func (wrapper *WebsocketMessageWrapper) GetMessage() (msg interface{}, err error) {
	res := GetValueByTypeId(wrapper.MessageType)

	//fmt.Printf("unmarshalling into %#v, (type %v)\n", res, reflect.TypeOf(res))

	err = json.Unmarshal([]byte(wrapper.Data), res)

	//fmt.Printf("unmarshalled: %#v error: %#v\n", res, err)

	return res, err
}

type WebsocketConnection struct {
	*BasicConnection
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
			log.Println("EOF")
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
		log.Println("error parsing message", err)
		return nil, err
	} else {
		return wrapper.GetMessage()
	}
}

func (connection *WebsocketConnection) StartReading(ch UserMessagesChannel) {
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
				case SyncTimeMessage:
					log.Println("Sync time message")
					log.Println(time.Now().UnixNano() / int64(time.Millisecond))
				default:
					ch <- UserMessage{connection.id, msg}
				}
			}
		}

		log.Println("Reading finished for", connection.id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		log.Println("Writing started ", connection.id)

		for message := range connection.GetResponseChannel() {
			log.Println(fmt.Sprintf("WsCon. Sending message %T for %d", message, connection.id))
			bytes, err := json.Marshal(message)
			if err == nil {
				connection.ws.WriteJSON(WebsocketMessageWrapper{GetMessageTypeId(message), string(bytes)})
			} else {
				log.Println("Error serializing message", message, err)
			}

		}

		log.Println("Writing finished for", connection.id)
	}()
}

func (connection *WebsocketConnection) Close(code int, message string) {
	log.Println("Closing websocket connection")
	//	if (len(message) > 0) {
	//		connection.responseChannel <-
	//	}

	close(connection.responseChannel)
	connection.ws.Close()
	connection.closingChannel <- connection.id
}

func (connection *WebsocketConnection) GetAuth() (*AuthMessage, error) {
	msg, err := connection.ReadMessage()
	if err == nil {
		if auth, ok := msg.(*AuthMessage); ok {
			return auth, nil
		} else {
			fmt.Printf("Converted message %#v\n", msg)
			return nil, errors.New("Wrong message type")
		}
	} else {
		return nil, err
	}
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	connection := &WebsocketConnection{ws: ws, BasicConnection: &BasicConnection{}}
	connection.StartWriting()
	log.Println("Created new ws connection", connection)
	return connection
}
