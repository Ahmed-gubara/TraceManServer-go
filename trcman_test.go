package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"trcman/parser"
)

func TestFrame(t *testing.T) {
	tframe := make([]byte, 0)
	parser.SetPayload(tframe, parser.ProtocolPrefix{ProtocolHead: [2]byte{2, 100}})

	data := "40 40 87 00 04 32 31 33 47 44 50 32 30 31 38 30 32 32 33 38 38 00 00 00 00 10 01 8d a8 3f 5e f8 ab 3f 5e 01 68 0d 00 91 02 00 00 dd 20 00 00 0c 00 00 00 00 00 07 2f 01 57 00 00 00 00 01 09 02 14 06 33 12 42 fe 58 03 34 80 fc 06 00 00 d8 0d bf 49 44 44 5f 32 31 33 57 30 31 5f 53 20 56 32 2e 32 2e 30 00 49 44 44 5f 32 31 33 57 30 31 5f 48 20 56 32 2e 32 2e 30 00 04 00 02 18 01 1f 02 1f 04 1f 5c 29 0d 0a"

	recieved := atb(data)
	// fmt.Printf("before procc %d |%s|\n", len(frame), string(frame))
	recieved, i := parser.GetPayload(recieved, parser.ProtocolPrefix{})
	prefix := i.(parser.ProtocolPrefix)

	recieved, i = parser.GetPayload(recieved, parser.Login0x1001{})
	login0x1001 := i.(parser.Login0x1001)

	frame := parser.Encapsulate(prefix.ProtocolVersion, prefix.DeviceID, 0x9001, login0x1001)

	t.Logf("%+ x\n", recieved)
	t.Logf("%+ x\n", frame)

	t.Error("www")

}
func atb(data string) []byte {
	frame := strings.Split(data, " ")
	byteFrame := make([]byte, 0, len(frame))
	for _, str := range frame {
		i, _ := strconv.ParseInt(str, 16, 0)
		byteFrame = append(byteFrame, byte(i))
	}
	return byteFrame
}
func getOutboundIP() string {
	res, err := http.Get("http://api.ipify.org/")
	handleError("erro", err)
	defer res.Body.Close()

	content, _ := ioutil.ReadAll(res.Body)
	return string(content)
}
func getIP() [4]byte {
	ip := getOutboundIP()

	ss := strings.Split(ip, ".")
	var iparray [4]byte
	for i := 0; i < 4; i++ {
		num, _ := strconv.Atoi(ss[i])
		iparray[i] = byte(num)
	}
	return iparray
}
func handleError(txt string, err error) {
	if err != nil {
		log.Print(txt)
	}
}
