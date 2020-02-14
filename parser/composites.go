package parser

import (
	"time"
)

//Date ...
type Date struct {
	Day   uint8
	Month uint8
	Year  uint8
}

//Time ...
type Time struct {
	Hour   uint8
	Minute uint8
	Second uint8
}

//GPSItem ..
type GPSItem struct {
	Date      Date
	Time      Time
	Latitude  uint32
	Longitude uint32
	Speed     uint16
	Direction uint16
	Flag      uint8
}

type GPSData struct {
	GPSCount uint8
	GPSArray []GPSItem `length:"GPSCount"`
}

type VState struct {
	S0 uint8
	S1 uint8
	S2 uint8
	S3 uint8
	// ExhaustEmission              bool `depend:"S0,7"`
	// IdleEngine                   bool `depend:"S0,6"`
	// HardDeceleration             bool `depend:"S0,5"`
	// HardAcceleration             bool `depend:"S0,4"`
	// HighEngineCoolantTemperature bool `depend:"S0,3"`
	// Speeding                     bool `depend:"S0,2"`
	// Towing                       bool `depend:"S0,1"`
	// LowVoltage                   bool `depend:"S0,0"`
	// Tamper                       bool `depend:"S1,7"`
	// Crash                        bool `depend:"S1,6"`
	// Emergency                    bool `depend:"S1,5"`
	// FatigueDriving               bool `depend:"S1,4"`
	// SharpTurn                    bool `depend:"S1,3"`
	// QuickLaneChange              bool `depend:"S1,2"`
	// PowerOn                      bool `depend:"S1,1"`
	// HighRPM                      bool `depend:"S1,0"`
	// Mil                          bool `depend:"S2,7"`
	// ObdCommunicationError        bool `depend:"S2,6"`
	// PowerOff                     bool `depend:"S2,5"`
	// NoGPSdevice                  bool `depend:"S2,4"`
	// PrivacyStatus                bool `depend:"S2,3"`
	// IgnitionOn                   bool `depend:"S2,2"`
	// IllegalIgnition              bool `depend:"S2,1"`
	// IllegalEnter                 bool `depend:"S2,0"`
	// Reserved1                    bool `depend:"S3,7"`
	// Reserved2                    bool `depend:"S3,6"`
	// Door2Stat                    bool `depend:"S3,5"`
	// Door1Stat                    bool `depend:"S3,4"`
	// Vibration                    bool `depend:"S3,3"`
	// DangerousDriving             bool `depend:"S3,2"`
	// NoCardPresented              bool `depend:"S3,1"`
	// Unlock                       bool `depend:"S3,0"`
}

type StatData struct {
	LastAccOnTime    time.Time
	UtcTime          time.Time
	TotalTripMilage  uint32
	CurentTripMilage uint32
	TotalFuel        uint32
	CurrentFuel      uint16
	Vstate           VState
	Reserved         [8]uint8
}

var boolBit = [8]uint8{1, 2, 4, 8, 16, 32, 64, 128}

type ProtocolPrefix struct {
	ProtocolHead    [2]uint8 //2
	ProtocolLength  uint16
	ProtocolVersion uint8
	DeviceID        string `string:"strf,20"`
	ProtocolID      uint16 `binary:"BE"`
}

func (f ProtocolPrefix) inFrameSize() int {
	return 2 + 2 + 1 + 20 + 2
}

// func GetProtocolPrefix(frame []byte) ([]byte, ProtocolPrefix) {
// 	frame, protocolHead := outU8(frame, 2)
// 	frame, protocolLength := outU16LE1(frame)
// 	frame, protocolVersion := outU8_1(frame)
// 	frame, deviceID := outStrF(frame, 20)
// 	frame, protocolID := outU16BE1(frame)
// 	return frame, ProtocolPrefix{
// 		ProtocolHead:    protocolHead,
// 		ProtocolLength:  protocolLength,
// 		ProtocolVersion: protocolVersion,
// 		DeviceID:        deviceID,
// 		ProtocolID:      protocolID,
// 	}
// }
// func SetProtocolPrefix(frame []byte, pre ProtocolPrefix) []byte {
// 	frame = inU8(frame, pre.ProtocolHead...)
// 	frame = inU16LE(frame, pre.ProtocolLength)
// 	frame = inU8(frame, pre.ProtocolVersion)
// 	frame = inStrF(frame, pre.DeviceID, 20)
// 	frame = inU16BE(frame, pre.ProtocolID)
// 	return frame
// }

type ProtocolSufix struct {
	crc          uint16
	protocolTail []uint8 // 2
}

func (f ProtocolSufix) inFrameSize() int {
	return 2 + 2
}
func GetProtocolSufix(frame []byte) ([]byte, ProtocolSufix) {
	frame, crc := outU16LE1(frame)
	frame, protocolTail := outU8(frame, 2)
	return frame, ProtocolSufix{crc: crc, protocolTail: protocolTail}
}
func SetProtocolSufix(frame []byte, suf ProtocolSufix) []byte {
	frame = inU16LE(frame, suf.crc)
	frame = inU8(frame, suf.protocolTail...)
	return frame
}
