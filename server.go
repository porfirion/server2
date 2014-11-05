package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

const (
	HTTP_HOST string = ""
	HTTP_PORT string = "8080"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var incomingConnections chan *Connection = make(chan *Connection, 10)
var incomingMessages chan Message = make(chan Message, 100)

func main() {
	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic := new(Logic)
	logic.incomingMessages = incomingMessages
	go logic.run()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Println("ADDR: " + HTTP_HOST + ":" + HTTP_PORT)
	if err := http.ListenAndServe(HTTP_HOST+":"+HTTP_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	log.Println("running")

	for {

	}

	log.Println("exit")
}

func wsHandler(rw http.ResponseWriter, request *http.Request) {
	webSocket, err := upgrader.Upgrade(rw, request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	conn := NewWebsocketConnection(webSocket)
	incomingConnections <- &conn
}

func indexHandler(rw http.ResponseWriter, request *http.Request) {
	indexTempl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct{}{}
	indexTempl.Execute(rw, data)
}

func assetsHandler(rw http.ResponseWriter, request *http.Request) {
	http.ServeFile(rw, request, request.URL.Path[1:])
}
