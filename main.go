package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"trcman/proto"

	guuid "github.com/google/uuid"
	term "github.com/nsf/termbox-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serviceServer struct {
}
type subscription struct {
	channel chan string
	filter  func(string) bool
}

var allRegisteredClients map[string]subscription = make(map[string]subscription)
var clientsLock = sync.RWMutex{}

func main() {

	term.Init()
	listener, err := net.Listen("tcp", ":4041")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	proto.RegisterTrcmanServiceServer(srv, &serviceServer{})
	reflection.Register(srv)
	fmt.Println("starting server")
	go srv.Serve(listener)
	fmt.Println("server started and listening at {}", listener.Addr().String())
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyCtrlC:
				fmt.Println("stopping")
				srv.Stop()
				fmt.Println("stopped")
				return
			}
		default:
			term.Sync()
		}
	}

}
func monitorServer(srv *grpc.Server, lis *net.Listener) {
	for {
		for name, info := range srv.GetServiceInfo() {
			if name == "proto.TrcmanService" {
				fmt.Print(name)
				fmt.Print(":")
				fmt.Println(info.Methods)
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
func (s *serviceServer) IsServiceRunning(ctx context.Context, request *proto.StringMessage) (*proto.StringMessage, error) {
	message := request.GetContent()
	return &proto.StringMessage{Content: fmt.Sprintf("hello %s", message)}, nil
}
func (s *serviceServer) Subscribe(m *proto.StringMessage, w proto.TrcmanService_SubscribeServer) error {
	filter := func(data string) bool {
		return strings.Contains(data, "yes")
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

func (s *serviceServer) Publish(ctx context.Context, msg *proto.StringMessage) (*proto.StringMessage, error) {
	if len(msg.GetContent()) > 0 {
		Broadcast(msg.GetContent())
	}
	return &proto.StringMessage{}, nil
}

//Broadcast ss
func Broadcast(msg string) {
	clientsLock.RLock()
	for _, x := range allRegisteredClients {
		if x.filter(msg) {
			select {
			case x.channel <- msg:
			default:
			}
		}
	}
	clientsLock.RUnlock()
}
