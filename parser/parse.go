package parser

import (
	"fmt"
	"strconv"
	"strings"
)

//Test dd
const Test = true

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

	tframe := make([]byte, 0)
	SetProtocolPrefix(tframe, ProtocolPrefix{protocolHead: []byte{2, 100}})

	data := "64 64 135 0 4 50 49 51 71 68 80 50 48 49 56 48 50 50 51 56 56 0 0 0 0 16 1 88 182 52 94 107 183 52 94 30 101"
	// data := "48 48 49 49 48 48 48 67 56 49 52 50 57 57 50 49 55 49 54 57 51 52 48 48 48 56 48 48 51 50 48 48 55 55 48 48 54 70 48 48 54 67 48 48 54 54 48 48 55 57 48 48 52 52 48 48 52 53 48 48 52 68 48 48 52 70 48 48 53 70 48 48 51 49 48 48 50 99 48 48 50 48 48 48 54 57 48 48 54 55 48 48 54 101 48 48 54 57 48 48 55 52 48 48 54 57 48 48 54 102 48 48 54 101 48 48 50 48 48 48 54 102 48 48 54 101 48 48 50 49 0 0 0 0 0 64 64 135 0 4 50 49 51 71 68 80 50 48 49 56 48 50 50 51 56 56 0 0 0 0 16 1 88 182 52 94 226 182 52 94 30 101"
	frame := atb(data)
	fmt.Printf("before procc %d |%s|\n", len(frame), string(frame))
	frame, ProtocolPrefix := GetProtocolPrefix(frame)
	fmt.Printf("ProtocolPrefix %+v\n", ProtocolPrefix)
	frame, statData := GetStatData(frame)
	fmt.Printf("statData %+v\n", statData)
	// frame, ProProtocolSufix := GetProtocolSufix(frame)
	// fmt.Printf("ProtocolPrefix %+v \nstatData %+v\nProProtocolSufix %+v\n", ProtocolPrefix, statData, ProProtocolSufix)
	fmt.Printf("after procc %d |%s|\n", len(frame), string(frame))

}
func parseProtocol(frame []byte) {

}
