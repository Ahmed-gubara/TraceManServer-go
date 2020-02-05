package parser

import "time"

type LoginMessage0x1001 struct {
	StatData          StatData
	GPSData           GPSData
	SoftwareVersion   string `string:"strz"`
	HardwareVersion   string `string:"strz"`
	NewParameterCount uint16
	NewParameterArray []uint16 `length:"NewParameterCount"`
}
type LoginResponseMessage0x9001 struct {
	IPAddress  [4]uint8
	Port       uint16
	ServerTime time.Time
}
