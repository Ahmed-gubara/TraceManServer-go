package main

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"os"
	services "trcman/services"
	tg "trcman/telegram"
)

func main() {
	OBDconnections := make(chan *services.OBDConnection, 10)

	println("A")
	server := services.StartTrcManServer(OBDconnections)

	println("B")
	services.StartTCPServer(OBDconnections)

	println("C")
	services.StartGRPCService(&server)

	//---------------------------------
	println("D")
	tg.StartTelegramService()
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyCtrlC:
				fmt.Println("stopping")
				os.Exit(0)
				return
			}
		default:
			term.Sync()
		}
	}
}
