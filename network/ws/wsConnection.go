package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/messages"
)

type WebsocketConnection struct {
	*network.BasicConnection
	ws *websocket.Conn
}

/*WARNING! это можно вызывать только из того треда, который отправляет собщения в канал*/
func (connection *WebsocketConnection) Close(message string) {
	log.Println("WSCon: Closing websocket connection")
	close(connection.OutgoingChannel)
	connection.ws.WriteControl(websocket.CloseMessage, []byte(message), time.Time{})
	connection.ws.Close()
}

func (connection *WebsocketConnection) StartReading() {
	go func() {
		defer func() {
			// нельзя закрывать самому себя - это может сделать только тот, кто пишет в канал
			//connection.Close(0, "unimplemented")

			connection.ws.Close()
			connection.NotifyPoolWeAreClosing()
		}()

		for {
			if msgType, buffer, err := connection.ws.ReadMessage(); err == nil {
				switch msgType {
				case websocket.CloseMessage:
					log.Println("Close message received")
				case websocket.PingMessage:
					log.Println("Ping message received")
				case websocket.TextMessage:
					if msg, err := messages.DeserializeFromJson(buffer); err != nil {
						log.Println("Error parsing wrapper", err)
					} else {
						connection.Notify(msg)

						/*
							case *network.SyncTimeMessage:
							// log.Println("Sync time message", time.Now().UnixNano()/int64(time.Millisecond), int64(time.Now().UnixNano()/int64(time.Millisecond)))
							connection.OutgoingChannel <- network.SyncTimeMessage{Time: int64(time.Now().UnixNano() / int64(time.Millisecond))}
						*/

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
				break
			}
		}

		log.Println("Reading finished for", connection.Id)
	}()
}
func (connection *WebsocketConnection) StartWriting() {
	go func() {
		//log.Println("Writing started ", connection.id)

		for message := range connection.OutgoingChannel {
			//log.Println(fmt.Sprintf("WsCon. Sending message %T for %d", message, connection.id))
			if bytes, err := messages.SerializeToJson(message.Data); err == nil {
				connection.ws.WriteMessage(websocket.TextMessage, bytes)
			}
		}

		log.Println("Writing finished for", connection.Id)
	}()
}

func NewWebsocketConnection(
	id uint64,
	incoming chan network.MessageFromClient,
	closingChannel chan uint64,
	ws *websocket.Conn,
) network.Connection {
	connection := &WebsocketConnection{
		ws: ws,
		BasicConnection: network.NewBasicConnection(
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
