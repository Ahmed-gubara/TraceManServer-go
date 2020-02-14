package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
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
		recievedPayload, i := parser.GetPayload(recieved, parser.ProtocolPrefix{})
		prefix := i.(parser.ProtocolPrefix)
		fmt.Printf("ProtocolPrefix %+v", prefix)
		msgType := parser.GetMessageType(prefix.ProtocolID)
		Broadcast(fmt.Sprintf("Received 0x%x %s (%d Byte) from %s hex : \n<code>% x</code>", prefix.ProtocolID, msgType, len(recieved), prefix.DeviceID, recieved))
		switch prefix.ProtocolID {
		case 0x1001:
			_, payload := parser.GetPayload(recievedPayload, parser.Login0x1001{})
			login0x1001 := payload.(parser.Login0x1001)
			Broadcast(fmt.Sprintf("Received and parsed 0x%x %s (%d Byte) from %s json : \n<code>%+v</code>", prefix.ProtocolID, msgType, len(recievedPayload), prefix.DeviceID, fr(login0x1001)))
			lResponse := parser.LoginResponse0x9001{IPAddress: getIP(), Port: 9000, ServerTime: time.Now().UTC()}
			Broadcast(fmt.Sprintf("Respoinding : \n<code>%+v</code>", lResponse))
			frame := parser.Encapsulate(prefix.ProtocolVersion, prefix.DeviceID, 0x9001, lResponse)
			Broadcast(fmt.Sprintf("sending 0x%x %s (%d Byte) from %s hex : \n<code>% x</code>", 0x9001, parser.GetMessageType(0x9001), len(frame), prefix.DeviceID, frame))
			obdconn.send <- frame

		case 0x1002:
			_, payload := parser.GetPayload(recievedPayload, parser.Logout0x1002{})
			logout0x1002 := payload.(parser.Logout0x1002)
			Broadcast(fmt.Sprintf("Received and parsed 0x%x %s (%d Byte) from %s json : \n<code>%+v</code>", prefix.ProtocolID, msgType, len(recievedPayload), prefix.DeviceID, fr(logout0x1002)))

		case 0x1003: // heartbeat message
			frame := parser.Encapsulate(prefix.ProtocolVersion, prefix.DeviceID, 0x9003, parser.HearbeatResponse0x9003{})
			obdconn.send <- frame

		case 0x4009:
			_, payload := parser.GetPayload(recievedPayload, parser.GPSinSleep0x4009{})
			gpsinSleep0x4009 := payload.(parser.GPSinSleep0x4009)
			Broadcast(fmt.Sprintf("Received and parsed 0x%x %s (%d Byte) from %s json : \n<code>%+v</code>", prefix.ProtocolID, msgType, len(recievedPayload), prefix.DeviceID, fr(gpsinSleep0x4009)))

		case 0x4001:
			_, payload := parser.GetPayload(recievedPayload, parser.GPSData0x4001{})
			gpsData0x4001 := payload.(parser.GPSData0x4001)
			Broadcast(fmt.Sprintf("Received and parsed 0x%x %s (%d Byte) from %s json : \n<code>%+v</code>", prefix.ProtocolID, msgType, len(recievedPayload), prefix.DeviceID, fr(gpsData0x4001)))

		default:
			Broadcast(fmt.Sprintf("not handled"))
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
func fr(data interface{}) string {
	str, _ := json.MarshalIndent(data, "", "    ")
	return string(str)
}
