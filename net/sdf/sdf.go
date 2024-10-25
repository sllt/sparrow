package sdf

import (
	"io"
	"sync"
)

// Options for encoding/decoding
type Options struct {
	AtomCache   *sync.Map // atom => id (encoding), id => atom (decoding)
	AtomMapping *sync.Map // atomX => atomY (encoding/decoding)
	RegCache    *sync.Map // type/name => id (encoding), id => type (for decoding)
	ErrCache    *sync.Map // error => id (for encoder), id => error (for decoder)
	Cache       *sync.Map // common cache (caching reflect.Type => encoder, string([]byte) => decoder)
}

const (
	sdtType = byte(130) // 0x82
	sdtReg  = byte(131) // 0x83
	sdtAny  = byte(132) // 0x84

	sdtAtom    = byte(140) // 0x8c
	sdtString  = byte(141) // 0x8d
	sdtBinary  = byte(142) // 0x8e
	sdtFloat32 = byte(143) // 0x8f
	sdtFloat64 = byte(144) // 0x90
	sdtBool    = byte(145) // 0x91
	sdtInt8    = byte(146) // 0x92
	sdtInt16   = byte(147) // 0x93
	sdtInt32   = byte(148) // 0x94
	sdtInt64   = byte(149) // 0x95
	sdtInt     = byte(150) // 0x96
	sdtUint8   = byte(151) // 0x97
	sdtUint16  = byte(152) // 0x98
	sdtUint32  = byte(153) // 0x99
	sdtUint64  = byte(154) // 0x9a
	sdtUint    = byte(155) // 0x9b
	sdtError   = byte(156) // 0x9c
	sdtSlice   = byte(157) // 0x9d
	sdtArray   = byte(158) // 0x9e
	sdtMap     = byte(159) // 0x9f

	sdtPID       = byte(170) // 0xaa
	sdtProcessID = byte(171) // 0xab
	sdtAlias     = byte(172) // 0xac
	sdtEvent     = byte(173) // 0xad
	sdtRef       = byte(174) // 0xae
	sdtTime      = byte(175) // 0xaf

	sdtNil = byte(255) // 0xff
)

type Marshaler interface {
	MarshalSDF(io.Writer) error
}

type Unmarshaler interface {
	UnmarshalSDF([]byte) error
}
