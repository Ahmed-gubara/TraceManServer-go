package telegram

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"trcman/proto"

	bot_api "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/grpc"
)

const token string = "970356809:AAFg6NOlJLtlIJF5OALFP9DWqNXAZjJsiVU"
const chatsFile string = "chat.txt"

var chatList []int64 = []int64{}

func StartTelegramService() {
	// done := make(chan os.Signal, 1)
	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	bot, error := bot_api.NewBotAPI(token)
	if error != nil {
		log.Panicf("tgbotapi.NewBotAPI() failed with %s", error)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	loadChats()
	broadCastMessage(bot, fmt.Sprintf("Server Started at ip <b>%s</b>", getOutboundIP()))
	// go startTCPServer(bot)
	u := bot_api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("bot.GetUpdatesChan(u) failed with %s", err)
	}
	// go func() {
	// 	for {
	// 		select {
	// 		case sig := <-done:
	// 			broadCastMessage(bot, fmt.Sprintf("Server signaled <code>%v</code> at ip %s, Exiting now", sig, getOutboundIP()))
	// 			os.Exit(0)
	// 		}
	// 	}
	// }()
	go startGrpcClient(bot)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		go handleUpdate(&update, bot)
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
	{
		m := bot_api.NewMessage(update.Message.Chat.ID, "<i>unauthorized use of bot, still under development, sorry for the inconvenience 😊</i>")
		m.ParseMode = "HTML"
		bot.Send(m)
	}
	return

authorized:
	// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := bot_api.NewMessage(update.Message.Chat.ID, fmt.Sprintf("no action for <code>(%s)</code>", msgCnt))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "HTML"
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
		message := bot_api.NewMessage(chatid, msg)
		// message.ParseMode = "MarkdownV2"
		message.ParseMode = "HTML"
		bot.Send(message)
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
func panicif(err error) {
	if err != nil {
		panic(err)
	}
}
func startGrpcClient(bot *bot_api.BotAPI) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":4041", grpc.WithInsecure())
	panicif(err)
	defer conn.Close()
	c := proto.NewTrcmanServiceClient(conn)
	broadCastMessage(bot, fmt.Sprintf("subscribed to"))

	for {
		broadCastMessage(bot, fmt.Sprintf("connecting to grpc server"))
		response, err := c.Subscribe(context.Background(), &proto.StringMessage{Content: "hi"})
		if err != nil {
			broadCastMessage(bot, fmt.Sprintf("connection to grpc server failed : %+v", err))
			break
		}
		broadCastMessage(bot, fmt.Sprintf("connected to grpc server"))
		for {
			resp, err := response.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				continue
			}
			broadCastMessage(bot, fmt.Sprintf("notification : %+v", resp.GetContent()))
		}
	}
}
