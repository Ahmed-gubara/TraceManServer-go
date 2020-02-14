package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetDate(frame []byte) ([]byte, Date) {
	frame, day := outU8_1(frame)
	frame, month := outU8_1(frame)
	frame, year := outU8_1(frame)
	return frame, Date{Day: day, Month: month, Year: year}
}
func SetDate(frame []byte, date Date) []byte {
	frame = inU8(frame, date.Day)
	frame = inU8(frame, date.Month)
	frame = inU8(frame, date.Year)
	return frame
}
func GetTime(frame []byte) ([]byte, Time) {
	frame, hour := outU8_1(frame)
	frame, minute := outU8_1(frame)
	frame, second := outU8_1(frame)
	return frame, Time{Hour: hour, Minute: minute, Second: second}
}
func SetTime(frame []byte, time Time) []byte {
	frame = inU8(frame, time.Hour)
	frame = inU8(frame, time.Minute)
	frame = inU8(frame, time.Second)
	return frame
}
func GetGPSItem(frame []byte) ([]byte, GPSItem) {
	frame, date := GetDate(frame)
	frame, time := GetTime(frame)
	frame, latitude := outU32LE1(frame)
	frame, longitude := outU32LE1(frame)
	frame, speed := outU16LE1(frame)
	frame, direction := outU16LE1(frame)
	frame, flag := outU8_1(frame)
	return frame, GPSItem{Date: date, Time: time, Latitude: latitude, Longitude: longitude, Speed: speed, Direction: direction, Flag: flag}
}
func getGPSData(frame []byte) ([]byte, GPSData) {
	frame, gpsCount := outU8_1(frame)
	gpsArray := make([]GPSItem, gpsCount)
	for index := 0; index < int(gpsCount); index++ {
		frame, gpsArray[index] = GetGPSItem(frame)
	}
	return frame, GPSData{GPSCount: gpsCount, GPSArray: gpsArray}
}
func GetVState(frame []byte) ([]byte, VState) {
	// bit := []uint8{1, 2, 4, 8, 16, 32, 64, 128}
	// frame, s0 := outU8_1(frame)
	// frame, s1 := outU8_1(frame)
	// frame, s2 := outU8_1(frame)
	// frame, s3 := outU8_1(frame)
	return frame, VState{
		// ExhaustEmission:              s0&bit[7] == 1, //s0
		// IdleEngine:                   s0&bit[6] == 1,
		// HardDeceleration:             s0&bit[5] == 1,
		// HardAcceleration:             s0&bit[4] == 1,
		// HighEngineCoolantTemperature: s0&bit[3] == 1,
		// Speeding:                     s0&bit[2] == 1,
		// Towing:                       s0&bit[1] == 1,
		// LowVoltage:                   s0&bit[0] == 1,
		// Tamper:                       s1&bit[7] == 1, //s1
		// Crash:                        s1&bit[6] == 1,
		// Emergency:                    s1&bit[5] == 1,
		// FatigueDriving:               s1&bit[4] == 1,
		// SharpTurn:                    s1&bit[3] == 1,
		// QuickLaneChange:              s1&bit[2] == 1,
		// PowerOn:                      s1&bit[1] == 1,
		// HighRPM:                      s1&bit[0] == 1,
		// Mil:                          s2&bit[7] == 1, //s2
		// ObdCommunicationError:        s2&bit[6] == 1,
		// PowerOff:                     s2&bit[5] == 1,
		// NoGPSdevice:                  s2&bit[4] == 1,
		// PrivacyStatus:                s2&bit[3] == 1,
		// IgnitionOn:                   s2&bit[2] == 1,
		// IllegalIgnition:              s2&bit[1] == 1,
		// IllegalEnter:                 s2&bit[0] == 1,
		// Reserved1:                    s3&bit[7] == 1, //s3
		// Reserved2:                    s3&bit[6] == 1,
		// Door2Stat:                    s3&bit[5] == 1,
		// Door1Stat:                    s3&bit[4] == 1,
		// Vibration:                    s3&bit[3] == 1,
		// DangerousDriving:             s3&bit[2] == 1,
		// NoCardPresented:              s3&bit[1] == 1,
		// Unlock:                       s3&bit[0] == 1,
	}
}
func GetStatData(frame []byte) ([]byte, StatData) {
	frame, lastAccOnTime := outDateTime(frame)
	frame, utcTime := outDateTime(frame)
	frame, totalTripMilage := outU32LE1(frame)
	frame, curentTripMilage := outU32LE1(frame)
	frame, totalFuel := outU32LE1(frame)
	frame, currentFuel := outU16LE1(frame)
	frame, vstate := GetVState(frame)
	frame, _ = outU8(frame, 8)
	return frame, StatData{
		LastAccOnTime:    lastAccOnTime,
		UtcTime:          utcTime,
		TotalTripMilage:  totalTripMilage,
		CurentTripMilage: curentTripMilage,
		TotalFuel:        totalFuel,
		CurrentFuel:      currentFuel,
		Vstate:           vstate,
		// reserved:         reserved,
	}
}
func OLDGetPayload(frame []byte, rType reflect.Type) ([]byte, reflect.Value) {

	newStuct := reflect.New(rType).Elem()
	switch rType.Kind() {

	case reflect.Struct:
		{

			for i := 0; i < rType.NumField(); i++ {
				field := rType.Field(i)
				if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
					panic(fmt.Sprintf("field is private %s.%s", rType.Name(), field.Name))
				}
				switch field.Type.Kind() {
				case reflect.Bool:
					if dep := field.Tag.Get("depend"); len(dep) > 0 {
						if specs := strings.Split(dep, ","); len(specs) == 2 {
							value, _ := newStuct.FieldByName(specs[0]).Interface().(uint8)
							bitindex, _ := strconv.Atoi(specs[1])
							val := (value & boolBit[bitindex]) == boolBit[bitindex]
							println(val)
							newStuct.Field(i).Set(reflect.ValueOf(val))
						}
					}
				case reflect.Uint32:
					frame1, val := outU32LE1(frame)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(val))
					break
				case reflect.Uint16:
					frame1, val := outU16LE1(frame)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(val))
					break
				case reflect.Uint8:
					frame1, val := outU8_1(frame)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(val))
					break
				case reflect.Array:

					frame1, val := outU8(frame, field.Type.Len())
					var arr [8]byte
					copy(arr[:], val)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(arr))
					break
				case reflect.Struct:
					{
						switch {
						case field.Type.String() == "time.Time":
							frame1, time := outDateTime(frame)
							frame = frame1
							newStuct.Field(i).Set(reflect.ValueOf(time))
							break
						case strings.HasPrefix(field.Type.String(), "parser."):
							frame1, vstate := OLDGetPayload(frame, field.Type)
							frame = frame1
							newStuct.Field(i).Set(reflect.ValueOf(vstate.Interface()))
							break
						default:
							panic(fmt.Sprintf("struct of type %s can't be parsed", field.Type.String()))
						}
					}
				default:
					panic(fmt.Sprintf("Kind not registered %s", field.Type.Kind()))
				}
			}
		}

	}

	// fmt.Printf("%+v", stype.Elem().Kind())
	return frame, newStuct
}
