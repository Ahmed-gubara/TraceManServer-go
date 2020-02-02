package parser

import "time"

import "fmt"

//Date ...
type Date struct {
	day   uint8
	month uint8
	year  uint8
}

func GetDate(frame []byte) ([]byte, Date) {
	frame, day := outU8_1(frame)
	frame, month := outU8_1(frame)
	frame, year := outU8_1(frame)
	return frame, Date{day: day, month: month, year: year}
}
func SetDate(frame []byte, date Date) []byte {
	frame = inU8(frame, date.day)
	frame = inU8(frame, date.month)
	frame = inU8(frame, date.year)
	return frame
}

//Time ...
type Time struct {
	hour   uint8
	minute uint8
	second uint8
}

func GetTime(frame []byte) ([]byte, Time) {
	frame, hour := outU8_1(frame)
	frame, minute := outU8_1(frame)
	frame, second := outU8_1(frame)
	return frame, Time{hour: hour, minute: minute, second: second}
}
func SetTime(frame []byte, time Time) []byte {
	frame = inU8(frame, time.hour)
	frame = inU8(frame, time.minute)
	frame = inU8(frame, time.second)
	return frame
}

//GPSItem ..
type GPSItem struct {
	date      Date
	time      Time
	latitude  uint32
	longitude uint32
	speed     uint16
	direction uint16
	flag      uint8
}

func GetGPSItem(frame []byte) ([]byte, GPSItem) {
	frame, date := GetDate(frame)
	frame, time := GetTime(frame)
	frame, latitude := outU32LE1(frame)
	frame, longitude := outU32LE1(frame)
	frame, speed := outU16LE1(frame)
	frame, direction := outU16LE1(frame)
	frame, flag := outU8_1(frame)
	return frame, GPSItem{date: date, time: time, latitude: latitude, longitude: longitude, speed: speed, direction: direction, flag: flag}
}

type GPSData struct {
	gpsCount uint8
	gpsArray []GPSItem
}

func getGPSData(frame []byte) ([]byte, GPSData) {
	frame, gpsCount := outU8_1(frame)
	gpsArray := make([]GPSItem, gpsCount)
	for index := 0; index < int(gpsCount); index++ {
		frame, gpsArray[index] = GetGPSItem(frame)
	}
	return frame, GPSData{gpsCount: gpsCount, gpsArray: gpsArray}
}

type VState struct {
	exhaustEmission              bool
	idleEngine                   bool
	hardDeceleration             bool
	hardAcceleration             bool
	highEngineCoolantTemperature bool
	speeding                     bool
	towing                       bool
	lowVoltage                   bool
	tamper                       bool
	crash                        bool
	emergency                    bool
	fatigueDriving               bool
	sharpTurn                    bool
	quickLaneChange              bool
	powerOn                      bool
	highRPM                      bool
	mil                          bool
	obdCommunicationError        bool
	powerOff                     bool
	noGPSdevice                  bool
	privacyStatus                bool
	ignitionOn                   bool
	illegalIgnition              bool
	IllegalEnter                 bool
	reserved1                    bool
	reserved2                    bool
	door2Stat                    bool
	door1Stat                    bool
	vibration                    bool
	dangerousDriving             bool
	noCardPresented              bool
	unlock                       bool
}

