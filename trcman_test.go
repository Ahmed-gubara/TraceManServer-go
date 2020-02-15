package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"trcman/parser"
)

func TestFrame(t *testing.T) {
	// return
	// tframe := make([]byte, 0)
	// parser.SetPayload(tframe, parser.ProtocolPrefix{ProtocolHead: [2]byte{2, 100}})
	data := "40 40 e8 01 04 32 31 33 47 44 50 32 30 31 38 30 32 32 33 38 38 00 00 00 00 40 01 01 76 01 3e 5e f7 03 3e 5e cb fd 0c 00 00 00 00 00 f1 1f 00 00 00 00 00 00 04 00 07 39 01 57 00 00 03 00 14 08 02 14 00 20 32 c0 b6 58 03 40 8c fc 06 00 00 08 0d 6f 08 02 14 00 21 14 2c a5 58 03 00 8d fc 06 99 02 ca 06 af 08 02 14 00 21 32 a8 92 58 03 e0 8e fc 06 60 04 1e 07 cf 08 02 14 00 22 14 70 67 58 03 d6 92 fc 06 e3 03 2c 07 bf 08 02 14 00 22 32 5e 40 58 03 a2 99 fc 06 d4 03 3a 07 bf 08 02 14 00 23 14 8e 2e 58 03 a4 ac fc 06 92 04 1a 03 cf 08 02 14 00 23 32 46 37 58 03 06 e4 fc 06 b0 06 22 03 bf 08 02 14 00 24 14 a0 40 58 03 c6 1d fd 06 b3 04 20 03 cf 08 02 14 00 24 32 ca 25 58 03 de 29 fd 06 85 04 93 06 af 08 02 14 00 25 14 28 f5 57 03 18 32 fd 06 4d 04 d0 06 af 08 02 14 00 25 32 06 e4 57 03 f0 20 fd 06 ce 04 56 0a cf 08 02 14 00 26 14 86 d6 57 03 a6 01 fd 06 e4 03 a6 06 cf 08 02 14 00 26 32 0a ad 57 03 08 09 fd 06 a1 02 d1 06 cf 08 02 14 00 27 14 5a a8 57 03 50 fd fc 06 b4 01 ca 06 cf 08 02 14 00 27 32 3a 92 57 03 70 01 fd 06 00 00 5b 07 bf 08 02 14 00 28 14 3a 92 57 03 70 01 fd 06 00 00 5b 07 bf 08 02 14 00 28 32 3a 92 57 03 70 01 fd 06 00 00 5b 07 cf 08 02 14 00 29 14 3a 92 57 03 70 01 fd 06 00 00 5b 07 cf 08 02 14 00 29 32 3a 92 57 03 70 01 fd 06 00 00 5b 07 cf 08 02 14 00 2a 14 3a 92 57 03 70 01 fd 06 00 00 5b 07 cf 14 25 04 f2 05 2e 08 15 07 24 07 9d 08 40 07 df 04 c3 06 4e 06 8c 09 7c 07 29 03 8a 07 31 03 ee 02 f1 02 e6 02 f6 02 03 03 5f f6 0d 0a"

	rframe := atb(data)
	// fmt.Printf("before procc %d |%s|\n", len(frame), string(frame))
	recieved, i := parser.GetPayload(rframe, parser.ProtocolPrefix{})
	prefix := i.(parser.ProtocolPrefix)

	recieved, i = parser.GetPayload(recieved, parser.GPSData0x4001{})
	login0x1001 := i.(parser.GPSData0x4001)
	x, _ := json.MarshalIndent(login0x1001, "", " ")
	log.Printf("%+v", string(x))
	frame := parser.Encapsulate(prefix.ProtocolVersion, prefix.DeviceID, 0x4001, login0x1001)

	if len(frame) != len(frame) {
		t.Errorf("length not equal r %d s %d", len(rframe), len(frame))
		return
	}
	var c bool
	for i := 0; i < len(rframe); i++ {
		if frame[i] != rframe[i] {
			c = true
			t.Logf("mismaching data index %d s %v r %v", i, frame[i], rframe[i])
		}
	}

	if c {
		t.Error("failed")
	}
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
