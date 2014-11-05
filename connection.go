package main

import (
	"github.com/gorilla/websocket"
	"net"
)

type WebsocketConnection struct {
	conn *websocket.Conn
}

func (conn *WebsocketConnection) StartReading(ch chan Message) {

}

func (conn *WebsocketConnection) WriteMessage(msg Message) {}

func NewWebsocketConnection(conn *websocket.Conn) Connection {
	return &WebsocketConnection{conn: conn}
}

type TcpConnection struct {
}

func (conn *TcpConnection) StartReading(ch chan Message) {

}
func (conn *TcpConnection) WriteMessage(msg Message) {}

func NewTcpConnection(conn net.Conn) Connection {
	return &TcpConnection{}
}

type Connection interface {
	StartReading(ch chan Message)
	WriteMessage(msg Message)
}
