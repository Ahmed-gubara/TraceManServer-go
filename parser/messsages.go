package parser

import "time"

var messageTypes = map[uint16]string{
	0x1001: "Login",
	0x9001: "Login response",
	0x1002: "Logout",
	0x1003: "Heartbeat",
	0x9003: "Heartbeat response",
	0x2001: "Set",
	0x2002: "Query",
	0xA001: "Setting response",
	0xA002: "Query response",
	0x3001: "Location query",
	0xB001: "Current location information",
	0x3002: "Clear DTC",
	0xB002: "Clear DTC response",
	0x3003: "Restore factory settings",
	0xB003: "Restore factory settings response",
	0x3004: "Remot controlling",
	0xB004: "Remote controlling response",
	0x3005: "Remot listening",
	0xB005: "Remote controlling listening response",
	0x3006: "Text information",
	0xB006: "Text information response",
	0x3007: "Remote querying",
	0xB007: "Remote querying response",
	0x3008: "Tak e photos",
	0xB008: "Take photos response",
	0x300F: "OBD alpha",
	0xB00F: "OBD alpha response",
	0x4001: "GPS data",
	0x4002: "PID data",
	0x4003: "G-Sensor data",
	0x4004: "Supported PID types",
	0x4005: "Snapshot/Frozen frame data",
	0x4006: "DTCs of Passenger car",
	0x400B: "DTCs of commercial vehicle",
	0x4007: "Alarm",
	0xC007: "Alarm received confirmation",
	0x4008: "Cell ID",
	0x4009: "GPS report in sleep",
	0x400C: "Driver card ID",
	0xC00C: "Driver card ID received confirmation",
	0x400D: "RFID card number and location information uploading",
	0xC00D: "RFID card number rand location information confirmation",
	0x400E: "fuel data uploading",
	0x400F: "Crash G-Sensor	Data upload",
	0x5001: "Update	notification",
	0xD001: "Update confirmation",
	0x5002: "Update message",
	0xD002: "Update message confirmation",
	0x5101: "A-GPS data request",
	0x5102: "A-GPS message",
	0xD102: "A-GPS message confirmation",
	0x6001: "RSA public key	request",
	0xF001: "RSA public key distribute",
	0x6002: "AES key uploading",
	0xF002: "AES key confirmation",
}

func GetMessageType(ID uint16) string {
	val, found := messageTypes[ID]
	if found {
		return val
	}
	return "UNKNOWN MESSAGE TYPE"
}

type Login0x1001 struct {
	StatData          StatData
	GPSData           GPSData
	SoftwareVersion   string `string:"strz"`
	HardwareVersion   string `string:"strz"`
	NewParameterCount uint16
	NewParameterArray []uint16 `length:"NewParameterCount"`
}
type LoginResponse0x9001 struct {
	IPAddress  [4]uint8
	Port       uint16
	ServerTime time.Time
}
type Logout0x1002 struct {
}

type Heartbeat0x1003 struct {
}
type HearbeatResponse0x9003 struct {
}

type GPSinSleep0x4009 struct {
	UTCtime time.Time
	GPSItem GPSItem
}
