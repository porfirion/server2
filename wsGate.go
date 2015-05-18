package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net"
	"net/http"
)

type WebSocketGate struct {
	addr                *net.TCPAddr
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
	listener, err := net.ListenTCP("tcp4", gate.addr)

	if err != nil {
		log.Fatal("error creating listener")
	}

	log.Println("Listening http:", gate.addr)

	server.Serve(listener)
}

/**
 * Обработчик входящих подключений по websocket
 * @param  {[type]} gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request [description]
 * @return {[type]} [description]
 */
func (gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request) {
	webSocket, err := upgrader.Upgrade(rw, request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("WSGate: Not a websocket handshake")
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println("WSGate: Unknown error", err)
		return
	}

	conn := NewWebsocketConnection(webSocket)
	log.Println("WSGate: new websocket connection")

	gate.incomingConnections <- conn
}

/**
 * Отдаёт главную (и единственную) страницу
 * @param  {[type]} gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request [description]
 * @return {[type]} [description]
 */
func (gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request) {
	indexTempl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct{}{}
	indexTempl.Execute(rw, data)
}

/**
 * Отвечает за отдачу статики
 * @param  {[type]} gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request [description]
 * @return {[type]} [description]
 */
func (gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request) {
	http.ServeFile(rw, request, request.URL.Path[1:])
}
