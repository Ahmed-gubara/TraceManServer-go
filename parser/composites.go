package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	S0                           uint8
	S1                           uint8
	S2                           uint8
	S3                           uint8
	ExhaustEmission              bool `depend:"S0,7"`
	IdleEngine                   bool `depend:"S0,6"`
	HardDeceleration             bool `depend:"S0,5"`
	HardAcceleration             bool `depend:"S0,4"`
	HighEngineCoolantTemperature bool `depend:"S0,3"`
	Speeding                     bool `depend:"S0,2"`
	Towing                       bool `depend:"S0,1"`
	LowVoltage                   bool `depend:"S0,0"`
	Tamper                       bool `depend:"S1,7"`
	Crash                        bool `depend:"S1,6"`
	Emergency                    bool `depend:"S1,5"`
	FatigueDriving               bool `depend:"S1,4"`
	SharpTurn                    bool `depend:"S1,3"`
	QuickLaneChange              bool `depend:"S1,2"`
	PowerOn                      bool `depend:"S1,1"`
	HighRPM                      bool `depend:"S1,0"`
	Mil                          bool `depend:"S2,7"`
	ObdCommunicationError        bool `depend:"S2,6"`
	PowerOff                     bool `depend:"S2,5"`
	NoGPSdevice                  bool `depend:"S2,4"`
	PrivacyStatus                bool `depend:"S2,3"`
	IgnitionOn                   bool `depend:"S2,2"`
	IllegalIgnition              bool `depend:"S2,1"`
	IllegalEnter                 bool `depend:"S2,0"`
	Reserved1                    bool `depend:"S3,7"`
	Reserved2                    bool `depend:"S3,6"`
	Door2Stat                    bool `depend:"S3,5"`
	Door1Stat                    bool `depend:"S3,4"`
	Vibration                    bool `depend:"S3,3"`
	DangerousDriving             bool `depend:"S3,2"`
	NoCardPresented              bool `depend:"S3,1"`
	Unlock                       bool `depend:"S3,0"`
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

func GetPayload(frame []byte, rType interface{}) ([]byte, interface{}) {
	frame, interfaceValue := getPayload(frame, reflect.TypeOf(rType), "", nil)
	return frame, interfaceValue.Interface()
}
func getPayload(frame []byte, rType reflect.Type, tag reflect.StructTag, parent *reflect.Value) ([]byte, reflect.Value) {
	switch rType.Kind() {
	case reflect.Bool:
		if dep := tag.Get("depend"); len(dep) > 0 {
			if specs := strings.Split(dep, ","); len(specs) == 2 {
				value, _ := parent.FieldByName(specs[0]).Interface().(uint8)
				bitindex, _ := strconv.Atoi(specs[1])
				val := (value & boolBit[bitindex]) == boolBit[bitindex]
				return frame, reflect.ValueOf(val)
			}
		}
	case reflect.String:
		strType := tag.Get("string")
		switch {
		case strType == "strz":
			{
				frame, str := outStrZ(frame)
				return frame, reflect.ValueOf(str)
			}
		case strings.HasPrefix(strType, "strf"):
			{
				splits := strings.Split(strType, ",")
				length, err := strconv.Atoi(splits[1])
				if err != nil {
					panic("no length specificed")
				}
				frame, str := outStrF(frame, length)
				return frame, reflect.ValueOf(str)
			}
		case strings.HasPrefix(strType, "str"):
			{
				splits := strings.Split(strType, ",")
				length := parent.FieldByName(splits[1]).Interface().(int)
				frame, str := outStrF(frame, length)
				return frame, reflect.ValueOf(str)
			}
		default:
			panic(fmt.Sprintf("string with no defined length %s.%s", parent.Type().Name(), rType.Name()))
		}
	case reflect.Uint32:
		frame, val := outU32LE1(frame)
		return frame, reflect.ValueOf(val)
	case reflect.Uint16:
		frame, val := outU16LE1(frame)
		return frame, reflect.ValueOf(val)
	case reflect.Uint8:
		frame, val := outU8_1(frame)
		return frame, reflect.ValueOf(val)
	case reflect.Array:
		println(rType.Elem().String())
		arrType := reflect.ArrayOf(rType.Len(), rType.Elem())
		arr := reflect.New(arrType).Elem()
		for i := 0; i < rType.Len(); i++ {
			print("list " + rType.Elem().Name())
			frame1, val := getPayload(frame, rType.Elem(), tag, parent)
			frame = frame1
			arr.Index(i).Set(val)
		}
		return frame, arr
	case reflect.Slice:

		if lengthfield, ok := tag.Lookup("length"); !ok {
			panic(fmt.Sprintf("slice with no defined length %s.%s", parent.Type().Name(), rType.Name()))
		} else {
			length, err := strconv.Atoi(fmt.Sprint(parent.FieldByName(lengthfield).Interface()))
			if err != nil {
				panic(err)
			}
			sliceType := reflect.SliceOf(rType.Elem())
			slice := reflect.New(sliceType).Elem()
			for i := 0; i < length; i++ {
				frame1, val := getPayload(frame, rType.Elem(), tag, parent)
				frame = frame1
				slice = reflect.Append(slice, val)
			}
			return frame, slice
		}

	case reflect.Struct:
		{
			switch {
			case rType.String() == "time.Time":
				frame1, time := outDateTime(frame)
				frame = frame1
				return frame, reflect.ValueOf(time)
			case strings.HasPrefix(rType.String(), "parser."):
				newStuct := reflect.New(rType).Elem()
				println(rType.String())
				for i := 0; i < rType.NumField(); i++ {
					field := rType.Field(i)
					if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
						panic(fmt.Sprintf("field is private %s.%s", rType.Name(), field.Name))
					}
					frame1, vstate := getPayload(frame, field.Type, field.Tag, &newStuct)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(vstate.Interface()))
				}
				return frame, newStuct
			default:
				panic(fmt.Sprintf("struct of type %s can't be parsed", rType.String()))
			}
		}
	default:
		panic(fmt.Sprintf("Kind not registered %s", rType.Kind()))
	}
	panic("this PANIC shouldn't be reached, this indicates a BUG somewhere!")
}

