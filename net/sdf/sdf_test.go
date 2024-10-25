package sdf

type integerCase struct {
	name    string
	integer any
	bin     []byte
}

func integerCases() []integerCase {

	return []integerCase{
		//
		// unsigned integers
		//
		{"uint8::255", uint8(255), []byte{sdtUint8, 255}},
		{"uint16::65535", uint16(65535), []byte{sdtUint16, 255, 255}},
		{"uint32::4294967295", uint32(4294967295), []byte{sdtUint32, 255, 255, 255, 255}},
		{"uint64::18446744073709551615", uint64(18446744073709551615), []byte{sdtUint64, 255, 255, 255, 255, 255, 255, 255, 255}},

		// fails on 32bit arch
		{"uint::18446744073709551615", uint(18446744073709551615), []byte{sdtUint, 255, 255, 255, 255, 255, 255, 255, 255}},

		//
		// signed integers
		//

		{"int8::-127", int8(-127), []byte{sdtInt8, 129}},
		{"int8::127", int8(127), []byte{sdtInt8, 127}},
		{"int16::-32767", int16(-32767), []byte{sdtInt16, 128, 1}},
		{"int16::32767", int16(32767), []byte{sdtInt16, 127, 255}},
		{"int32::-2147483647", int32(-2147483647), []byte{sdtInt32, 128, 0, 0, 1}},
		{"int32::2147483647", int32(2147483647), []byte{sdtInt32, 127, 255, 255, 255}},
		{"int64::-9223372036854775807", int64(-9223372036854775807), []byte{sdtInt64, 128, 0, 0, 0, 0, 0, 0, 1}},
		{"int64::9223372036854775807", int64(9223372036854775807), []byte{sdtInt64, 127, 255, 255, 255, 255, 255, 255, 255}},
		// fails on 32bit arch
		{"int::-9223372036854775807", int(-9223372036854775807), []byte{sdtInt, 128, 0, 0, 0, 0, 0, 0, 1}},
		// fails on 32bit arch
		{"int::9223372036854775807", int(9223372036854775807), []byte{sdtInt, 127, 255, 255, 255, 255, 255, 255, 255}},
	}
}
