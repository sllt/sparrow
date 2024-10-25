package handshake

import (
	"github.com/sllt/sparrow/gen"
	"github.com/sllt/sparrow/net/sdf"
	"sync"
)

const (
	handshakeName    string = "EHS"
	handshakeRelease string = "R1" // Sparrow Handshake (Rev.1)

	handshakeMagic   byte = 87
	handshakeVersion byte = 1

	defaultPoolSize int = 3
)

var (
	DefaultPoolSize int = 1
)

type MessageHello struct {
	Salt       string
	Digest     string
	DigestCert string
}

type MessageJoin struct {
	Node         gen.Atom
	ConnectionID string
	Salt         string
	Digest       string
}

type MessageIntroduce struct {
	Node     gen.Atom
	Version  gen.Version
	Flags    gen.NetworkFlags
	Creation int64

	MaxMessageSize int

	AtomCache map[uint16]gen.Atom
	RegCache  map[uint16]string
	ErrCache  map[uint16]error
	Digest    string
}

type MessageAccept struct {
	ID         string
	PoolSize   int
	PoolDSN    []string
	Digest     string
	DigestCert string
}

type ConnectionOptions struct {
	PoolSize int
	PoolDSN  []string

	EncodeAtomCache *sync.Map
	EncodeRegCache  *sync.Map
	EncodeErrCache  *sync.Map

	DecodeAtomCache *sync.Map
	DecodeRegCache  *sync.Map
	DecodeErrCache  *sync.Map
}

func init() {
	types := []any{
		MessageHello{},
		MessageJoin{},
		MessageIntroduce{},
		MessageAccept{},
	}

	for _, t := range types {
		err := sdf.RegisterTypeOf(t)
		if err == nil || err == gen.ErrTaken {
			continue
		}
		panic(err)
	}
}