func GetVState(frame []byte) ([]byte, VState) {
	bit := []uint8{1, 2, 4, 8, 16, 32, 64}
	frame, s0 := outU8_1(frame)
	frame, s1 := outU8_1(frame)
	frame, s2 := outU8_1(frame)
	frame, s3 := outU8_1(frame)
	return frame, VState{
		exhaustEmission:              s0&bit[7] == 1, //s0
		idleEngine:                   s0&bit[6] == 1,
		hardDeceleration:             s0&bit[5] == 1,
		hardAcceleration:             s0&bit[4] == 1,
		highEngineCoolantTemperature: s0&bit[3] == 1,
		speeding:                     s0&bit[2] == 1,
		towing:                       s0&bit[1] == 1,
		lowVoltage:                   s0&bit[0] == 1,
		tamper:                       s1&bit[7] == 1, //s1
		crash:                        s1&bit[6] == 1,
		emergency:                    s1&bit[5] == 1,
		fatigueDriving:               s1&bit[4] == 1,
		sharpTurn:                    s1&bit[3] == 1,
		quickLaneChange:              s1&bit[2] == 1,
		powerOn:                      s1&bit[1] == 1,
		highRPM:                      s1&bit[0] == 1,
		mil:                          s2&bit[7] == 1, //s2
		obdCommunicationError:        s2&bit[6] == 1,
		powerOff:                     s2&bit[5] == 1,
		noGPSdevice:                  s2&bit[4] == 1,
		privacyStatus:                s2&bit[3] == 1,
		ignitionOn:                   s2&bit[2] == 1,
		illegalIgnition:              s2&bit[1] == 1,
		IllegalEnter:                 s2&bit[0] == 1,
		reserved1:                    s3&bit[7] == 1, //s3
		reserved2:                    s3&bit[6] == 1,
		door2Stat:                    s3&bit[5] == 1,
		door1Stat:                    s3&bit[4] == 1,
		vibration:                    s3&bit[3] == 1,
		dangerousDriving:             s3&bit[2] == 1,
		noCardPresented:              s3&bit[1] == 1,
		unlock:                       s3&bit[0] == 1,
	}
}

type StatData struct {
	lastAccOnTime    time.Time
	utcTime          time.Time
	totalTripMilage  uint32
	curentTripMilage uint32
	totalFuel        uint32
	currentFuel      uint16
	vstate           VState
	reserved         []uint8 // should be fixed to 8 length
}

func GetStatData(frame []byte) ([]byte, StatData) {
	frame, lastAccOnTime := outDatetime(frame)
	fmt.Println(lastAccOnTime)
	frame, utcTime := outDatetime(frame)
	fmt.Println(utcTime)
	frame, totalTripMilage := outU32LE1(frame)
	frame, curentTripMilage := outU32LE1(frame)
	frame, totalFuel := outU32LE1(frame)
	frame, currentFuel := outU16LE1(frame)
	frame, vstate := GetVState(frame)
	frame, reserved := outU8(frame, 8)
	return frame, StatData{
		lastAccOnTime:    lastAccOnTime,
		utcTime:          utcTime,
		totalTripMilage:  totalTripMilage,
		curentTripMilage: curentTripMilage,
		totalFuel:        totalFuel,
		currentFuel:      currentFuel,
		vstate:           vstate,
		reserved:         reserved,
	}
}

type ProtocolPrefix struct {
	protocolHead    []uint8 //2
	protocolLength  uint16
	protocolVersion uint8
	deviceID        string
	protocolID      uint16
}

func (f ProtocolPrefix) inFrameSize() int {
	return 2 + 2 + 1 + 20 + 2
}
func GetProtocolPrefix(frame []byte) ([]byte, ProtocolPrefix) {
	frame, protocolHead := outU8(frame, 2)
	frame, protocolLength := outU16LE1(frame)
	frame, protocolVersion := outU8_1(frame)
	frame, deviceID := outStrF(frame, 20)
	frame, protocolID := outU16BE1(frame)
	return frame, ProtocolPrefix{
		protocolHead:    protocolHead,
		protocolLength:  protocolLength,
		protocolVersion: protocolVersion,
		deviceID:        deviceID,
		protocolID:      protocolID,
	}
}
func SetProtocolPrefix(frame []byte, pre ProtocolPrefix) []byte {
	frame = inU8(frame, pre.protocolHead...)
	frame = inU16LE(frame, pre.protocolLength)
	frame = inU8(frame, pre.protocolVersion)
	frame = inStrF(frame, pre.deviceID, 20)
	frame = inU16BE(frame, pre.protocolID)
	return frame
}

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
