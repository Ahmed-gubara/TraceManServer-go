package services

import (
	"fmt"
	"net"
	"trcman/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCService(server *serviceServer) {
	listener, err := net.Listen("tcp", ":4041")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	proto.RegisterTrcmanServiceServer(srv, server)
	reflection.Register(srv)
	fmt.Println("starting server")
	go srv.Serve(listener)

	fmt.Printf("server started and listening at %s", listener.Addr().String())

}

// func monitorServer(srv *grpc.Server, lis *net.Listener) {
// 	for {
// 		for name, info := range srv.GetServiceInfo() {
// 			if name == "proto.TrcmanService" {
// 				fmt.Print(name)
// 				fmt.Print(":")
// 				fmt.Println(info.Methods)
// 			}
// 		}
// 		time.Sleep(time.Millisecond * 100)
// 	}
// }
