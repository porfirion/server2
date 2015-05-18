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
	ws              *websocket.Conn
	responseChannel MessagesChannel
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
		connection.Close()
		return nil, io.EOF
	} else if msgType == websocket.CloseMessage {
		log.Println("Close message received")
		connection.Close()
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

		log.Println("Reading finished")
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		log.Println("Writing started")
		defer log.Println("Writing finished")

		for message := range connection.responseChannel {
			connection.ws.WriteJSON(message)
		}
	}()
}

func (connection *WebsocketConnection) GetResponseChannel() MessagesChannel {
	return connection.responseChannel
}

func (connection *WebsocketConnection) Close() {
	log.Println("Closing websocket connection")
	log.Println(connection.responseChannel)
	log.Println(len(connection.responseChannel))
	log.Println(cap(connection.responseChannel))
	//close(connection.responseChannel)
	connection.ws.Close()
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	connection := &WebsocketConnection{ws, make(MessagesChannel)}
	connection.StartWriting()
	return connection
}
