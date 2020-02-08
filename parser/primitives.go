package parser

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type inout struct {
	in  func(frame []byte, value ...interface{}) []byte //interface{}
	out func(frame []byte) interface{}
}

func outU8_1(frame []byte) ([]byte, uint8) {
	f, u := outU8(frame, 1)
	return f, u[0]
}

func outU8(frame []byte, length int) ([]byte, []uint8) {
	// output := make([]byte, length)
	// for index := 0; index < length; index++ {
	// 	output[index] = frame[index]
	// }
	return frame[length:], frame[:length]
}
func inU8(frame []byte, params ...uint8) []byte {
	return append(frame, params...)
}
func outS8_1(frame []byte) ([]byte, int8) {
	f, u := outS8(frame, 1)
	return f, u[0]
}
func outS8(frame []byte, length int) ([]byte, []int8) {
	output := make([]int8, length)
	for index := 0; index < length; index++ {
		output[index] = int8(frame[index])
	}
	return frame[length:], output
}
func inS8(frame []byte, params ...int8) []byte {
	temp := make([]byte, len(params))
	for _, num := range params {
		temp = append(temp, byte(num))
	}
	return append(frame, temp...)
}
func outU16LE1(frame []byte) ([]byte, uint16) {
	f, u := outU16LE(frame, 1)
	return f, u[0]
}
func outU16LE(frame []byte, length int) ([]byte, []uint16) {
	return outU16(frame, length, binary.LittleEndian.Uint16)
}
func outU16BE1(frame []byte) ([]byte, uint16) {
	f, u := outU16BE(frame, 1)
	return f, u[0]
}
func outU16BE(frame []byte, length int) ([]byte, []uint16) {
	return outU16(frame, length, binary.BigEndian.Uint16)
}
func outU16(frame []byte, length int, f func([]byte) uint16) ([]byte, []uint16) {
	size := 2
	output := make([]uint16, length)
	for index := 0; index < length*size; index += size {
		output[index/size] = f(frame[index : index+size])
	}
	return frame[length*size:], output
}
func inU16LE(frame []byte, params ...uint16) []byte {
	return inU16(frame, binary.LittleEndian.PutUint16, params...)
}
func inU16BE(frame []byte, params ...uint16) []byte {
	return inU16(frame, binary.BigEndian.PutUint16, params...)
}
func inU16(frame []byte, f func([]byte, uint16), params ...uint16) []byte {
	size := 2
	temp := make([]byte, 0)
	for _, num := range params {
		u8s := make([]byte, size)
		f(u8s, num)
		temp = append(temp, u8s...)
	}
	return append(frame, temp...)
}

func outS16LE1(frame []byte) ([]byte, int16) {
	f, u := outS16LE(frame, 1)
	return f, u[0]
}
func outS16LE(frame []byte, length int) ([]byte, []int16) {
	return outS16(frame, length, binary.LittleEndian.Uint16)
}
func outS16BE1(frame []byte) ([]byte, int16) {
	f, u := outS16BE(frame, 1)
	return f, u[0]
}
func outS16BE(frame []byte, length int) ([]byte, []int16) {
	return outS16(frame, length, binary.BigEndian.Uint16)
}
func outS16(frame []byte, length int, f func([]byte) uint16) ([]byte, []int16) {
	size := 2
	output := make([]int16, length)
	for index := 0; index < length*size; index += size {
		output[index/size] = int16(f(frame[index : index+size]))
	}
	return frame[length*size:], output
}

func inS16LE(frame []byte, params ...int16) []byte {
	return inS16(frame, binary.LittleEndian.PutUint16, params...)
}
func inS16BE(frame []byte, params ...int16) []byte {
	return inS16(frame, binary.BigEndian.PutUint16, params...)
}
func inS16(frame []byte, f func([]byte, uint16), params ...int16) []byte {
	size := 2
	temp := make([]byte, 0)
	for _, num := range params {
		u8s := make([]byte, size)
		f(u8s, uint16(num))
		temp = append(temp, u8s...)
	}
	return append(frame, temp...)
}

