package network

import (
	"github.com/porfirion/server2/network/http"
	"github.com/porfirion/server2/network/pool"
	"github.com/porfirion/server2/network/tcp"
	"github.com/porfirion/server2/network/ws"
	"github.com/porfirion/server2/service"
	"log"
	"net"
)

type NetworkService struct {
	*service.BasicService
	pool *pool.ConnectionsPool
}

func (s *NetworkService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (s *NetworkService) Start() {
	go s.startReceivingFromBroker()
	go s.startReadingFromClients()
}

func (s *NetworkService) startReceivingFromBroker() {
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		s.pool.OutgoingMessages <- pool.MessageForClient{
			Targets: msg.DestinationServiceClients,
			Data:    msg.MessageData,
		}
	}
}

func (s *NetworkService) startReadingFromClients() {
	// TODO здесь ещё нужна проверка на то, зарегистрировали ли нас и есть ли нам куда писать
	for msg := range s.pool.IncomingMessages {
		//log.Printf("NetworkService: Received msg %T", msg)
		// TODO сейчас в пробкер отправляются сырые байты и никакого парсинга не происходит
		// также не указывается целевой сервис, в который мы отправляем эти данные
		s.SendMessageToBroker(msg.Data, msg.ClientId, 0, 0, nil)
	}

}

func NewService(wsport, tcpport, httpport int, nostatic bool) *NetworkService {
	pool := pool.NewConnectionsPool()
	go pool.Start()
	log.Println("Pool started")

	wsGate := &ws.WebSocketGate{
		Addr: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: wsport},
		Pool: pool,
	}
	go wsGate.Start()
	log.Println("WsGate started")

	tcpGate := &tcp.TcpGate{
		Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: tcpport},
		Pool: pool,
	}
	go tcpGate.Start()
	log.Println("TcpGate started")

	if !nostatic {
		go http.ServeStatic(httpport)
		log.Println("HttpStaticServer started")
	}

	return &NetworkService{
		BasicService: service.NewBasicService(service.TypeNetwork),
		pool:         pool,
	}
}
