package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"log"
)

type WebsocketMessageWrapper struct {
	MessageType int
	Data        []byte
}

func (wrapper *WebsocketMessageWrapper) GetMessage() (msg Message, err error) {
	switch wrapper.MessageType {
	case 1:
		var res AuthMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
	case 1000:
		var res DataMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
	case 1001:
		var res TextMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
	default:
		log.Println("Unknown message type: ", wrapper.MessageType)
	}

	return
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

func (connection *WebsocketConnection) ReadMessage() (Message, error) {
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

func (connection *WebsocketConnection) StartReading(ch MessagesChannel) {
	go func() {
		defer connection.Close()

		for {
			msg, err := connection.ReadMessage()
			if err != nil {
				log.Println("Error reading message")
				break
			}
			if msg != nil {
				ch <- msg
			}
		}

		log.Println("Reading finished for", connection.id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		log.Println("Writing started ", connection.id)

		for message := range connection.GetResponseChannel() {
			log.Println("For ", connection.id, " message ", message)
			connection.ws.WriteJSON(message)
		}

		log.Println("Writing finished for", connection.id)
	}()
}

func (connection *WebsocketConnection) Close() {
	log.Println("Closing websocket connection")
	close(connection.responseChannel)
	connection.ws.Close()
	connection.closingChannel <- connection.id
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	connection := &WebsocketConnection{ws: ws, BasicConnection: &BasicConnection{}}
	connection.StartWriting()
	log.Println("Created new ws connection", connection)
	return connection
}
