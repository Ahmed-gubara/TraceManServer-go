package services

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func StartTCPServer(obdConnection chan<- *OBDConnection) <-chan string {
	outChan := make(chan string, 10)
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
	go func() {
		for {
			// accept connection
			conn, _ := ln.Accept()
			obdcon := OBDConnection{recieved: make(chan []byte), send: make(chan []byte), isConnected: true}
			obdConnection <- &obdcon
			go handleTCPConnection(conn, &obdcon)
			// run loop forever (or until ctrl-c)
		}
	}()
	return outChan
}
func handleTCPConnection(conn net.Conn, obdcon *OBDConnection) {
	// broadCastMessage(bot, fmt.Sprintf("connection started with ip %s", conn.RemoteAddr().String()))

	defer func() {
		conn.Close()
		obdcon.isConnected = false
		close(obdcon.recieved)
		// close(obdcon.send)
		// broadCastMessage(bot, fmt.Sprintf("connection closed from ip %s", conn.RemoteAddr().String()))
	}()
	scanner := bufio.NewScanner(conn)
	scanner.Split(splitter)
	for scanner.Scan() {
		message := scanner.Bytes()
		fmt.Printf("maching message recieved %v", message)
		obdcon.recieved <- message
		select {
		case send := <-obdcon.send:
			conn.Write(send)
		}
		//fmt.Sprintf("A message Received (%d Byte) hex : \n<code>% x</code>", len(message), message)
		// broadCastMessage(bot, fmt.Sprintf("A message Received (%d Byte) hex : \n<code>% x</code>", len(message), message))
		// message := nil
		// // get message, output
		// // message, err := bufio.NewReader(conn).ReadBytes('\r') //	 add \n to match \r\n pattern
		// //conn.Write([]byte(gen0x9001()))
		// // broadCastMessage(bot, fmt.Sprintf("Message Received : %s", message))
		// broadCastMessage(bot, fmt.Sprintf("Message Received (%v Byte) : %v", len(message), message))
		// temp := strings.TrimSpace(string(message))
		// if temp == "STOP" {
		// 	break
		// }
	}
}
func splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	fmt.Printf("a %+v\n", data)

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{0x40, 0x40}); i >= 0 && i < (len(data)-3) {

		data = data[i:]

		size := binary.LittleEndian.Uint16(data[2:4])
		if int(size) > (len(data)) {
			fmt.Printf("b size %d actual %d %+v\n", size, len(data), data)

			// return int(size) - (len(data)), nil, nil
			return
		}
		fmt.Printf("d size %d %+v\n", size, data)
		return int(size), data[:size], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
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
