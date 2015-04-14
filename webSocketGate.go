package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
)

type WebsocketMessageWrapper struct {
	MessageType int
	Data        []byte
}

func (wrapper *WebsocketMessageWrapper) GetMessage() (msg Message) {
	// log.Println("Getting message from ", string(wrapper.Data))
	// fmt.Println(string(wrapper.Data))
	var err error
	switch wrapper.MessageType {
	case 1:
		var res AuthMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
		fmt.Println(err, res)
	case 1000:
		var res DataMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
		fmt.Println(err, res)
	case 1001:
		var res TextMessage
		err = json.Unmarshal(wrapper.Data, &res)
		msg = res
		fmt.Println(err, res)
	default:
		log.Println("Unknown message type: ", wrapper.MessageType)
	}

	fmt.Println("Parsed data message:", msg, err)
	if err == nil {
		return
	} else {
		return nil
	}
}

/**
 * Web Socket Connection
 */

type WebsocketConnection struct {
	ws *websocket.Conn
}

func ParseMessage(data []byte) (*WebsocketMessageWrapper, error) {
	log.Println("Parsing message: ", string(data), data)
	wrapper := new(WebsocketMessageWrapper)

	err := json.Unmarshal(data, wrapper)

	if err != nil {
		log.Println("error: ", err)
		return nil, errors.New("Can't parse message")
	}

	log.Println("Parsed wrapper", wrapper)

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
	} else if wrapper, err := ParseMessage(data); err != nil {
		log.Println("error parsing message", err)
		return nil, err
	} else {
		var msg Message = wrapper.GetMessage()
		fmt.Printf(format, msg)
		fmt.Println(msg)
		return msg, nil
	}
}

func (connection *WebsocketConnection) StartReading(ch MessagesChannel) {
	defer connection.Close()

	for {
		msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error reading message")
			break
		}
		if msg != nil {
			fmt.Println("Sending message to channel", msg)
			ch <- msg
		}
	}
}

func (connection *WebsocketConnection) WriteMessage(msg Message) {}

func (connection *WebsocketConnection) Close() {
	log.Println("Closing websocket connection")
	connection.ws.Close()
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	connection := &WebsocketConnection{ws}
	return connection
}

/**
 * Web Socket Gate
 */

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

	go server.Serve(listener)
}

/**
 * Обработчик входящих подключений по websocket
 * @param  {[type]} gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request [description]
 * @return {[type]} [description]
 */
func (gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request) {
	log.Println("WSGate: new websocket connection")
	webSocket, err := upgrader.Upgrade(rw, request, nil)
	log.Println("Upgraded")

	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("WSGate: Not a websocket handshake")
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println("WSGate: Unknown error", err)
		return
	}

	conn := NewWebsocketConnection(webSocket)
	log.Println("WSGate: connection sent to channel")
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
