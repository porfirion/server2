package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net"
	"net/http"
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

type WebSocketGate struct {
	addr                string
	incomingConnections ConnectionsChannel
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (gate *WebSocketGate) Start() {

	http.HandleFunc("/", gate.indexHandler)
	http.HandleFunc("/assets/", gate.assetsHandler)
	http.HandleFunc("/ws", gate.wsHandler)

	server := &http.Server{}
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})

	if err != nil {
		log.Fatal("error creating listener")
	}

	go server.Serve(listener)
}

func (gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request) {
	log.Println("new websocket connection")
	webSocket, err := upgrader.Upgrade(rw, request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	conn := NewWebsocketConnection(webSocket)
	gate.incomingConnections <- conn
}

func (gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request) {
	indexTempl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct{}{}
	indexTempl.Execute(rw, data)
}

func (gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request) {
	http.ServeFile(rw, request, request.URL.Path[1:])
}
