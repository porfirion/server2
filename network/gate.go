package network

type Gate interface {
	Start(ConnectionsChannel, MessagesChannel)
}
