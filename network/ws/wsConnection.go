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

func (connection *WebsocketConnection) ParseWrapper(data []byte) (*WebsocketMessageWrapper, error) {
	wrapper := new(WebsocketMessageWrapper)

	err := json.Unmarshal(data, wrapper)

	if err != nil {
		log.Println("error: ", err)
		return nil, errors.New("Can't parse message")
	}

	return wrapper, nil
}

/*WARNING! это можно вызывать только из того треда, который отправляет собщения в канал*/
func (connection *WebsocketConnection) Close(code int, message string) {
	log.Println("WSCon: Closing websocket connection")
	//	if (len(message) > 0) {
	//		connection.responseChannel <-
	//	}

	close(connection.OutgoingChannel)
	connection.ws.Close()
}

func (connection *WebsocketConnection) Write(msg interface{}) {
	connection.OutgoingChannel <- msg
}

func (connection *WebsocketConnection) StartReading() {
	go func() {
		defer func() {
			// нельзя закрывать самому себя - это может сделать только тот, кто пишет в канал
			//connection.Close(0, "unimplemented")

			connection.ws.Close()
			connection.ClosingChannel <- connection.Id
		}()

		for {
			if msgType, buffer, err := connection.ws.ReadMessage(); err == nil {
				switch msgType {
				case websocket.CloseMessage:
					log.Println("Close message received")
				case websocket.PingMessage:
					log.Println("Ping message received")
				case websocket.TextMessage:
					if wrapper, err := connection.ParseWrapper(buffer); err != nil {
						log.Println("Error parsing wrapper", err)
					} else {
						if msg, err := wrapper.GetMessage(); err == nil {
							connection.IncomingChannel <- network.UserMessage{connection.Id, msg}
						} else {
							log.Println(err)
						}

						/*
						case *network.SyncTimeMessage:
						// log.Println("Sync time message", time.Now().UnixNano()/int64(time.Millisecond), int64(time.Now().UnixNano()/int64(time.Millisecond)))
						connection.OutgoingChannel <- network.SyncTimeMessage{Time: int64(time.Now().UnixNano() / int64(time.Millisecond))}
						*/

					}
				case websocket.BinaryMessage:
					log.Println("Binary message!")
				}
			} else {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					// всё штатно закрылось
				} else {
					fmt.Println("Error reading from connection", err.Error())
				}
				break
			}
		}

		log.Println("Reading finished for", connection.Id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		//log.Println("Writing started ", connection.id)

		for message := range connection.OutgoingChannel {
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

func NewWebsocketConnection(id uint64, ws *websocket.Conn, ch network.UserMessagesChannel, closingChannel chan uint64) network.Connection {
	connection := &WebsocketConnection{
		ws: ws,
		BasicConnection: &network.BasicConnection{
			Id: id,
			OutgoingChannel: make (chan interface{}),
			IncomingChannel: make (chan interface{}),
			ClosingChannel: closingChannel,
		},
	}

	go connection.StartReading()
	go connection.StartWriting()

	log.Println("Created new ws connection", connection)
	return connection
}
