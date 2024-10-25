package sdf

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sllt/sparrow/app/system/inspect"
	"github.com/sllt/sparrow/gen"
)

var (
	// register generic Sparrow Framework types for the networking
	genTypes = []any{

		gen.Env(""),
		gen.LogLevel(0),
		gen.ProcessState(0),
		gen.MetaState(0),
		gen.NetworkMode(0),
		gen.MessagePriority(0),
		gen.CompressionType(""),
		gen.CompressionLevel(0),
		gen.ApplicationMode(0),
		gen.ApplicationState(0),

		gen.Version{},

		gen.ApplicationDepends{},

		gen.LoggerInfo{},
		gen.NodeInfo{},
		gen.Compression{},
		gen.ProcessFallback{},
		gen.MailboxQueues{},
		gen.ProcessInfo{},
		gen.ProcessShortInfo{},
		gen.ProcessOptions{},
		gen.ProcessOptionsExtra{},
		gen.ApplicationOptions{},
		gen.ApplicationOptionsExtra{},
		gen.MetaInfo{},

		gen.NetworkFlags{},
		gen.NetworkProxyFlags{},
		gen.NetworkSpawnInfo{},
		gen.NetworkApplicationStartInfo{},
		gen.RemoteNodeInfo{},
		gen.RouteInfo{},
		gen.ProxyRouteInfo{},
		gen.Route{},
		gen.ApplicationRoute{},
		gen.ProxyRoute{},
		gen.RegisterRoutes{},
		gen.RegistrarInfo{},
		gen.AcceptorInfo{},
		gen.NetworkInfo{},
		gen.MessageEvent{},
		gen.MessageEventStart{},
		gen.MessageEventStop{},

		// inspector messages

		inspect.RequestInspectNode{},
		inspect.ResponseInspectNode{},
		inspect.MessageInspectNode{},

		inspect.RequestInspectNetwork{},
		inspect.ResponseInspectNetwork{},
		inspect.MessageInspectNetwork{},

		inspect.RequestInspectConnection{},
		inspect.ResponseInspectConnection{},
		inspect.MessageInspectConnection{},

		inspect.RequestInspectProcessList{},
		inspect.ResponseInspectProcessList{},
		inspect.MessageInspectProcessList{},

		inspect.RequestInspectLog{},
		inspect.ResponseInspectLog{},
		inspect.MessageInspectLogNode{},
		inspect.MessageInspectLogNetwork{},
		inspect.MessageInspectLogProcess{},
		inspect.MessageInspectLogMeta{},

		inspect.RequestInspectProcess{},
		inspect.ResponseInspectProcess{},
		inspect.MessageInspectProcess{},

		inspect.RequestInspectProcessState{},
		inspect.ResponseInspectProcessState{},
		inspect.MessageInspectProcessState{},

		inspect.RequestInspectMeta{},
		inspect.ResponseInspectMeta{},
		inspect.MessageInspectMeta{},

		inspect.RequestInspectMetaState{},
		inspect.ResponseInspectMetaState{},
		inspect.MessageInspectMetaState{},

		inspect.RequestDoSend{},
		inspect.ResponseDoSend{},

		inspect.RequestDoSendMeta{},
		inspect.ResponseDoSendMeta{},

		inspect.RequestDoSendExit{},
		inspect.ResponseDoSendExit{},

		inspect.RequestDoSendExitMeta{},
		inspect.ResponseDoSendExitMeta{},

		inspect.RequestDoKill{},
		inspect.ResponseDoKill{},

		inspect.RequestDoSetLogLevel{},
		inspect.RequestDoSetLogLevelProcess{},
		inspect.RequestDoSetLogLevelMeta{},
		inspect.ResponseDoSetLogLevel{},
	}

	// register standard errors of the Sparrow Framework
	genErrors = []error{
		gen.ErrIncorrect,
		gen.ErrTimeout,
		gen.ErrUnsupported,
		gen.ErrUnknown,
		gen.ErrNameUnknown,
		gen.ErrNotAllowed,
		gen.ErrProcessUnknown,
		gen.ErrProcessTerminated,
		gen.ErrMetaUnknown,
		gen.ErrApplicationUnknown,
		gen.ErrTaken,
		gen.TerminateReasonNormal,
		gen.TerminateReasonShutdown,
		gen.TerminateReasonKill,
		gen.TerminateReasonPanic,
	}
)

