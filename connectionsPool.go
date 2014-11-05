package main

type ConnectionsPool struct {
	peer2sserver        chan Message    // сообщения, которые отправляются на обработку в сервер
	server2peer         chan Message    // сообщения, которые приходят на отправку из сервера
	incomingConnections chan Connection // входящие соединения
}

func (pool *ConnectionsPool) processConnection(conn *Connection) {

}

func (pool *ConnectionsPool) start() {

}
