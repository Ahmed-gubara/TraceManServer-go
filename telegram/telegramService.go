package telegram

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	bot_api "github.com/go-telegram-bot-api/telegram-bot-api"
)

const token string = "970356809:AAFg6NOlJLtlIJF5OALFP9DWqNXAZjJsiVU"
const chatsFile string = "chat.txt"

var chatList []int64 = []int64{}

func StartService() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	bot, error := bot_api.NewBotAPI(token)
	if error != nil {
		log.Panicf("tgbotapi.NewBotAPI() failed with %s", error)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	loadChats()
	broadCastMessage(bot, fmt.Sprintf("Server Started at ip %s", getOutboundIP()))
	go startTCPServer(bot)
	u := bot_api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("bot.GetUpdatesChan(u) failed with %s", err)
	}
	for {
		select {
		case msg := <-done:
			broadCastMessage(bot, fmt.Sprintf("Server Killed at ip %s %s", msg.String(), getOutboundIP()))
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}
			go handleUpdate(&update, bot)
		default:

		}
	}

}
func handleUpdate(update *bot_api.Update, bot *bot_api.BotAPI) {
	msgCnt := update.Message.Text
	for i, chatid := range chatList {
		if chatid == update.Message.Chat.ID {
			if msgCnt == "unbind" {
				chatList = append(chatList[:i], chatList[i+1:]...)
				go saveChats()
				msg := bot_api.NewMessage(update.Message.Chat.ID, "chat unsaved!")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return
			}
			goto authorized
		}
	}
	if msgCnt == "bind" {
		chatList = append(chatList, update.Message.Chat.ID)
		go saveChats()
		msg := bot_api.NewMessage(update.Message.Chat.ID, "chat saved!")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		return
		//goto authorized
	}

	bot.Send(bot_api.NewMessage(update.Message.Chat.ID, "unauthorized use of bot, still under development, sorry for the inconvenience 😊"))
	return

authorized:
	// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := bot_api.NewMessage(update.Message.Chat.ID, fmt.Sprintf("no action for (%s)", msgCnt))
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}
func saveChats() {
	file, err := os.Create(chatsFile)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(arrayToString(chatList, ","))
}
func loadChats() {
	data, err := ioutil.ReadFile(chatsFile)
	if err != nil || len(data) == 0 {
		log.Printf("loadChats :: err %s , len %d", err, len(data))
		return
	}
	str := string(data)

	strings := strings.Split(str, ",")
	chatList = make([]int64, len(strings))

	for i, s := range strings {
		chatList[i], _ = strconv.ParseInt(s, 10, 64)
	}

}
func arrayToString(a []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}
func broadCastMessage(bot *bot_api.BotAPI, msg string) {
	for _, chatid := range chatList {
		log.Printf("sending to chatid %d", chatid)
		bot.Send(bot_api.NewMessage(chatid, msg))
	}
}
func startTCPServer(bot *bot_api.BotAPI) {
	// listen on port 8000
	var ln net.Listener
	var err error
	switch runtime.GOOS {
	case "windows":
		ln, err = net.Listen("tcp", ":9000")

	default:
		ln, err = net.Listen("tcp", ":9000")
	}
	if err != nil {
		panic(err)
	}
	for {
		// accept connection
		conn, _ := ln.Accept()
		go handleTCPConnection(conn, bot)
		// run loop forever (or until ctrl-c)
	}

}
func handleTCPConnection(conn net.Conn, bot *bot_api.BotAPI) {
	broadCastMessage(bot, fmt.Sprintf("connection started with ip %s", conn.RemoteAddr().String()))
	defer func() {
		conn.Close()
		broadCastMessage(bot, fmt.Sprintf("connection closed from ip %s", conn.RemoteAddr().String()))
	}()
	for {
		bufio.NewScanner(conn)
		// get message, output
		message, err := bufio.NewReader(conn).ReadBytes('\r')
		if err != nil {
			fmt.Println(err)
			return
		}
		message = append(message, 10) // add \n to match \r\n pattern
		//conn.Write([]byte(gen0x9001()))
		// broadCastMessage(bot, fmt.Sprintf("Message Received : %s", message))
		broadCastMessage(bot, fmt.Sprintf("Message Received (%v Byte) : %v", len(message), message))
		temp := strings.TrimSpace(string(message))
		if temp == "STOP" {
			break
		}
	}
}
func getOutboundIP() string {
	res, err := http.Get("http://api.ipify.org/")
	handleError("erro", err)
	defer res.Body.Close()

	content, _ := ioutil.ReadAll(res.Body)
	return string(content)
}
func handleError(txt string, err error) {
	if err != nil {
		log.Print(txt)
	}
}
func gen0x9001() string {
	str := []string{}
	frame := strings.Join(str, "") + "8000" + fmt.Sprint(time.Now().Unix())
	return frame
}