type ProtocolPrefix struct {
	ProtocolHead    [2]uint8 //2
	ProtocolLength  uint16
	ProtocolVersion uint8
	DeviceID        string `string:"strf,20"`
	ProtocolID      uint16
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
func SetPayload(frame []byte, rType interface{}) []byte {
	return setPayload(frame, reflect.ValueOf(rType), "", nil)

}
func Encapsulate(protocolversion uint8, deviceID string, protocolID uint16, payload interface{}) []byte {
	pload := SetPayload([]byte{}, payload)
	var protocolLength uint16 = uint16(len(pload)) + 2 + 2 + 1 + 20 + 2 + 2 + 2
	prefix := ProtocolPrefix{ProtocolHead: [2]byte{0x40, 0x40}, ProtocolLength: protocolLength, ProtocolVersion: protocolversion, DeviceID: deviceID, ProtocolID: protocolID}
	frame := SetPayload([]byte{}, prefix)
	frame = append(frame, pload...)
	crc := CRC_MakeCrc(frame)
	frame = inU16LE(frame, crc)
	frame = inU8(frame, 0x0d, 0x0a)
	return frame
}
func setPayload(frame []byte, rValue reflect.Value, tag reflect.StructTag, parent *reflect.Value) []byte {
	switch rValue.Kind() {
	case reflect.Bool:
		panic("logic not ready,can't be set")
		// if dep := tag.Get("depend"); len(dep) > 0 {
		// 	if specs := strings.Split(dep, ","); len(specs) == 2 {
		// 		value, _ := parent.FieldByName(specs[0]).Interface().(uint8)
		// 		bitindex, _ := strconv.Atoi(specs[1])
		// 		val := (value & boolBit[bitindex]) == boolBit[bitindex]
		// 		return frame, reflect.ValueOf(val)
		// 	}
		// }
	case reflect.String:
		strType := tag.Get("string")
		switch {
		case strType == "strz":
			{
				frame = inStrZ(frame, rValue.Interface().(string))
				return frame
			}
		case strings.HasPrefix(strType, "strf"):
			{
				splits := strings.Split(strType, ",")
				length, err := strconv.Atoi(splits[1])
				if err != nil {
					panic("no length specificed")
				}
				frame = inStrF(frame, rValue.Interface().(string), length)
				return frame
			}
		case strings.HasPrefix(strType, "str"):
			{
				splits := strings.Split(strType, ",")
				str := rValue.Interface().(string)
				parent.FieldByName(splits[1]).Set(reflect.ValueOf(len(str)))
				frame = inStrF(frame, str, len(str))
				return frame
			}
		default:
			panic(fmt.Sprintf("string with no defined length %s.%s", parent.Type().Name(), rValue.Type().Name()))
		}
	case reflect.Uint32:
		frame = inU32LE(frame, rValue.Interface().(uint32))
		return frame
	case reflect.Uint16:
		frame = inU16LE(frame, rValue.Interface().(uint16))
		return frame
	case reflect.Uint8:
		frame = inU8(frame, rValue.Interface().(uint8))
		return frame
	case reflect.Array:
		for i := 0; i < rValue.Len(); i++ {

			frame1 := setPayload(frame, rValue.Index(i), tag, parent)
			frame = frame1

		}
		return frame
	case reflect.Slice:

		if lengthfield, ok := tag.Lookup("length"); !ok {
			panic(fmt.Sprintf("slice with no defined length %s.%s", parent.Type().Name(), rValue.Type().Name()))
		} else {
			parent.FieldByName(lengthfield).Set(reflect.ValueOf(rValue.Len()))

			for i := 0; i < rValue.Len(); i++ {
				frame1 := setPayload(frame, rValue.Index(i), tag, parent)
				frame = frame1
			}
			return frame
		}

	case reflect.Struct:
		{
			switch {
			case rValue.Type().String() == "time.Time":
				frame1 := inDateTime(frame, rValue.Interface().(time.Time))
				frame = frame1
				return frame
			case strings.HasPrefix(rValue.Type().String(), "parser."):

				for i := 0; i < rValue.NumField(); i++ {

					field := rValue.Field(i)
					// if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
					// 	panic(fmt.Sprintf("field is private %s.%s", rValue.Name(), field.Name))
					// }
					frame1 := setPayload(frame, field, rValue.Type().Field(i).Tag, &rValue)
					frame = frame1

				}
				return frame
			default:
				panic(fmt.Sprintf("struct of type %s can't be parsed", rValue.Type().String()))
			}
		}
	default:
		panic(fmt.Sprintf("Kind not registered %s", rValue.Kind()))
	}
	panic("this PANIC shouldn't be reached, this indicates a BUG somewhere!")
}
