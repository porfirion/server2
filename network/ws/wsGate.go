package ws

import (
	"github.com/gorilla/websocket"
	"github.com/porfirion/server2/network/pool"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type WebSocketGate struct {
	Addr *net.TCPAddr
	Pool *pool.ConnectionsPool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (gate *WebSocketGate) Start() error {
	http.HandleFunc("/", gate.indexHandler)
	http.HandleFunc("/assets/", gate.assetsHandler)
	http.HandleFunc("/ws", gate.wsHandler)

	server := &http.Server{}
	listener, err := net.ListenTCP("tcp4", gate.Addr)

	if err != nil {
		log.Printf("Error creating listener %v", err)
		return err
	} else {
		log.Println("Listening http:", gate.Addr)
		server.Serve(listener)
		return nil
	}
}

// Обработчик входящих подключений по websocket
// @param  {[type]} gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request [description]
// @return {[type]} [description]
func (gate *WebSocketGate) wsHandler(rw http.ResponseWriter, request *http.Request) {
	webSocket, err := upgrader.Upgrade(rw, request, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			log.Println("WSGate: Not a websocket handshake")
			http.Error(rw, "Not a websocket handshake", 400)
			return
		} else {
			log.Println("WSGate: Unknown error", err)
			return
		}
	}

	conn := NewWebSocketConnection(
		<-gate.Pool.ConnectionsEnumerator,
		gate.Pool.IncomingMessages,
		gate.Pool.ClosingChannel,
		webSocket,
	)

	log.Println("WSGate: new websocket connection", conn)

	// отправляем соединение в пул, пусть дальше он разбирается
	gate.Pool.IncomingConnections <- conn
}

// Отдаёт главную (и единственную) страницу
// @param  {[type]} gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request [description]
// @return {[type]} [description]
func (gate *WebSocketGate) indexHandler(rw http.ResponseWriter, request *http.Request) {
	var path string
	if request.RequestURI == "/" {
		path = "templates/index.html"
	} else {
		path = "templates/" + strings.Replace(request.RequestURI[1:], "..", "", -1)
	}
	if _, err := os.Stat(path); err == nil {
		indexTempl := template.Must(template.ParseFiles(path))
		data := struct{}{}
		indexTempl.Execute(rw, data)
	} else {
		rw.WriteHeader(404)
		rw.Write([]byte(path))
		rw.Write([]byte("404 Not Found"))
	}
}

// Отвечает за отдачу статики
// @param  {[type]} gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request [description]
// @return {[type]} [description]
func (gate *WebSocketGate) assetsHandler(rw http.ResponseWriter, request *http.Request) {
	http.ServeFile(rw, request, request.URL.Path[1:])
}
