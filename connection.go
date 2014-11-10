package main

import (
	"github.com/gorilla/websocket"
	"net"
)

type WebsocketConnection struct {
	conn *websocket.Conn
}

func (conn *WebsocketConnection) StartReading(ch MessagesChannel) {

}

func (conn *WebsocketConnection) WriteMessage(msg Message) {}

func NewWebsocketConnection(conn *websocket.Conn) Connection {
	connection := &WebsocketConnection{conn}
	return connection
}

type TcpConnection struct {
	conn net.Conn
}

func (conn *TcpConnection) StartReading(ch MessagesChannel) {

}
func (conn *TcpConnection) WriteMessage(msg Message) {}

func NewTcpConnection(conn net.Conn) Connection {
	connection := &TcpConnection{conn}
	return connection
}

type Connection interface {
	StartReading(ch MessagesChannel)
	WriteMessage(msg Message)
}
