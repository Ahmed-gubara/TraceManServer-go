package services

import (
	"context"
	"fmt"
	"sync"
	"trcman/parser"
	"trcman/proto"

	guuid "github.com/google/uuid"
)

type subscription struct {
	channel chan string
	filter  func(string) bool
}

type OBDConnection struct {
	ID          string
	recieved    chan []byte
	send        chan []byte
	isConnected bool
}
type serviceServer struct {
}

var allRegisteredClients map[string]subscription = make(map[string]subscription)

var clientsLock = sync.RWMutex{}

func StartTrcManServer(connection <-chan *OBDConnection) serviceServer {

	go func() {
		for obdconn := range connection {
			go handleOBDConnection(obdconn)

		}
	}()

	return serviceServer{}
}
func handleOBDConnection(obdconn *OBDConnection) {
	for recieved := range obdconn.recieved {
		recieved, prefix := parser.GetProtocolPrefix(recieved)
		msgType := parser.GetMessageType(prefix.ProtocolID)
		Broadcast(fmt.Sprintf("Received 0x%x %s (%d Byte) from %s hex : \n<code>% x</code>", prefix.ProtocolID, msgType, len(recieved), prefix.DeviceID, recieved))
		switch prefix.ProtocolID {
		case 0x1001:
			_, payload := parser.GetPayload(recieved, parser.Login0x1001{})
			login0x1001 := payload.(parser.Login0x1001)
			Broadcast(fmt.Sprintf("Received and parsed 0x%x %s (%d Byte) from %s hex : \n<code>%+v</code>", prefix.ProtocolID, msgType, len(recieved), prefix.DeviceID, login0x1001))

		}
	}
}

// func ()  {

// }
func (s *serviceServer) Subscribe(m *proto.StringMessage, w proto.TrcmanService_SubscribeServer) error {
	filter := func(data string) bool {
		return true
	}
	id := guuid.New().String()
	clientsLock.Lock()
	ch := make(chan string, 10)
	allRegisteredClients[id] = subscription{channel: ch, filter: filter}
	clientsLock.Unlock()
	fmt.Printf("client connected, current active clients : %d\n", len(allRegisteredClients))
	w.Send(&proto.StringMessage{Content: fmt.Sprintf("welcome")})
	for msg := range ch {
		if err := w.Send(&proto.StringMessage{Content: msg}); err != nil {
			break
		}
		// send message
		// Deal with errors
		// Deal with client terminations
	}
	fmt.Println("client disconnected")
	clientsLock.Lock()
	delete(allRegisteredClients, id)
	clientsLock.Unlock()
	return nil
}
func (s *serviceServer) IsServiceRunning(context.Context, *proto.StringMessage) (*proto.StringMessage, error) {
	return &proto.StringMessage{Content: "service is running correctly"}, nil
}

func (s *serviceServer) Publish(ctx context.Context, msg *proto.StringMessage) (*proto.StringMessage, error) {
	if len(msg.GetContent()) > 0 {
		Broadcast(msg.GetContent())
	}
	return &proto.StringMessage{}, nil
}

//Broadcast ss
func Broadcast(msg string) {
	clientsLock.RLock()
	fmt.Printf("sending to  %d client", len(allRegisteredClients))
	for _, x := range allRegisteredClients {
		if x.filter(msg) {
			select {
			case x.channel <- msg:
			default:
				fmt.Printf("scipped client by filter")
			}
		}
	}
	clientsLock.RUnlock()
}
