package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
)

type WebsocketMessageWrapper struct {
	MessageType int
	Data        string
}

type WebsocketConnection struct {
	ws *websocket.Conn
}

func ParseMessage(data []byte) Message {
	wrapper := new(WebsocketMessageWrapper)

	err := json.Unmarshal(data, wrapper)

	if err != nil {
		log.Println("error: ", err)
		return nil
	}

	log.Println("Parsed")
	log.Println(wrapper)

	return nil
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
	} else if msg := ParseMessage(data); msg == nil {
		log.Println("error parsing message", err)
		return nil, errors.New("error parsing message")
	} else {
		return msg, nil
	}

	return new(DataMessage), nil
}

func (connection *WebsocketConnection) StartReading(ch MessagesChannel) {
	defer connection.Close()

	for {
		if msg, err := connection.ReadMessage(); err != nil {
			log.Println("Error reading message")
			break
		} else {
			ch <- msg
		}
	}
}

func (connection *WebsocketConnection) WriteMessage(msg Message) {}

func (connection *WebsocketConnection) Close() {
	log.Println("Closing websocket connection")
	connection.ws.Close()
}
func (connection *WebsocketConnection) GetAuth() *Player {
	msg, err := connection.ReadMessage()
	if err != nil {
		log.Println("couldn't read logic from websocket")
		return nil
	} else if auth, ok := msg.(AuthMessage); ok {
		player := &Player{auth.name, connection}
		return player
	} else {
		log.Println("Not an auth mesage")
		return nil
	}
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	connection := &WebsocketConnection{ws}
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
