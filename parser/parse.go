package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Test dd
const Test = false

func atb(data string) []byte {
	frame := strings.Split(data, " ")
	byteFrame := make([]byte, 0, len(frame))
	for _, str := range frame {
		i, _ := strconv.Atoi(str)

		byteFrame = append(byteFrame, byte(i))
	}
	return byteFrame
}

//Maino test
func Maino() {

	// tframe := make([]byte, 0)
	// SetProtocolPrefix(tframe, ProtocolPrefix{protocolHead: []byte{2, 100}})

	// data := "64 64 135 0 4 50 49 51 71 68 80 50 48 49 56 48 50 50 51 56 56 0 0 0 0 16 1 213 213 57 94 35 232 57 94 153 7 11 0 0 0 0 0 79 27 0 0 4 0 0 0 0 0 7 46 1 87 0 0 0 0 1 4 2 20 21 54 42 116 174 88 3 68 133 252 6 0 0 244 8 143 73 68 68 95 50 49 51 87 48 49 95 83 32 86 50 46 50 46 48 0 73 68 68 95 50 49 51 87 48 49 95 72 32 86 50 46 50 46 48 0 4 0 2 24 1 31 2 31 4 31 1 210 13"
	// data += " 10"
	// // data := "48 48 49 49 48 48 48 67 56 49 52 50 57 57 50 49 55 49 54 57 51 52 48 48 48 56 48 48 51 50 48 48 55 55 48 48 54 70 48 48 54 67 48 48 54 54 48 48 55 57 48 48 52 52 48 48 52 53 48 48 52 68 48 48 52 70 48 48 53 70 48 48 51 49 48 48 50 99 48 48 50 48 48 48 54 57 48 48 54 55 48 48 54 101 48 48 54 57 48 48 55 52 48 48 54 57 48 48 54 102 48 48 54 101 48 48 50 48 48 48 54 102 48 48 54 101 48 48 50 49 0 0 0 0 0 64 64 135 0 4 50 49 51 71 68 80 50 48 49 56 48 50 50 51 56 56 0 0 0 0 16 1 88 182 52 94 226 182 52 94 30 101"
	// frame := atb(data)
	// // fmt.Printf("before procc %d |%s|\n", len(frame), string(frame))
	// frame, _ = GetProtocolPrefix(frame)
	// //fmt.Printf("ProtocolPrefix %+v\n", ProtocolPrefix)
	// if false {
	// 	_, statData := GetStatData(frame)
	// 	fmt.Printf("\n\n\nstatData %+v\n", statData)
	// }

	// frame, statData := GetPayload(frame, LoginMessage0x1001{})
	// fmt.Printf("\n\n\nstatData %+v\n", statData)

	d := LoginResponseMessage0x9001{
		IPAddress:  [...]uint8{1, 2, 3, 4},
		Port:       9000,
		ServerTime: time.Now(),
	}
	f := make([]byte, 0)
	f = SetPayload(f, d)
	fmt.Printf("%+v", f)

}
func parseProtocol(frame []byte) {

}
