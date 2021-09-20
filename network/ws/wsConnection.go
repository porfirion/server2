package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/network/pool"
	"github.com/porfirion/server2/service"
)

type WebsocketConnection struct {
	*pool.BasicConnection
	ws *websocket.Conn
}

// WARNING! это можно вызывать только из того треда, который отправляет собщения в канал
// А именно - это либо гейт, либо пул
func (connection *WebsocketConnection) Close(message string) {
	log.Println("WSCon: Closing websocket connection")
	close(connection.OutgoingChannel)
	err := connection.ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, message), time.Time{})
	if err != nil {
		log.Println("Error: ", err)
	}
	_ = connection.ws.Close()
}

func (connection *WebsocketConnection) StartReading() {
	go func() {
		defer func() {
			// нельзя закрывать самому себя - это может сделать только тот, кто пишет в канал
			//connection.Close(0, "unimplemented")

			connection.ws.Close()
			connection.NotifyPoolWeAreClosing()
		}()

	ReadingLoop:
		for {
			if msgType, buffer, err := connection.ws.ReadMessage(); err == nil {
				switch msgType {
				case websocket.CloseMessage:
					log.Println("Close message received")
					break ReadingLoop
				case websocket.PingMessage:
					log.Println("Ping message received")
				case websocket.TextMessage:
					if msg, err := messages.DeserializeFromJson(buffer); err == nil {
						switch msg.(type) {
						case *messages.SyncTimeMessage:
							connection.WriteMessage(messages.SyncTimeMessage{Time: time.Now().UnixNano() / int64(time.Millisecond)})
						default:
							//log.Printf("WsConnection: sending message %T to pool\n", msg)
							err := connection.NotifyPoolMessage(msg)
							if err != nil {
								log.Println("Error notify pool message", err)
							}
						}

					} else {
						log.Println("Error parsing wrapper", err)
					}
				case websocket.BinaryMessage:
					log.Println("Binary message!")
				}
			} else {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					// всё штатно закрылось
				} else {
					fmt.Println("Error reading from connection", err.Error())
				}
				break ReadingLoop
			}
		}

		log.Println("Reading finished for", connection.Id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		//log.Println("Writing started ", connection.id)

		for message := range connection.OutgoingChannel {
			if data, ok := message.Data.(service.TypedMessage); ok {
				// log.Println(fmt.Sprintf("WsCon. Sending message %T for %d", message, connection.id))
				if bytes, err := messages.SerializeToJson(data); err == nil {
					err := connection.ws.WriteMessage(websocket.TextMessage, bytes)
					if err != nil {
						log.Println("Error sending to websocket", err)
					}
				} else {
					log.Println("WsConnection: error serializing message", err)
				}
			} else {
				if bts, err := json.Marshal(message.Data); err == nil {
					connection.ws.WriteMessage(websocket.TextMessage, bts)
				} else {
					log.Println("errro serializing message", err)
				}
			}
		}

		log.Println("Writing finished for", connection.Id)
	}()
}

func NewWebSocketConnection(
	id uint64,
	incoming chan pool.MessageFromClient,
	closingChannel chan uint64,
	ws *websocket.Conn,
) pool.Connection {
	connection := &WebsocketConnection{
		ws: ws,
		BasicConnection: pool.NewBasicConnection(
			id,
			incoming,
			closingChannel,
		),
	}

	go connection.StartReading()
	go connection.StartWriting()

	log.Println("Created new ws connection", connection)
	return connection
}
