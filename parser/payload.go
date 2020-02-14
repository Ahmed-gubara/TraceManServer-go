package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var indent int = 0

func Echo(format string, a ...interface{}) {
	return
	fmt.Printf(fmt.Sprintf(fmt.Sprintf("%%%ds", indent), "")+format, a...)
}
func Encapsulate(protocolversion uint8, deviceID string, protocolID uint16, payload interface{}) []byte {

	bpayload := SetPayload([]byte{}, payload)
	var protocolLength uint16 = uint16(len(bpayload)) + 2 + 2 + 1 + 20 + 2 + 2 + 2
	prefix := ProtocolPrefix{ProtocolHead: [2]byte{0x40, 0x40}, ProtocolLength: protocolLength, ProtocolVersion: protocolversion, DeviceID: deviceID, ProtocolID: protocolID}
	frame := make([]byte, 0)
	frame = SetPayload(frame, prefix)
	frame = append(frame, bpayload...)

	crc := CRC_MakeCrc(frame)
	frame = inU16LE(frame, crc)
	frame = inU8(frame, 0x0d, 0x0a)
	return frame
}

func SetPayload(frame []byte, rType interface{}) []byte {
	return setPayload(frame, reflect.ValueOf(rType), "", nil)

}

func setPayload(frame []byte, rValue reflect.Value, tag reflect.StructTag, parent *reflect.Value) []byte {
	framesize := len(frame)
	Echo("{")
	defer func() {
		Echo("}added type %v , with tag %v , frame increase by %v \n", rValue.Type().Name(), tag, (len(frame) - framesize))
	}()
	switch rValue.Kind() {
	case reflect.Bool:
		// return frame
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
				Echo("before : %d\n", len(frame))

				splits := strings.Split(strType, ",")
				length, err := strconv.Atoi(splits[1])

				if err != nil {
					panic("no length specificed")
				}
				frame = inStrF(frame, rValue.Interface().(string), length)
				Echo("after : %d\n", len(frame))

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
	case reflect.Int32:
		frame = inS32LE(frame, rValue.Interface().(int32))
		return frame
	case reflect.Int16:
		frame = inS16LE(frame, rValue.Interface().(int16))
		return frame
	case reflect.Uint16:
		if v, f := tag.Lookup("binary"); f && strings.ToUpper(v) == "BE" {
			frame = inU16BE(frame, rValue.Interface().(uint16))
			return frame
		}
		frame = inU16LE(frame, rValue.Interface().(uint16))
		return frame
	case reflect.Int8:
		frame = inS8(frame, rValue.Interface().(int8))
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

		if _, ok := tag.Lookup("length"); !ok {
			panic(fmt.Sprintf("slice with no defined length %s.%s", parent.Type().Name(), rValue.Type().Name()))
		} else {
			//fmt.Printf("%s ..........", parent.FieldByName(lengthfield))
			//parent.FieldByName(lengthfield).Set(reflect.ValueOf(uint8(rValue.Len())))
			indent += 2

			for i := 0; i < rValue.Len(); i++ {
				frame1 := setPayload(frame, rValue.Index(i), tag, parent)
				frame = frame1
			}
			indent -= 2
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
				indent += 2

				for i := 0; i < rValue.NumField(); i++ {

					field := rValue.Field(i)
					// if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
					// 	panic(fmt.Sprintf("field is private %s.%s", rValue.Name(), field.Name))
					// }
					frame1 := setPayload(frame, field, rValue.Type().Field(i).Tag, &rValue)
					frame = frame1

				}
				indent -= 2
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

func GetPayload(frame []byte, rType interface{}) ([]byte, interface{}) {
	frame, interfaceValue := getPayload(frame, reflect.TypeOf(rType), "", nil)
	return frame, interfaceValue.Interface()
}
func getPayload(frame []byte, rType reflect.Type, tag reflect.StructTag, parent *reflect.Value) ([]byte, reflect.Value) {
	framesize := len(frame)
	Echo("{")
	defer func() {
		Echo("}getted type %v , with tag %v , frame decrease by %v \n", rType.Name(), tag, (len(frame) - framesize))
	}()
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
		if v, f := tag.Lookup("binary"); f && strings.ToUpper(v) == "BE" {
			frame, val := outU16BE1(frame)
			return frame, reflect.ValueOf(val)
		}
		frame, val := outU16LE1(frame)
		return frame, reflect.ValueOf(val)

	case reflect.Uint8:
		frame, val := outU8_1(frame)
		return frame, reflect.ValueOf(val)
	case reflect.Array:
		arrType := reflect.ArrayOf(rType.Len(), rType.Elem())
		arr := reflect.New(arrType).Elem()
		for i := 0; i < rType.Len(); i++ {
			Echo("list " + rType.Elem().Name())
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
			indent += 2
			for i := 0; i < length; i++ {
				frame1, val := getPayload(frame, rType.Elem(), tag, parent)
				frame = frame1
				slice = reflect.Append(slice, val)
			}
			indent -= 2
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
				Echo(rType.String())
				indent += 2
				for i := 0; i < rType.NumField(); i++ {
					field := rType.Field(i)
					if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
						panic(fmt.Sprintf("field is private %s.%s", rType.Name(), field.Name))
					}
					frame1, vstate := getPayload(frame, field.Type, field.Tag, &newStuct)
					frame = frame1
					newStuct.Field(i).Set(reflect.ValueOf(vstate.Interface()))
				}
				indent -= 2
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
