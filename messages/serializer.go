package messages

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/porfirion/server2/service"
)

// когда мы читаем сообщение от клиента - мы никогда не знаем заранее что это будет
// поэтому используется RawMessage, чтобы не десериализовывать само сообщение в непонятный тип
type jsonIncomingMessageWrapper struct {
	MessageType uint64          `json:"type"`
	Data        json.RawMessage `json:"data"`
}

func (w *jsonIncomingMessageWrapper) GetType() uint64 {
	return w.MessageType
}

// а вот когда мы отправляем сообщение - мы уже заранее всё знаем,
// поэтому сериализуем всё за один проход
type jsonOutgoingMessageWrapper struct {
	MessageType uint64      `json:"type"`
	Data        interface{} `json:"data"`
}

func DeserializeFromJson(bytes []byte) (service.TypedMessage, error) {
	msg := &jsonIncomingMessageWrapper{}
	if err := json.Unmarshal(bytes, msg); err == nil {
		var typedMessage service.TypedMessage

		switch msg.GetType() {
		case 1:
			typedMessage = &AuthMessage{}
		case 10:
			typedMessage = &LoginMessage{}
		case 1001:
			typedMessage = &TextMessage{}
		case 20002:
			typedMessage = &SyncTimeMessage{}
		case 1000001:
			typedMessage = &SimulateMessage{}
		case 1000002:
			typedMessage = &ChangeSimulationMode{}
		default:
			return nil, fmt.Errorf("Unknown message type %d", msg.GetType())
		}

		err := json.Unmarshal(msg.Data, typedMessage)
		return typedMessage, err
	} else {
		return msg, err
	}
}

func SerializeToJson(msg service.TypedMessage) ([]byte, error) {
	return json.Marshal(jsonOutgoingMessageWrapper{MessageType: msg.GetType(), Data: msg})
}

type TypedBytesMessage []byte

func (t TypedBytesMessage) GetType() uint64 {
	return binary.BigEndian.Uint64(t[:8])
}

func DeserializeFromBinary(bytes []byte) (msg service.TypedMessage, err error) {
	return TypedBytesMessage(bytes), nil
}

func SerializeToBinary(msg service.TypedMessage) ([]byte, error) {
	if bytes, err := json.Marshal(msg); err == nil {
		var res = make([]byte, 0, 8+len(bytes))
		binary.BigEndian.PutUint64(res, msg.GetType())
		res = append(res, bytes...)
		return res, nil
	} else {
		return nil, err
	}
}
