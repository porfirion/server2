package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

type WebSocketGate struct {
	incomingConnections chan *Connection
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (gate *WebSocketGate) Start(incomingConnections chan *Connection) {
	gate.incomingConnections = incomingConnections
	http.HandleFunc("/", gate.indexHandler)
	http.HandleFunc("/assets/", gate.assetsHandler)
	http.HandleFunc("/ws", gate.wsHandler)

	log.Println("ADDR: " + HTTP_HOST + ":" + HTTP_PORT)
	if err := http.ListenAndServe(HTTP_HOST+":"+HTTP_PORT, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

func (gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request) {
	webSocket, err := upgrader.Upgrade(rw, request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	conn := NewWebsocketConnection(webSocket)
	gate.incomingConnections <- &conn
}

func (gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request) {
	indexTempl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct{}{}
	indexTempl.Execute(rw, data)
}

func (gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request) {
	http.ServeFile(rw, request, request.URL.Path[1:])
}
