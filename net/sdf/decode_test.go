package sdf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sllt/sparrow/gen"
)

func TestDecodeBool(t *testing.T) {
	expect := true
	packet := []byte{sdtBool, 1}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
	expect = false
	packet = []byte{sdtBool, 0}
	value, _, err = Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceBool(t *testing.T) {
	expect := []bool{false, true, false}
	packet := []byte{sdtType, 0, 2,
		sdtSlice, sdtBool,
		sdtSlice,
		0, 0, 0, 3,
		0, 1, 0,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSliceAnyBool(t *testing.T) {
	expect := []any{false, true, false}
	packet := []byte{sdtType, 0, 2,
		sdtSlice, sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtBool, 0,
		sdtBool, 1,
		sdtBool, 0,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeAtom(t *testing.T) {
	expect := gen.Atom("hello world")
	packet := []byte{sdtAtom,
		0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeAtomCache(t *testing.T) {
	expect := gen.Atom("hello world")
	packet := []byte{sdtAtom,
		0x01, 0x2c, // cached "hello world" => 300
	}

	atomCache := new(sync.Map)
	atomCache.Store(uint16(300), expect)

	value, _, err := Decode(packet, Options{AtomCache: atomCache})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeAtomMapping(t *testing.T) {
	expect := gen.Atom("hello world")
	mapped := gen.Atom("hi")
	packet := []byte{sdtAtom,
		0, 0x02, // len
		0x68, 0x69, // "hi"
	}

	atomMapping := new(sync.Map)
	atomMapping.Store(mapped, expect)

	value, _, err := Decode(packet, Options{AtomMapping: atomMapping})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeAtomMappingCache(t *testing.T) {
	expect := gen.Atom("hello world")
	mapped := gen.Atom("hi")

	packet := []byte{sdtAtom,
		0x01, 0x2c, // mapped "hello world" => "hi", cached "hi" => 300
	}

	atomCache := new(sync.Map)
	atomCache.Store(uint16(300), mapped)
	atomMapping := new(sync.Map)
	atomMapping.Store(mapped, expect)

	value, _, err := Decode(packet, Options{AtomCache: atomCache, AtomMapping: atomMapping})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAtom(t *testing.T) {

	v := gen.Atom("hello world")
	expect := []gen.Atom{
		v, v, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAtom,
		sdtSlice,
		0, 0, 0, 3,
		0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
		0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
		0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSliceAtomCache(t *testing.T) {

	v := gen.Atom("hello world")
	expect := []gen.Atom{
		v, v, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAtom,
		sdtSlice,
		0, 0, 0, 3,
		0x01, 0x2c, // cached "hello world" => 300
		0x01, 0x2c, // cached "hello world" => 300
		0x01, 0x2c, // cached "hello world" => 300
	}
	atomCache := new(sync.Map)
	atomCache.Store(uint16(300), v)
	value, _, err := Decode(packet, Options{AtomCache: atomCache})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyAtom(t *testing.T) {

	v := gen.Atom("hello world")
	expect := []any{
		v, nil, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtAtom, 0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
		sdtNil,
		sdtAtom, 0, 0x0b, // len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // "hello world"
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyAtomCache(t *testing.T) {

	v := gen.Atom("hello world")
	expect := []any{
		v, nil, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtAtom, 0x01, 0x2c, // cached "hello world" => 300
		sdtNil,
		sdtAtom, 0x01, 0x2c, // cached "hello world" => 300
	}
	atomCache := new(sync.Map)
	atomCache.Store(uint16(300), v)
	value, _, err := Decode(packet, Options{AtomCache: atomCache})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeString(t *testing.T) {
	expect := "abc"
	packet := []byte{sdtString, 0, 3, 97, 98, 99}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceString(t *testing.T) {
	expect := []string{"abc", "def", "ghi"}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtString,
		sdtSlice,
		0, 0, 0, 3,
		0, 3, 97, 98, 99, // "abc"
		0, 3, 100, 101, 102, // "def"
		0, 3, 103, 104, 105, // "ghi"
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSliceAnyString(t *testing.T) {
	expect := []any{"abc", "def", "ghi"}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtString, 0, 3, 97, 98, 99, // "abc"
		sdtString, 0, 3, 100, 101, 102, // "def"
		sdtString, 0, 3, 103, 104, 105, // "ghi"
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeBinary(t *testing.T) {
	expect := []byte{1, 2, 3, 4, 5}
	packet := []byte{sdtBinary,
		0x0, 0x0, 0x0, 0x05, // len
		0x1, 0x2, 0x3, 0x4, 0x5,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceBinary(t *testing.T) {
	expect := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8}, {9}}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtBinary,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x0, 0x0, 0x05, // len
		0x1, 0x2, 0x3, 0x4, 0x5,
		0x0, 0x0, 0x0, 0x03, // len
		0x6, 0x7, 0x8,
		0x0, 0x0, 0x0, 0x01, // len
		0x9,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyBinary(t *testing.T) {
	expect := []any{[]byte{1, 2, 3, 4, 5}, []byte{6, 7, 8}, []byte{9}}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtBinary, 0x0, 0x0, 0x0, 0x05, // len
		0x1, 0x2, 0x3, 0x4, 0x5,
		sdtBinary, 0x0, 0x0, 0x0, 0x03, // len
		0x6, 0x7, 0x8,
		sdtBinary, 0x0, 0x0, 0x0, 0x01, // len
		0x9,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeFloat32(t *testing.T) {
	expect32 := float32(3.14)
	packet := []byte{sdtFloat32, 64, 72, 245, 195}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect32) {
		fmt.Println("exp", expect32)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSliceFloat32(t *testing.T) {

	expect := []float32{3.14, 3.15, 3.16}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtFloat32,
		sdtSlice,
		0, 0, 0, 3,
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyFloat32(t *testing.T) {

	expect := []any{float32(3.14), float32(3.15), float32(3.16)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat32, 0x40, 0x49, 0x99, 0x9a, // 3.15
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeFloat64(t *testing.T) {
	expect64 := float64(3.14)
	packet := []byte{sdtFloat64, 64, 9, 30, 184, 81, 235, 133, 31}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect64) {
		fmt.Println("exp", expect64)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceFloat64(t *testing.T) {

	expect := []float64{3.14, 3.15, 3.16}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtFloat64,
		sdtSlice,
		0, 0, 0, 3,
		0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		0x40, 0x9, 0x47, 0xae, 0x14, 0x7a, 0xe1, 0x48, // 3.16

	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyFloat64(t *testing.T) {

	expect := []any{float64(3.14), float64(3.15), float64(3.16)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtFloat64, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
		sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtFloat64, 0x40, 0x9, 0x47, 0xae, 0x14, 0x7a, 0xe1, 0x48, // 3.16
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeInteger(t *testing.T) {
	for _, c := range integerCases() {
		t.Run(c.name, func(t *testing.T) {
			value, _, err := Decode(c.bin, Options{})
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(value, c.integer) {
				fmt.Println("exp ", c.integer)
				fmt.Println("got ", value)
				t.Fatal("incorrect value")
			}
		})
	}
}

func TestDecodeSliceInt(t *testing.T) {
	expect := []int{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtInt,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyInt(t *testing.T) {
	expect := []any{int(1), int(2), int(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtInt, 0, 0, 0, 0, 0, 0, 0, 1,
		sdtInt, 0, 0, 0, 0, 0, 0, 0, 2,
		sdtInt, 0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceInt8(t *testing.T) {
	expect := []int8{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtInt8,
		sdtSlice,
		0, 0, 0, 3,
		1,
		2,
		3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyInt8(t *testing.T) {
	expect := []any{int8(1), int8(2), int8(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtInt8, 1,
		sdtInt8, 2,
		sdtInt8, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceInt16(t *testing.T) {
	expect := []int16{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtInt16,
		sdtSlice,
		0, 0, 0, 3,
		0, 1,
		0, 2,
		0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyInt16(t *testing.T) {
	expect := []any{int16(1), int16(2), int16(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtInt16, 0, 1,
		sdtInt16, 0, 2,
		sdtInt16, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceInt32(t *testing.T) {
	expect := []int32{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtInt32,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 1,
		0, 0, 0, 2,
		0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyInt32(t *testing.T) {
	expect := []any{int32(1), int32(2), int32(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtInt32, 0, 0, 0, 1,
		sdtInt32, 0, 0, 0, 2,
		sdtInt32, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceInt64(t *testing.T) {
	expect := []int64{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtInt64,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyInt64(t *testing.T) {
	expect := []any{int64(1), int64(2), int64(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtInt64, 0, 0, 0, 0, 0, 0, 0, 1,
		sdtInt64, 0, 0, 0, 0, 0, 0, 0, 2,
		sdtInt64, 0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUint(t *testing.T) {
	expect := []uint{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtUint,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyUint(t *testing.T) {
	expect := []any{uint(1), uint(2), uint(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtUint, 0, 0, 0, 0, 0, 0, 0, 1,
		sdtUint, 0, 0, 0, 0, 0, 0, 0, 2,
		sdtUint, 0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUint8(t *testing.T) {
	// basically, []uint8 == []byte, which means it should be encoded as a binary,
	// but check this way of encoding anyway.
	expect := []uint8{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtUint8,
		sdtSlice,
		0, 0, 0, 3,
		1,
		2,
		3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyUint8(t *testing.T) {
	expect := []any{uint8(1), byte(2), uint8(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtUint8, 1,
		sdtUint8, 2,
		sdtUint8, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUint16(t *testing.T) {
	expect := []uint16{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtUint16,
		sdtSlice,
		0, 0, 0, 3,
		0, 1,
		0, 2,
		0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyUint16(t *testing.T) {
	expect := []any{uint16(1), uint16(2), uint16(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtUint16, 0, 1,
		sdtUint16, 0, 2,
		sdtUint16, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUint32(t *testing.T) {
	expect := []uint32{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtUint32,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 1,
		0, 0, 0, 2,
		0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyUint32(t *testing.T) {
	expect := []any{uint32(1), uint32(2), uint32(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtUint32, 0, 0, 0, 1,
		sdtUint32, 0, 0, 0, 2,
		sdtUint32, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUint64(t *testing.T) {
	expect := []uint64{1, 2, 3}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtUint64,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyUint64(t *testing.T) {
	expect := []any{uint64(1), uint64(2), uint64(3)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtUint64, 0, 0, 0, 0, 0, 0, 0, 1,
		sdtUint64, 0, 0, 0, 0, 0, 0, 0, 2,
		sdtUint64, 0, 0, 0, 0, 0, 0, 0, 3,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSliceAnyInteger(t *testing.T) {
	expect := []any{
		int(1), nil, int8(2), nil, int16(3), nil, int32(4), nil, int64(5), nil,
		uint(6), nil, uint8(7), nil, uint16(8), nil, uint32(9), nil, uint64(10),
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 19,
		sdtInt, 0, 0, 0, 0, 0, 0, 0, 1,
		sdtNil,
		sdtInt8, 2,
		sdtNil,
		sdtInt16, 0, 3,
		sdtNil,
		sdtInt32, 0, 0, 0, 4,
		sdtNil,
		sdtInt64, 0, 0, 0, 0, 0, 0, 0, 5,
		sdtNil,
		sdtUint, 0, 0, 0, 0, 0, 0, 0, 6,
		sdtNil,
		sdtUint8, 7,
		sdtNil,
		sdtUint16, 0, 8,
		sdtNil,
		sdtUint32, 0, 0, 0, 9,
		sdtNil,
		sdtUint64, 0, 0, 0, 0, 0, 0, 0, 10,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnySlice(t *testing.T) {
	expect := []any{
		[]int{4},
		nil,
		[]float32{3.14, 3.15, 3.16},
		nil,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 4,

		sdtType, 0, 2,
		sdtSlice, sdtInt,
		sdtSlice, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 4,

		sdtNil,

		sdtType, 0, 2,
		sdtSlice, sdtFloat32,
		sdtSlice, 0, 0, 0, 3,
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16

		sdtNil,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeTime(t *testing.T) {
	expect := time.Date(1399, time.January, 26, 0, 0, 0, 0, time.UTC)
	packet := []byte{sdtTime,
		0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceTime(t *testing.T) {
	v := time.Date(1399, time.January, 26, 0, 0, 0, 0, time.UTC)
	expect := []time.Time{
		v, v, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtTime,
		sdtSlice,
		0, 0, 0, 3,
		0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
		0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
		0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyTime(t *testing.T) {
	v := time.Date(1399, time.January, 26, 0, 0, 0, 0, time.UTC)
	expect := []any{
		v, v, v,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtTime, 0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
		sdtTime, 0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
		sdtTime, 0xf, // len
		0x1, 0x0, 0x0, 0x0, 0xa, 0x45, 0xaf, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceReg(t *testing.T) {
	type MyFloat float32
	var x MyFloat

	expect := []MyFloat{3.14, 3.15, 3.16}
	packet := []byte{sdtType, 0, 39,
		sdtSlice,
		sdtReg, 0, 35,
		// name: #github.com/sllt/sparrow/net/sdf/MyFloat
		0x23, 0x65, 0x72, 0x67, 0x6f, 0x2e, 0x73, 0x65,
		0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65,
		0x72, 0x67, 0x6f, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
		0x65, 0x64, 0x66, 0x2f, 0x4d, 0x79, 0x46, 0x6c,
		0x6f, 0x61, 0x74,
		sdtSlice,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}
	RegisterTypeOf(x)

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceRegCache(t *testing.T) {
	type MyFloat123 float32
	var x MyFloat123

	expect := []MyFloat123{3.14, 3.15, 3.16}
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtReg, 0x13, 0x88, // cache id uint16(5000) => name: #github.com/sllt/sparrow/net/sdf/MyFloat123
		sdtSlice,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}

	RegisterTypeOf(x)

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MyFloat123")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegSlice(t *testing.T) {
	type MySlice29 []float32

	expect := MySlice29{3.14, 3.15, 3.16}
	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}
	if err := RegisterTypeOf(expect); err != nil {
		t.Fatal(err)
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MySlice29")

	opts := Options{
		RegCache: regCache,
	}

	value, _, err := Decode(packet, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegSliceRegSlice(t *testing.T) {
	type MyDecSliceFloat []float32
	type MyDecSliceOfSlice []MyDecSliceFloat

	expect := MyDecSliceOfSlice{
		{3.14, 3.15, 3.16},
		nil,
		{3.14},
	}
	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0x0, 0x0, 0x0, 0x3,
		sdtSlice,
		0x0, 0x0, 0x0, 0x3,
		0x40, 0x48, 0xf5, 0xc3,
		0x40, 0x49, 0x99, 0x9a,
		0x40, 0x4a, 0x3d, 0x71,
		sdtNil,
		sdtSlice,
		0x0, 0x0, 0x0, 0x1,
		0x40, 0x48, 0xf5, 0xc3,
	}

	if err := RegisterTypeOf(MyDecSliceOfSlice{}); err != nil {
		t.Fatal(err)
	}
	if err := RegisterTypeOf(MyDecSliceFloat{}); err != nil {
		t.Fatal(err)
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MyDecSliceOfSlice")

	opts := Options{
		RegCache: regCache,
	}

	value, _, err := Decode(packet, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyReg(t *testing.T) {
	type MyFloat20 float32
	x := MyFloat20(3.16)

	expect := []any{float32(3.14), float64(3.15), float32(3.16), nil, x}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 5,
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtNil,
		sdtReg, 0x13, 0x88,
		0x40, 0x4a, 0x3d, 0x71, // MyFloat20(3.16)
	}

	RegisterTypeOf(x)

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MyFloat20")

	opts := Options{
		RegCache: regCache,
	}

	value, _, err := Decode(packet, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %t\n", expect)
		fmt.Printf("got %t\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegSliceReg(t *testing.T) {
	type MyFloat19 float32
	type MySlice19 []MyFloat19
	var x MyFloat19

	expect := MySlice19{3.14, 3.15, 3.16}
	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}

	if err := RegisterTypeOf(x); err != nil {
		t.Fatal(err)
	}
	if err := RegisterTypeOf(expect); err != nil {
		t.Fatal(err)
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MySlice19")

	opts := Options{
		RegCache: regCache,
	}

	value, _, err := Decode(packet, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodePID(t *testing.T) {
	expect := gen.PID{Node: "abc@def", ID: 32767, Creation: 2}
	packet := []byte{sdtPID,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSlicePID(t *testing.T) {
	v := gen.PID{Node: "abc@def", ID: 32767, Creation: 2}
	expect := []gen.PID{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtPID,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyPID(t *testing.T) {
	v := gen.PID{Node: "abc@def", ID: 32767, Creation: 2}
	expect := []any{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtPID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		sdtPID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		sdtPID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7f, 0xff, // id
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeProcessID(t *testing.T) {
	expect := gen.ProcessID{Node: "abc@def", Name: "ghi"}
	packet := []byte{sdtProcessID,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceProcessID(t *testing.T) {
	v := gen.ProcessID{Node: "abc@def", Name: "ghi"}
	expect := []gen.ProcessID{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtProcessID,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyProcessID(t *testing.T) {
	v := gen.ProcessID{Node: "abc@def", Name: "ghi"}
	expect := []any{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtProcessID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		sdtProcessID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
		sdtProcessID, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeEvent(t *testing.T) {
	expect := gen.Event{Node: "abc@def", Name: "ghi"}
	packet := []byte{sdtEvent,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceEvent(t *testing.T) {
	v := gen.Event{Node: "abc@def", Name: "ghi"}
	expect := []gen.Event{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtEvent,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyEvent(t *testing.T) {
	v := gen.Event{Node: "abc@def", Name: "ghi"}
	expect := []any{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtEvent, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (process name)
		0x67, 0x68, 0x69,
		sdtEvent, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
		sdtEvent, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x3, // len atom (node name)
		0x67, 0x68, 0x69,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRef(t *testing.T) {
	expect := gen.Ref{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	packet := []byte{sdtRef,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceRef(t *testing.T) {
	v := gen.Ref{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	expect := []gen.Ref{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtRef,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyRef(t *testing.T) {
	v := gen.Ref{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	expect := []any{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtRef, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		sdtRef, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		sdtRef, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeAlias(t *testing.T) {
	expect := gen.Alias{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	packet := []byte{sdtAlias,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAlias(t *testing.T) {
	v := gen.Alias{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	expect := []gen.Alias{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAlias,
		sdtSlice,
		0, 0, 0, 3,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyAlias(t *testing.T) {
	v := gen.Alias{Node: "abc@def", ID: [3]uint64{4, 5, 6}, Creation: 2}
	expect := []any{v, v, v}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtAlias, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		sdtAlias, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		sdtAlias, 0x0, 0x7, // len atom (node name)
		0x61, 0x62, 0x63, 0x40, 0x64, 0x65, 0x66,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, // creation
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeError(t *testing.T) {
	packet := []byte{sdtError,
		0, 3, // len
		97, 98, 99, // "abc"
	}
	expect := errors.New("abc")

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceError(t *testing.T) {
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtError,
		sdtSlice,
		0, 0, 0, 3,
		0, 4, // len
		97, 98, 99, 100, // "abcd"
		0, 4, // len
		97, 98, 99, 100, // "abcd"
		0, 4, // len
		97, 98, 99, 100, // "abcd"
	}
	v := errors.New("abcd")
	expect := []error{v, v, v}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceErrorNil(t *testing.T) {
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtError,
		sdtSlice,
		0, 0, 0, 3,
		0, 4, // len
		97, 98, 99, 100, // "abcd"
		0xff, 0xff, // nil error
		0, 4, // len
		97, 98, 99, 100, // "abcd"
	}
	v := errors.New("abcd")
	expect := []error{v, nil, v}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegError(t *testing.T) {
	packet := []byte{sdtError,
		0x88, 0xb8, // 35000 => error "abc"
	}

	expect := errors.New("abc")

	errCache := new(sync.Map)
	errCache.Store(uint16(35000), expect)

	value, _, err := Decode(packet, Options{ErrCache: errCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}

}

func TestDecodeSliceAnyError(t *testing.T) {
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtError, 0, 4, // len
		97, 98, 99, 100, // "abcd"
		sdtNil,
		sdtError, 0, 4, // len
		97, 98, 99, 100, // "abcd"
	}
	v := errors.New("abcd")
	expect := []any{v, nil, v}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceTypeReg(t *testing.T) {
	type MyFloa float32
	var x MyFloa
	expect := []MyFloa{3.14, 3.15, 3.16}
	packet := []byte{sdtType, 0, 38,
		sdtSlice,
		sdtReg, 0, 34,
		// name: #github.com/sllt/sparrow/net/proto/sdf/MyFloa
		0x23, 0x65, 0x72, 0x67, 0x6f, 0x2e, 0x73, 0x65,
		0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65,
		0x72, 0x67, 0x6f, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
		0x65, 0x64, 0x66, 0x2f, 0x4d, 0x79, 0x46, 0x6c,
		0x6f, 0x61,
		sdtSlice,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}
	RegisterTypeOf(x)

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}

}

func TestDecodeSliceRegTypeReg(t *testing.T) {
	type MyFloatD19 float32
	type MySliceD19 []MyFloatD19
	var x MyFloatD19

	expect := MySliceD19{3.14, 3.15, 3.16}
	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0, 0, 0, 3, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}
	RegisterTypeOf(x)
	if err := RegisterTypeOf(expect); err != nil {
		t.Fatal(err)
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MySliceD19")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}

}

func TestDecodeSliceAny(t *testing.T) {
	expect := []any{float32(3.14), float64(3.15), float32(3.16)}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}

}

func TestDecodeSliceNil(t *testing.T) {

	expect := []any{nil, nil, nil}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 3,
		sdtNil,
		sdtNil,
		sdtNil,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceNil2(t *testing.T) {

	expect := []any{}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 0,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceNest(t *testing.T) {
	expect := []any{
		[]any{float32(3.15)},
		float32(3.14),
		float32(3.16),
		[]any{float64(3.15)},
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 4,
		sdtType, 0, 2,
		sdtSlice, sdtAny,
		sdtSlice,
		0, 0, 0, 1, sdtFloat32, 0x40, 0x49, 0x99, 0x9a, // 3.15
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtType, 0, 2,
		sdtSlice, sdtAny,
		sdtSlice,
		0, 0, 0, 1, sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceSlice(t *testing.T) {
	expect := [][]float32{
		{3.14, 3.15, 3.16},
		{3.16},
		nil,
		{3.14, 3.15},
		{},
	}
	packet := []byte{sdtType, 0, 3,
		sdtSlice,
		sdtSlice,
		sdtFloat32,
		sdtSlice,
		0, 0, 0, 5,
		sdtSlice,
		0, 0, 0, 3, // first slice with 3 items
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtSlice,
		0, 0, 0, 1, // second slice with 1 item
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtNil, // third one
		sdtSlice,
		0, 0, 0, 2, // 4th
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x49, 0x99, 0x9a, // 3.15
		sdtSlice,
		0, 0, 0, 0, // 5th
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceSliceAny(t *testing.T) {
	expect := [][]any{
		{float32(3.14), float32(3.16), float64(3.15)},
		{float64(3.15)},
		nil,
		{float32(3.14), float32(3.16)},
		{},
	}
	packet := []byte{sdtType, 0, 3,
		sdtSlice,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 5,
		sdtSlice,
		0, 0, 0, 3, // first slice with 3 items
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtSlice,
		0, 0, 0, 1, // second slice with 1 item
		sdtFloat64, 0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtNil, // third one
		sdtSlice,
		0, 0, 0, 2, // 4th
		sdtFloat32, 0x40, 0x48, 0xf5, 0xc3, // 3.14
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtSlice,
		0, 0, 0, 0, // 5th
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceSliceNil(t *testing.T) {
	expect := [][]any{nil, []any{}, nil, nil}
	packet := []byte{sdtType, 0, 3,
		sdtSlice,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 4,
		sdtNil,
		sdtSlice,
		0, 0, 0, 0,
		sdtNil,
		sdtNil,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceSliceReg(t *testing.T) {
	type MySlice555 []float32

	if err := RegisterTypeOf(MySlice555{}); err != nil {
		if err != gen.ErrTaken {
			t.Fatal(err)
		}
	}

	expect := []MySlice555{
		MySlice555{3.14, 3.16, 3.15},
		MySlice555{3.15},
		nil,
		MySlice555{3.14, 3.16},
		MySlice555{},
	}
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtReg, 0x13, 0x88,
		sdtSlice,
		0, 0, 0, 5,
		sdtReg,
		0, 0, 0, 3, // first slice with 3 items
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x49, 0x99, 0x9a, // 3.15
		sdtReg,
		0, 0, 0, 1, // second slice with 1 item
		0x40, 0x49, 0x99, 0x9a, // 3.15
		sdtNil,
		sdtReg,
		0, 0, 0, 2, // 4th
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtReg,
		0, 0, 0, 0, // third one
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MySlice555")
	opts := Options{
		RegCache: regCache,
	}

	value, _, err := Decode(packet, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSlice3DZero(t *testing.T) {
	expect := [][][]float32{}
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtSlice,
		sdtSlice,
		sdtFloat32,
		sdtSlice,
		0, 0, 0, 0,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeSlice3D(t *testing.T) {
	expect := [][][]float32{ /* len 3 */
		{ /* len 5 */
			{ /* len 7 */ 2.21018848, 2.94523878, 1.67807658, 1.30014748, 1.1873558, 8.1819557, 3.2368748},
			{ /* len 10 */ 2.17948558, 2.95483828, 3.29734688, 2.72996818, 2.50011478, 2.98767788, 1.31364818, 8.06395757, 2.53354848, 2.38570578},
			{ /* len 4 */ 2.9838078, 1.61728128, 1.8756628, 1.5756598},
			{ /* len 10 */ 8.5187367, 2.79348, 4.3456557, 1.29794587, 3.38391948, 1.4460748, 5.0206397, 2.02001097, 1.77825548, 2.33810328},
			{ /* len 8 */ 3.15617888, 2.21068618, 3.01507718, 7.0342597, 2.12085158, 7.9914467, 2.92003388, 3.19992137},
		}, { /* len 6 */
			{ /* len 3 */ 3.3188187, 2.82300078, 7.3257346},
			{ /* len 10 */ 1.47951058, 1.47638718, 3.1678068, 1.24334058, 1.48100658, 1.8274938, 2.07265258, 1.83188888, 5.8776197, 1.64099568},
			{ /* len 6 */ 2.26154558, 9.5987497, 3.24544727, 1.34864688, 2.47839448, 2.0456888},
			{ /* len 5 */ 9.0369537, 3.69528477, 3.04563028, 1.4488858, 3.80179227},
			{ /* len 5 */ 1.53326348, 2.77105168, 1.05977548, 2.75297638, 8.9171847},
			{ /* len 10 */ 1.65367358, 9.4070457, 3.06440548, 2.4763148, 2.22120158, 2.3734938, 3.37481478, 2.22900497, 6.2138987, 2.80613798},
		}, { /* len 1 */
			{ /* len 10 */ 8.03434337, 2.55059418, 2.20168828, 2.86517478, 4.38993137, 8.6655217, 2.22159657, 3.0119788, 1.19758818, 2.58799087},
		},
	}

	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtSlice,
		sdtSlice,
		sdtFloat32,

		sdtSlice,
		0x0, 0x0, 0x0, 0x3, // len 3 { x, x, x}
		sdtSlice,
		0x0, 0x0, 0x0, 0x5, // len 5 { {y, y, y, y, y}, x, x}
		sdtSlice,
		0x0, 0x0, 0x0, 0x7, // len 7 { { {z, z, z, z, z, z, z}, y, y, y, y}, x, x}
		0x40, 0xd, 0x73, 0xba, // z
		0x40, 0x3c, 0x7e, 0xcb, // z
		0x3f, 0xd6, 0xcb, 0x37, // z
		0x3f, 0xa6, 0x6b, 0x3c, // z
		0x3f, 0x97, 0xfb, 0x46, // z
		0x41, 0x2, 0xe9, 0x4a, // z
		0x40, 0x4f, 0x28, 0xf5, // z
		sdtSlice,
		0x0, 0x0, 0x0, 0xa, // len 10
		0x40, 0xb, 0x7c, 0xb1,
		0x40, 0x3d, 0x1c, 0x12,
		0x40, 0x53, 0x7, 0xbb,
		0x40, 0x2e, 0xb7, 0xcc,
		0x40, 0x20, 0x1, 0xe1,
		0x40, 0x3f, 0x36, 0x1d,
		0x3f, 0xa8, 0x25, 0xa0,
		0x41, 0x1, 0x5, 0xf8,
		0x40, 0x22, 0x25, 0xa9,
		0x40, 0x18, 0xaf, 0x67,
		sdtSlice,
		0x0, 0x0, 0x0, 0x4, // len 4
		0x40, 0x3e, 0xf6, 0xb5,
		0x3f, 0xcf, 0x3, 0x13,
		0x3f, 0xf0, 0x15, 0xb8,
		0x3f, 0xc9, 0xaf, 0x38,
		sdtSlice,
		0x0, 0x0, 0x0, 0xa, // len 10
		0x41, 0x8, 0x4c, 0xbf,
		0x40, 0x32, 0xc8, 0x60,
		0x40, 0x8b, 0xf, 0x9d,
		0x3f, 0xa6, 0x23, 0x17,
		0x40, 0x58, 0x92, 0x23,
		0x3f, 0xb9, 0x18, 0xfb,
		0x40, 0xa0, 0xa9, 0x15,
		0x40, 0x1, 0x47, 0xdc,
		0x3f, 0xe3, 0x9d, 0xe0,
		0x40, 0x15, 0xa3, 0x7c,
		sdtSlice,
		0x0, 0x0, 0x0, 0x8, // len 8
		0x40, 0x49, 0xfe, 0xd6,
		0x40, 0xd, 0x7b, 0xe2,
		0x40, 0x40, 0xf7, 0x6,
		0x40, 0xe1, 0x18, 0xa8,
		0x40, 0x7, 0xbc, 0x8,
		0x40, 0xff, 0xb9, 0xee,
		0x40, 0x3a, 0xe1, 0xd6,
		0x40, 0x4c, 0xcb, 0x83,
		sdtSlice,
		0x0, 0x0, 0x0, 0x6, // len 6
		sdtSlice,
		0x0, 0x0, 0x0, 0x3, // len 3
		0x40, 0x54, 0x67, 0x87,
		0x40, 0x34, 0xac, 0xb,
		0x40, 0xea, 0x6c, 0x6b,
		sdtSlice,
		0x0, 0x0, 0x0, 0xa, // len 10
		0x3f, 0xbd, 0x60, 0x9a,
		0x3f, 0xbc, 0xfa, 0x41,
		0x40, 0x4a, 0xbd, 0x59,
		0x3f, 0x9f, 0x25, 0xc9,
		0x3f, 0xbd, 0x91, 0xa0,
		0x3f, 0xe9, 0xeb, 0x51,
		0x40, 0x4, 0xa6, 0x57,
		0x3f, 0xea, 0x7b, 0x56,
		0x40, 0xbc, 0x15, 0x76,
		0x3f, 0xd2, 0xc, 0x25,
		sdtSlice,
		0x0, 0x0, 0x0, 0x6, // len 6
		0x40, 0x10, 0xbd, 0x2a,
		0x41, 0x19, 0x94, 0x7b,
		0x40, 0x4f, 0xb5, 0x68,
		0x3f, 0xac, 0xa0, 0x76,
		0x40, 0x1e, 0x9e, 0x4,
		0x40, 0x2, 0xec, 0x91,
		sdtSlice,
		0x0, 0x0, 0x0, 0x5, // len 5
		0x41, 0x10, 0x97, 0x5d,
		0x40, 0x6c, 0x7f, 0x8c,
		0x40, 0x42, 0xeb, 0x9b,
		0x3f, 0xb9, 0x75, 0x17,
		0x40, 0x73, 0x50, 0x91,
		sdtSlice,
		0x0, 0x0, 0x0, 0x5, // len 5
		0x3f, 0xc4, 0x41, 0xfa,
		0x40, 0x31, 0x58, 0xe9,
		0x3f, 0x87, 0xa6, 0xb9,
		0x40, 0x30, 0x30, 0xc4,
		0x41, 0xe, 0xac, 0xca,
		sdtSlice,
		0x0, 0x0, 0x0, 0xa, // len 10
		0x3f, 0xd3, 0xab, 0x93,
		0x41, 0x16, 0x83, 0x42,
		0x40, 0x44, 0x1f, 0x38,
		0x40, 0x1e, 0x7b, 0xf1,
		0x40, 0xe, 0x28, 0x2b,
		0x40, 0x17, 0xe7, 0x53,
		0x40, 0x57, 0xfc, 0xf7,
		0x40, 0xe, 0xa8, 0x4,
		0x40, 0xc6, 0xd8, 0x42,
		0x40, 0x33, 0x97, 0xc4,
		sdtSlice,
		0x0, 0x0, 0x0, 0x1, // len 1
		sdtSlice,
		0x0, 0x0, 0x0, 0xa, // len 10
		0x41, 0x0, 0x8c, 0xac,
		0x40, 0x23, 0x3c, 0xef,
		0x40, 0xc, 0xe8, 0x76,
		0x40, 0x37, 0x5f, 0x6,
		0x40, 0x8c, 0x7a, 0x51,
		0x41, 0xa, 0xa5, 0xfa,
		0x40, 0xe, 0x2e, 0xa3,
		0x40, 0x40, 0xc4, 0x43,
		0x3f, 0x99, 0x4a, 0x92,
		0x40, 0x25, 0xa1, 0xa4,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type testUnmarshal struct {
	X float32
}

func (testUnmarshal) MarshalSDF(io.Writer) error {
	return nil
}

func (u *testUnmarshal) UnmarshalSDF(data []byte) error {
	u.X = math.Float32frombits(binary.BigEndian.Uint32(data))
	return nil
}

func TestDecodeUnmarshal(t *testing.T) {
	packet := []byte{sdtReg, 0x13, 0x88,
		0, 0, 0, 4, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
	}
	expect := testUnmarshal{X: 3.14}

	if err := RegisterTypeOf(testUnmarshal{}); err != nil {
		t.Fatal(err)
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testUnmarshal")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceUnmarshal(t *testing.T) {
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtReg, 0x13, 0x88,
		sdtSlice,
		0, 0, 0, 3,
		0, 0, 0, 4, // len
		0x40, 0x48, 0xf5, 0xc3, // 3.14
		0, 0, 0, 4, // len
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0, 0, 0, 4, // len
		0x40, 0x4a, 0x3d, 0x71, // 3.16
	}
	expect := []testUnmarshal{{3.14}, {3.15}, {3.16}}

	RegisterTypeOf(testUnmarshal{})

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testUnmarshal")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Println("exp", expect)
		fmt.Println("got", value)
		t.Fatal("incorrect value")
	}
}

type testDecStruct struct {
	A float32
	B float64
}

func TestDecodeStruct(t *testing.T) {
	packet := []byte{sdtReg, 0x13, 0x88,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
	}
	expect := testDecStruct{3.16, 3.15}

	RegisterTypeOf(testDecStruct{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStruct")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceStruct(t *testing.T) {
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtReg, 0x13, 0x88,
		sdtSlice,
		0, 0, 0, 2,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
	}
	expect := []testDecStruct{{3.16, 3.15}, {3.15, 3.14}}

	RegisterTypeOf(testDecStruct{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStruct")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type testDecStructWithAny struct {
	A float32
	B float64
	C any
}

func TestDecodeStructWithAny(t *testing.T) {
	packet := []byte{sdtReg, 0x13, 0x88,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtNil,
	}
	expect := testDecStructWithAny{3.16, 3.15, nil}

	RegisterTypeOf(testDecStructWithAny{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStructWithAny")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}

	packet = []byte{sdtReg, 0x13, 0x88,
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
		sdtFloat64, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
	}
	expect = testDecStructWithAny{3.15, 3.14, float64(3.14)}

	value, _, err = Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type testDecSliceString []string
type testDecStructWithSlice struct {
	A float32
	B float64
	C []bool
	D testDecSliceString
	E []int
}

func TestDecodeStructWithSlice(t *testing.T) {
	packet := []byte{sdtReg, 0x13, 0x88,
		0x40, 0x4a, 0x3d, 0x71, // 3.16 (float32)
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15 (float64)
		sdtSlice,
		0x0, 0x0, 0x0, 0x2, // len of []bool
		0x1, 0x0, // true, false
		sdtReg,             // testDecSliceString
		0x0, 0x0, 0x0, 0x2, // len of testDecSliceString
		0x0, 0x4, // len of "true"
		0x74, 0x72, 0x75, 0x65, // "true"
		0x0, 0x5, // len of "false"
		0x66, 0x61, 0x6c, 0x73, 0x65, // "false"
		sdtNil, // nil value of []int
	}
	expect := testDecStructWithSlice{
		3.16,
		3.15,
		[]bool{true, false},
		testDecSliceString{"true", "false"},
		nil,
	}

	RegisterTypeOf(testDecSliceString{})
	RegisterTypeOf(testDecStructWithSlice{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStructWithSlice")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceStructWithAny(t *testing.T) {
	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtReg, 0x13, 0x88,
		sdtSlice,
		0, 0, 0, 3,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtNil,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
		sdtFloat64, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
	}
	expect := []testDecStructWithAny{
		{3.16, 3.15, nil},
		{3.16, 3.15, float32(3.16)},
		{3.15, 3.14, float64(3.14)},
	}
	RegisterTypeOf(testDecStructWithAny{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStructWithAny")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyWithStruct(t *testing.T) {
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,
		sdtSlice,
		0, 0, 0, 4,
		sdtNil,
		sdtReg, 0x13, 0x88,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		0x40, 0x9, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 3.15
		sdtNil,
		sdtReg, 0x13, 0x88,
		0x40, 0x49, 0x99, 0x9a, // 3.15
		0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, // 3.14
	}
	expect := []any{
		nil,
		testDecStruct{3.16, 3.15},
		nil,
		testDecStruct{3.15, 3.14},
	}
	RegisterTypeOf(testDecStruct{})
	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testDecStruct")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMap(t *testing.T) {
	expect := map[int16]string{
		8: "hello",
		9: "world",
	}

	packet := []byte{sdtType, 0, 3,
		sdtMap,
		sdtInt16,
		sdtString,
		sdtMap,
		0, 0, 0, 2,
		0, 8, // key 8
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 9, // key 9
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapAnyString(t *testing.T) {
	expect := map[any]string{
		nil:      "hello",
		int16(9): "world",
	}

	packet := []byte{sdtType, 0, 3,
		sdtMap,
		sdtAny,
		sdtString,
		sdtMap,
		0, 0, 0, 2,
		sdtNil,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtInt16, 0, 9, // key 9
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapStringAny(t *testing.T) {
	expect := map[string]any{
		"hello": nil,
		"world": int16(9),
		"helloo": map[float32]any{
			3.16: uint16(3),
		},
		"worldd": map[float64]any{},
	}

	packet := []byte{sdtType, 0, 3,
		sdtMap,
		sdtString,
		sdtAny,
		sdtMap, 0, 0, 0, 4,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtNil,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		sdtInt16, 0, 9, // key 9
		0, 6, // len of value "helloo"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x6f, // "helloo"
		sdtType, 0, 3,
		sdtMap, sdtFloat32, sdtAny,
		sdtMap, 0, 0, 0, 1,
		0x40, 0x4a, 0x3d, 0x71, // 3.16
		sdtUint16, 0, 3,
		0, 6, // len of value "worldd"
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x64, // "worldd"
		sdtType, 0, 3,
		sdtMap, sdtFloat64, sdtAny,
		sdtMap, 0, 0, 0, 0,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapStringMapAny(t *testing.T) {
	expect := map[string]map[any]int16{
		"hello": nil,
		"world": {
			float32(3.16): int16(9),
		},
	}

	packet := []byte{sdtType, 0, 5,
		sdtMap,
		sdtString,
		sdtMap,
		sdtAny,
		sdtInt16,
		sdtMap, 0, 0, 0, 2,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtNil,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		sdtMap, 0, 0, 0, 1,
		sdtFloat32, 0x40, 0x4a, 0x3d, 0x71, // 3.16
		0, 9, // 9
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeMapStringMapNilZero(t *testing.T) {
	expect := map[string]map[any]int16{
		"hello": nil,
		"world": {},
	}

	packet := []byte{sdtType, 0, 5,
		sdtMap,
		sdtString,
		sdtMap,
		sdtAny,
		sdtInt16,
		sdtMap, 0, 0, 0, 2,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtNil,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		sdtMap, 0, 0, 0, 0,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMap3DZero(t *testing.T) {
	expect := map[int16]map[string]map[float32]int{}
	packet := []byte{sdtType, 0, 7,
		sdtMap,
		sdtInt16,
		sdtMap,
		sdtString,
		sdtMap,
		sdtFloat32,
		sdtInt,
		sdtMap, 0, 0, 0, 0,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapZero(t *testing.T) {
	expect := map[int16]string{}

	packet := []byte{sdtType, 0, 3,
		sdtMap,
		sdtInt16, sdtString,
		sdtMap, 0, 0, 0, 0,
	}
	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceMap(t *testing.T) {
	expect := []map[int16]string{
		{
			8: "hello",
		}, {

			10: "helloo",
		},
		{
			12: "hellooo",
		},
	}

	packet := []byte{sdtType, 0, 4,
		sdtSlice,
		sdtMap,
		sdtInt16,
		sdtString,
		sdtSlice,
		0, 0, 0, 3,
		sdtMap,
		0, 0, 0, 1,
		0, 8, // key 8
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtMap,
		0, 0, 0, 1, // len of second map
		0, 0xa, // key 10
		0, 6, // len of value "helloo"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x6f, // "helloo"
		sdtMap,
		0, 0, 0, 1, // len of 3rd map
		0, 0xc, // key 12
		0, 7, // len of value "helloo"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x6f, 0x6f, // "hellooo"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}
func TestDecodeMapValueSliceNil(t *testing.T) {
	expect := map[int16][]any{
		int16(8): nil,
		int16(9): []any{"world"},
	}

	packet := []byte{sdtType, 0, 4,
		sdtMap,
		sdtInt16,
		sdtSlice,
		sdtAny,
		sdtMap,
		0, 0, 0, 2,

		0, 8, // key 8
		sdtNil,

		0, 9, // key 9
		sdtSlice,
		0, 0, 0, 1,
		sdtString, 0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapValueMap(t *testing.T) {
	expect := map[int16]map[string]int{
		int16(8): nil,
		int16(9): {
			"world": 10,
		},
	}

	packet := []byte{sdtType, 0, 5,
		sdtMap,
		sdtInt16,
		sdtMap,
		sdtString,
		sdtInt,
		sdtMap,
		0, 0, 0, 2,

		0, 8, // key
		sdtNil,

		0, 9, // 9 => map
		sdtMap,
		0, 0, 0, 1,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, 0, 0, 0, 0, 0, 0, 0xa, // key 10
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type testMapKey string

func TestDecodeMapValueMapRegKey(t *testing.T) {
	var x testMapKey

	if err := RegisterTypeOf(x); err != nil {
		if err != gen.ErrTaken {
			t.Fatal(err)
		}
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testMapKey")

	expect := map[int16]map[testMapKey]int{
		int16(8): nil,
		int16(9): {
			"world": 10,
		},
	}

	packet := []byte{sdtType, 0, 7,
		sdtMap,
		sdtInt16,
		sdtMap,
		sdtReg, 0x13, 0x88,
		sdtInt,

		sdtMap,
		0, 0, 0, 2,

		0, 8, // 8 => map
		sdtNil,

		0, 9, // 9 => map
		sdtMap,
		0, 0, 0, 1,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, 0, 0, 0, 0, 0, 0, 0xa, // value 10
	}

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeMapValueMapAnyWithRegKey(t *testing.T) {
	var x testMapKey = "world"

	if err := RegisterTypeOf(x); err != nil {
		if err != gen.ErrTaken {
			t.Fatal(err)
		}
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/testMapKey")

	expect := map[int16]map[any]int{
		int16(8): nil,
		int16(9): {
			x: 10,
		},
	}

	packet := []byte{sdtType, 0, 5,
		sdtMap,
		sdtInt16,
		sdtMap,
		sdtAny,
		sdtInt,

		sdtMap,
		0, 0, 0, 2,

		0, 8, // 8 => map
		sdtNil,

		0, 9, // 9 => map
		sdtMap,
		0, 0, 0, 1,
		sdtReg, 0x13, 0x88, 0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, 0, 0, 0, 0, 0, 0, 0xa, // value 10
	}

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type MyMap map[string]bool

func TestDecodeRegMap(t *testing.T) {
	var x MyMap

	if err := RegisterTypeOf(x); err != nil {
		if err != gen.ErrTaken {
			t.Fatal(err)
		}
	}

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/MyMap")

	expect := MyMap{
		"hello": true,
		"world": false,
	}
	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0, 0, 0, 2,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		1,    // true
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, // false
	}

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegMapRegSlice(t *testing.T) {
	type myDecSlice90 []bool
	type myDecMap90 map[string]myDecSlice90

	RegisterTypeOf(myDecSlice90{})
	RegisterTypeOf(myDecMap90{})

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/myDecMap90")

	expect := myDecMap90{
		"world": nil,
		"hello": {true, false, true},
	}

	packet := []byte{sdtReg, 0x13, 0x88,
		sdtReg,
		0, 0, 0, 2,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		sdtReg,
		0, 0, 0, 3,
		1, 0, 1, // true, false, true
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		sdtNil,
	}

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeArray(t *testing.T) {
	expect := [2]string{
		"hello", "world",
	}
	packet := []byte{sdtType, 0, 6,
		sdtArray, 0, 0, 0, 2,
		sdtString,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeArrayZero(t *testing.T) {
	expect := [0]string{}
	packet := []byte{sdtType, 0, 6,
		sdtArray, 0, 0, 0, 0,
		sdtString,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceArray(t *testing.T) {
	expect := [][2]string{
		{"hello", "world"},
	}
	packet := []byte{sdtType, 0, 7,
		sdtSlice,
		sdtArray, 0, 0, 0, 2,
		sdtString,
		sdtSlice,
		0, 0, 0, 1,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeSliceAnyArray(t *testing.T) {
	expect := []any{
		nil,
		[2]string{"hello", "world"},
		nil,
	}
	packet := []byte{sdtType, 0, 2,
		sdtSlice,
		sdtAny,

		sdtSlice,
		0, 0, 0, 3,

		sdtNil,

		sdtType, 0, 6,
		sdtArray, 0, 0, 0, 2,
		sdtString,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"

		sdtNil,
	}

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

type myArrayDec [2]string

func TestDecodeRegArray(t *testing.T) {
	expect := myArrayDec{"hello", "world"}

	packet := []byte{sdtReg,
		0, 38, // len of the type name #github.com/sllt/sparrow/net/sdf/myArrayDec
		0x23, 0x65, 0x72, 0x67, 0x6f, 0x2e, 0x73, 0x65,
		0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65,
		0x72, 0x67, 0x6f, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
		0x65, 0x64, 0x66, 0x2f, 0x6d, 0x79, 0x41, 0x72,
		0x72, 0x61, 0x79, 0x44, 0x65, 0x63,

		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	RegisterTypeOf(expect)

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeArrayReg(t *testing.T) {
	// type myArrayStr is declared in encode_test.go (before the TestEncodeArrayReg test)
	expect := [2]myArrayStr{"hello", "world"}

	packet := []byte{sdtType, 0, 46,
		sdtArray, 0, 0, 0, 2,
		sdtReg,
		0, 38, // len of the type name #github.com/sllt/sparrow/net/sdf/myArrayStr
		0x23, 0x65, 0x72, 0x67, 0x6f, 0x2e, 0x73, 0x65,
		0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65,
		0x72, 0x67, 0x6f, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
		0x65, 0x64, 0x66, 0x2f, 0x6d, 0x79, 0x41, 0x72,
		0x72, 0x61, 0x79, 0x53, 0x74, 0x72,

		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
	}

	RegisterTypeOf(expect[0])

	value, _, err := Decode(packet, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}

func TestDecodeRegArrayRegArray(t *testing.T) {
	type myDecArrayMyStr1 [2]string
	type myDecArrayArray1 [3]myDecArrayMyStr1

	expect := myDecArrayArray1{
		{"hello", "world"},
		{"", ""},
		{"world", "hello"},
	}

	packet := []byte{sdtReg, 0x13, 0x88,
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, 0,
		0, 0,
		0, 5, // len of value "world"
		0x77, 0x6f, 0x72, 0x6c, 0x64, // "world"
		0, 5, // len of value "hello"
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // "hello"
	}

	RegisterTypeOf(myDecArrayMyStr1{})
	RegisterTypeOf(myDecArrayArray1{})

	regCache := new(sync.Map)
	regCache.Store(uint16(5000), "#github.com/sllt/sparrow/net/sdf/myDecArrayArray1")

	value, _, err := Decode(packet, Options{RegCache: regCache})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, expect) {
		fmt.Printf("exp %#v\n", expect)
		fmt.Printf("got %#v\n", value)
		t.Fatal("incorrect value")
	}
}