func outU32LE1(frame []byte) ([]byte, uint32) {
	f, u := outU32LE(frame, 1)
	return f, u[0]
}
func outU32LE(frame []byte, length int) ([]byte, []uint32) {
	return outU32(frame, length, binary.LittleEndian.Uint32)
}
func outU32BE1(frame []byte) ([]byte, uint32) {
	f, u := outU32BE(frame, 1)
	return f, u[0]
}
func outU32BE(frame []byte, length int) ([]byte, []uint32) {
	return outU32(frame, length, binary.BigEndian.Uint32)
}
func outU32(frame []byte, length int, f func([]byte) uint32) ([]byte, []uint32) {
	size := 4
	output := make([]uint32, length)
	for index := 0; index < length*size; index += size {
		output[index/size] = f(frame[index : index+size])
	}
	return frame[length*size:], output
}

func inU32LE(frame []byte, params ...uint32) []byte {
	return inU32(frame, binary.LittleEndian.PutUint32, params...)
}
func inU32BE(frame []byte, params ...uint32) []byte {
	return inU32(frame, binary.BigEndian.PutUint32, params...)
}
func inU32(frame []byte, f func([]byte, uint32), params ...uint32) []byte {
	size := 4
	temp := make([]byte, 0)
	for _, num := range params {
		u8s := make([]byte, size)
		f(u8s, num)
		temp = append(temp, u8s...)
	}
	return append(frame, temp...)
}

func outS32LE1(frame []byte) ([]byte, int32) {
	f, u := outS32LE(frame, 1)
	return f, u[0]
}
func outS32LE(frame []byte, length int) ([]byte, []int32) {
	return outS32(frame, length, binary.LittleEndian.Uint32)
}
func outS32BE1(frame []byte) ([]byte, int32) {
	f, u := outS32BE(frame, 1)
	return f, u[0]
}
func outS32BE(frame []byte, length int) ([]byte, []int32) {
	return outS32(frame, length, binary.BigEndian.Uint32)
}
func outS32(frame []byte, length int, f func([]byte) uint32) ([]byte, []int32) {
	size := 4
	output := make([]int32, length)
	for index := 0; index < length*size; index += size {
		output[index/size] = int32(f(frame[index : index+size]))
	}
	return frame[length*size:], output
}

func inS32LE(frame []byte, params ...int32) []byte {
	return inS32(frame, binary.LittleEndian.PutUint32, params...)
}
func inS32BE(frame []byte, params ...int32) []byte {
	return inS32(frame, binary.BigEndian.PutUint32, params...)
}
func inS32(frame []byte, f func([]byte, uint32), params ...int32) []byte {
	size := 4
	temp := make([]byte, 0)
	for _, num := range params {
		u8s := make([]byte, size)
		f(u8s, uint32(num))
		temp = append(temp, u8s...)
	}
	return append(frame, temp...)
}
func outStrZ(frame []byte) ([]byte, string) {
	var str strings.Builder
	var length = 0
	for {
		c := frame[length]
		if c == 0 {
			length++
			break
		}
		str.WriteByte(c)
		length++
	}
	return frame[length:], str.String()
}
func inStrZ(frame []byte, str string) []byte {
	return append(frame, []byte(str)...)

}
func outStrF(frame []byte, length int) ([]byte, string) {

	return frame[length:], string(frame[:length])
}
func inStrF(frame []byte, str string, length int) []byte {
	return append(frame, []byte(fmt.Sprintf(fmt.Sprintf("%%%ds", length), str))...)
}
func outDateTime(frame []byte) ([]byte, time.Time) {
	// size := 4
	frame, seconds := outU32LE(frame, 1)
	return frame, time.Unix(int64(seconds[0]), 0)
}
func inDateTime(frame []byte, time time.Time) []byte {
	return inU32LE(frame, uint32(time.Unix()))
}