func init() {
	//
	// encoders
	//
	encoders.Store(reflect.TypeOf(gen.PID{}), &encoder{Prefix: []byte{sdtPID}, Encode: encodePID})
	encoders.Store(reflect.TypeOf(gen.ProcessID{}), &encoder{Prefix: []byte{sdtProcessID}, Encode: encodeProcessID})
	encoders.Store(reflect.TypeOf(gen.Ref{}), &encoder{Prefix: []byte{sdtRef}, Encode: encodeRef})
	encoders.Store(reflect.TypeOf(gen.Alias{}), &encoder{Prefix: []byte{sdtAlias}, Encode: encodeAlias})
	encoders.Store(reflect.TypeOf(gen.Event{}), &encoder{Prefix: []byte{sdtEvent}, Encode: encodeEvent})
	encoders.Store(reflect.TypeOf(true), &encoder{Prefix: []byte{sdtBool}, Encode: encodeBool})
	encoders.Store(reflect.TypeOf(gen.Atom("atom")), &encoder{Prefix: []byte{sdtAtom}, Encode: encodeAtom})
	encoders.Store(reflect.TypeOf("string"), &encoder{Prefix: []byte{sdtString}, Encode: encodeString})
	encoders.Store(reflect.TypeOf(int(0)), &encoder{Prefix: []byte{sdtInt}, Encode: encodeInt})
	encoders.Store(reflect.TypeOf(int8(0)), &encoder{Prefix: []byte{sdtInt8}, Encode: encodeInt8})
	encoders.Store(reflect.TypeOf(int16(0)), &encoder{Prefix: []byte{sdtInt16}, Encode: encodeInt16})
	encoders.Store(reflect.TypeOf(int32(0)), &encoder{Prefix: []byte{sdtInt32}, Encode: encodeInt32})
	encoders.Store(reflect.TypeOf(int64(0)), &encoder{Prefix: []byte{sdtInt64}, Encode: encodeInt64})
	encoders.Store(reflect.TypeOf(uint(0)), &encoder{Prefix: []byte{sdtUint}, Encode: encodeUint})
	encoders.Store(reflect.TypeOf(uint8(0)), &encoder{Prefix: []byte{sdtUint8}, Encode: encodeUint8})
	encoders.Store(reflect.TypeOf(uint16(0)), &encoder{Prefix: []byte{sdtUint16}, Encode: encodeUint16})
	encoders.Store(reflect.TypeOf(uint32(0)), &encoder{Prefix: []byte{sdtUint32}, Encode: encodeUint32})
	encoders.Store(reflect.TypeOf(uint64(0)), &encoder{Prefix: []byte{sdtUint64}, Encode: encodeUint64})
	encoders.Store(reflect.TypeOf([]byte(nil)), &encoder{Prefix: []byte{sdtBinary}, Encode: encodeBinary})
	encoders.Store(reflect.TypeOf(float32(0.0)), &encoder{Prefix: []byte{sdtFloat32}, Encode: encodeFloat32})
	encoders.Store(reflect.TypeOf(float64(0.0)), &encoder{Prefix: []byte{sdtFloat64}, Encode: encodeFloat64})
	encoders.Store(reflect.TypeOf(time.Time{}), &encoder{Prefix: []byte{sdtTime}, Encode: encodeTime})
	encoders.Store(anyType, &encoder{Prefix: []byte{sdtAny}, Encode: encodeAny})

	// error types
	encoders.Store(errType, &encoder{Prefix: []byte{sdtError}, Encode: encodeError})
	encoders.Store(reflect.TypeOf(fmt.Errorf("")), &encoder{Prefix: []byte{sdtError}, Encode: encodeError})
	// wrapped error has a different type
	encoders.Store(reflect.TypeOf(fmt.Errorf("%w", nil)), &encoder{Prefix: []byte{sdtError}, Encode: encodeError})

	//
	// decoders
	//
	decPID := &decoder{reflect.TypeOf(gen.PID{}), decodePID}
	decoders.Store(sdtPID, decPID)
	decoders.Store(decPID.Type, decPID)

	decProcessID := &decoder{reflect.TypeOf(gen.ProcessID{}), decodeProcessID}
	decoders.Store(sdtProcessID, decProcessID)
	decoders.Store(decProcessID.Type, decProcessID)

	decRef := &decoder{reflect.TypeOf(gen.Ref{}), decodeRef}
	decoders.Store(sdtRef, decRef)
	decoders.Store(decRef.Type, decRef)

	decAlias := &decoder{reflect.TypeOf(gen.Alias{}), decodeAlias}
	decoders.Store(sdtAlias, decAlias)
	decoders.Store(decAlias.Type, decAlias)

	decEvent := &decoder{reflect.TypeOf(gen.Event{}), decodeEvent}
	decoders.Store(sdtEvent, decEvent)
	decoders.Store(decEvent.Type, decEvent)

	decTime := &decoder{reflect.TypeOf(time.Time{}), decodeTime}
	decoders.Store(sdtTime, decTime)
	decoders.Store(decTime.Type, decTime)

	decBool := &decoder{reflect.TypeOf(true), decodeBool}
	decoders.Store(sdtBool, decBool)
	decoders.Store(decBool.Type, decBool)

	decAtom := &decoder{reflect.TypeOf(gen.Atom("atom")), decodeAtom}
	decoders.Store(sdtAtom, decAtom)
	decoders.Store(decAtom.Type, decAtom)

	decString := &decoder{reflect.TypeOf("string"), decodeString}
	decoders.Store(sdtString, decString)
	decoders.Store(decString.Type, decString)

	decInt := &decoder{reflect.TypeOf(int(0)), decodeInt}
	decoders.Store(sdtInt, decInt)
	decoders.Store(decInt.Type, decInt)

	decInt8 := &decoder{reflect.TypeOf(int8(0)), decodeInt8}
	decoders.Store(sdtInt8, decInt8)
	decoders.Store(decInt8.Type, decInt8)

	decInt16 := &decoder{reflect.TypeOf(int16(0)), decodeInt16}
	decoders.Store(sdtInt16, decInt16)
	decoders.Store(decInt16.Type, decInt16)

	decInt32 := &decoder{reflect.TypeOf(int32(0)), decodeInt32}
	decoders.Store(sdtInt32, decInt32)
	decoders.Store(decInt32.Type, decInt32)

	decInt64 := &decoder{reflect.TypeOf(int64(0)), decodeInt64}
	decoders.Store(sdtInt64, decInt64)
	decoders.Store(decInt64.Type, decInt64)

	decUint := &decoder{reflect.TypeOf(uint(0)), decodeUint}
	decoders.Store(sdtUint, decUint)
	decoders.Store(decUint.Type, decUint)

	decUint8 := &decoder{reflect.TypeOf(uint8(0)), decodeUint8}
	decoders.Store(sdtUint8, decUint8)
	decoders.Store(decUint8.Type, decUint8)

	decUint16 := &decoder{reflect.TypeOf(uint16(0)), decodeUint16}
	decoders.Store(sdtUint16, decUint16)
	decoders.Store(decUint16.Type, decUint16)

	decUint32 := &decoder{reflect.TypeOf(uint32(0)), decodeUint32}
	decoders.Store(sdtUint32, decUint32)
	decoders.Store(decUint32.Type, decUint32)

	decUint64 := &decoder{reflect.TypeOf(uint64(0)), decodeUint64}
	decoders.Store(sdtUint64, decUint64)
	decoders.Store(decUint64.Type, decUint64)

	decBinary := &decoder{reflect.TypeOf([]byte(nil)), decodeBinary}
	decoders.Store(sdtBinary, decBinary)
	decoders.Store(decBinary.Type, decBinary)

	decFloat32 := &decoder{reflect.TypeOf(float32(0.0)), decodeFloat32}
	decoders.Store(sdtFloat32, decFloat32)
	decoders.Store(decFloat32.Type, decFloat32)

	decFloat64 := &decoder{reflect.TypeOf(float64(0.0)), decodeFloat64}
	decoders.Store(sdtFloat64, decFloat64)
	decoders.Store(decFloat64.Type, decFloat64)

	decAny := &decoder{anyType, decodeAny}
	decoders.Store(sdtAny, decAny)
	decoders.Store(anyType, decAny)

	decErr := &decoder{errType, decodeError}
	decoders.Store(sdtError, decErr)
	decoders.Store(decErr.Type, decErr)

	for _, t := range genTypes {
		err := RegisterTypeOf(t)
		if err == nil || err == gen.ErrTaken {
			continue
		}
		panic(err)
	}

	for _, e := range genErrors {
		err := RegisterError(e)
		if err == nil || err == gen.ErrTaken {
			continue
		}
		panic(err)
	}
}
