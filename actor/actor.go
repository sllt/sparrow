package actor

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/sllt/sparrow/gen"
	"github.com/sllt/sparrow/lib"
)

// ActorBehavior interface
type ActorBehavior interface {
	gen.ProcessBehavior

	// Init invoked on a spawn Actor for the initializing.
	Init(args ...any) error

	// HandleMessage invoked if Actor received a message sent with gen.Process.Send(...).
	// Non-nil value of the returning error will cause termination of this process.
	// To stop this process normally, return gen.TerminateReasonNormal
	// or any other for abnormal termination.
	HandleMessage(from gen.PID, message any) error

	// HandleCall invoked if Actor got a synchronous request made with gen.Process.Call(...).
	// Return nil as a result to handle this request asynchronously and
	// to provide the result later using the gen.Process.SendResponse(...) method.
	HandleCall(from gen.PID, ref gen.Ref, request any) (any, error)

	// Terminate invoked on a termination process
	Terminate(reason error)

	// HandleMessageName invoked if split handling was enabled using SetSplitHandle(true)
	// and message has been sent by name
	HandleMessageName(name gen.Atom, from gen.PID, message any) error
	HandleMessageAlias(alias gen.Alias, from gen.PID, message any) error
	HandleCallName(name gen.Atom, from gen.PID, ref gen.Ref, request any) (any, error)
	HandleCallAlias(alias gen.Alias, from gen.PID, ref gen.Ref, request any) (any, error)

	// HandleLog invoked on a log message if this process was added as a logger.
	HandleLog(message gen.MessageLog) error

	// HandleEvent invoked on an event message if this process got subscribed on
	// this event using gen.Process.LinkEvent or gen.Process.MonitorEvent
	HandleEvent(message gen.MessageEvent) error

	// HandleInspect invoked on the request made with gen.Process.Inspect(...)
	HandleInspect(from gen.PID, item ...string) map[string]string
}

// Actor implementats ProcessBehavior interface and provides callbacks for
// - initialization
// - handling messages/requests.
// - termination
// All callbacks of the ActorBehavior are optional for the implementation.
type Actor struct {
	gen.Process

	behavior ActorBehavior
	mailbox  gen.ProcessMailbox

	trap  bool // trap exit
	split bool // split handle callback
}

// SetTrapExit enables/disables the trap on exit requests sent by SendExit(...).
// Enabled trap makes the actor ignore such requests transforming them into
// regular messages (gen.MessageExitPID) except for the request from the parent
// process with the reason gen.TerminateReasonShutdown.
// With disabled trap, actor gracefully terminates by invoking Terminate callback
// with the given reason
func (a *Actor) SetTrapExit(trap bool) {
	a.trap = trap
}

// TrapExit returns whether the trap was enabled on this actor
func (a *Actor) TrapExit() bool {
	return a.trap
}

// SetSplitHandle enables/disables splitting invoke callback  depending on the target type.
// Enabled splitting makes this process to invoke
//   - HandleCall/HandleMessage for the request/message addressed by gen.PID
//   - HandleCallName/HandleMessageName for the request/message addressed by gen.ProcessID
//   - HandleCallAlias/HandleMessageAlias for the request/message addressed by gen.Alias
func (a *Actor) SetSplitHandle(split bool) {
	a.split = split
}

// SplitHandle returns whether the splitting was enabled on this actor
func (a *Actor) SplitHandle() bool {
	return a.split
}

//
// ProcessBehavior implementation
//

// ProcessInit
func (a *Actor) ProcessInit(process gen.Process, args ...any) (rr error) {
	var ok bool

	if a.behavior, ok = process.Behavior().(ActorBehavior); ok == false {
		unknown := strings.TrimPrefix(reflect.TypeOf(process.Behavior()).String(), "*")
		return fmt.Errorf("ProcessInit: not an ActorBehavior %s", unknown)
	}

	if lib.Recover() {
		defer func() {
			if r := recover(); r != nil {
				pc, fn, line, _ := runtime.Caller(2)
				a.Log().Panic("Actor initialization failed. Panic reason: %#v at %s[%s:%d]",
					r, runtime.FuncForPC(pc).Name(), fn, line)
				rr = gen.TerminateReasonPanic
			}
		}()
	}

	a.Process = process
	a.mailbox = process.Mailbox()

	return a.behavior.Init(args...)
}

func (a *Actor) ProcessRun() (rr error) {
	if lib.Recover() {
		defer func() {
			if r := recover(); r != nil {
				pc, fn, line, _ := runtime.Caller(2)
				a.Log().Panic("Actor terminated. Panic reason: %#v at %s[%s:%d]",
					r, runtime.FuncForPC(pc).Name(), fn, line)
				rr = gen.TerminateReasonPanic
			}
		}()
	}

	for {
		if a.State() != gen.ProcessStateRunning {
			return gen.TerminateReasonKill
		}

		message := a.getNextMessage()
		if message == nil {
			return nil
		}

		if err := a.handleMessage(message); err != nil {
			gen.ReleaseMailboxMessage(message)
			return err
		}

		gen.ReleaseMailboxMessage(message)
	}
}

func (a *Actor) getNextMessage() *gen.MailboxMessage {
	for {
		if msg, ok := a.mailbox.Urgent.Pop(); ok {
			return msg.(*gen.MailboxMessage)
		}
		if msg, ok := a.mailbox.System.Pop(); ok {
			return msg.(*gen.MailboxMessage)
		}
		if msg, ok := a.mailbox.Main.Pop(); ok {
			return msg.(*gen.MailboxMessage)
		}
		if msg, ok := a.mailbox.Log.Pop(); ok {
			if reason := a.behavior.HandleLog(msg.(gen.MessageLog)); reason != nil {
				return nil
			}
			continue
		}
		return nil
	}
}

func (a *Actor) handleMessage(message *gen.MailboxMessage) error {
	switch message.Type {
	case gen.MailboxMessageTypeRegular:
		return a.handleRegularMessage(message)
	case gen.MailboxMessageTypeRequest:
		return a.handleRequestMessage(message)
	case gen.MailboxMessageTypeEvent:
		return a.behavior.HandleEvent(message.Message.(gen.MessageEvent))
	case gen.MailboxMessageTypeExit:
		return a.handleExitMessage(message)
	case gen.MailboxMessageTypeInspect:
		a.handleInspectMessage(message)
		return nil
	default:
		return fmt.Errorf("unknown message type: %v", message.Type)
	}
}

func (a *Actor) handleRegularMessage(message *gen.MailboxMessage) error {
	if a.split {
		switch target := message.Target.(type) {
		case gen.Atom:
			return a.behavior.HandleMessageName(target, message.From, message.Message)
		case gen.Alias:
			return a.behavior.HandleMessageAlias(target, message.From, message.Message)
		default:
			return a.behavior.HandleMessage(message.From, message.Message)
		}
	}
	return a.behavior.HandleMessage(message.From, message.Message)
}

func (a *Actor) handleRequestMessage(message *gen.MailboxMessage) error {
	var result any
	var err error

	if a.split {
		switch target := message.Target.(type) {
		case gen.Atom:
			result, err = a.behavior.HandleCallName(target, message.From, message.Ref, message.Message)
		case gen.Alias:
			result, err = a.behavior.HandleCallAlias(target, message.From, message.Ref, message.Message)
		default:
			result, err = a.behavior.HandleCall(message.From, message.Ref, message.Message)
		}
	} else {
		result, err = a.behavior.HandleCall(message.From, message.Ref, message.Message)
	}

	if err != nil {
		if err == gen.TerminateReasonNormal && result != nil {
			a.SendResponse(message.From, message.Ref, result)
		}
		return err
	}

	if result != nil {
		a.SendResponse(message.From, message.Ref, result)
	}

	return nil
}

func (a *Actor) handleExitMessage(message *gen.MailboxMessage) error {
	switch exit := message.Message.(type) {
	case gen.MessageExitPID:
		if a.trap && message.From != a.Parent() {
			return a.handleRegularMessage(message)
		}
		return fmt.Errorf("%s: %w", exit.PID, exit.Reason)
	case gen.MessageExitProcessID:
		if a.trap {
			return a.handleRegularMessage(message)
		}
		return fmt.Errorf("%s: %w", exit.ProcessID, exit.Reason)
	case gen.MessageExitAlias:
		if a.trap {
			return a.handleRegularMessage(message)
		}
		return fmt.Errorf("%s: %w", exit.Alias, exit.Reason)
	case gen.MessageExitEvent:
		if a.trap {
			return a.handleRegularMessage(message)
		}
		return fmt.Errorf("%s: %w", exit.Event, exit.Reason)
	case gen.MessageExitNode:
		if a.trap {
			return a.handleRegularMessage(message)
		}
		return fmt.Errorf("%s: %w", exit.Name, gen.ErrNoConnection)
	default:
		return fmt.Errorf("unknown exit message: %#v", exit)
	}
}

func (a *Actor) handleInspectMessage(message *gen.MailboxMessage) {
	result := a.behavior.HandleInspect(message.From, message.Message.([]string)...)
	a.SendResponse(message.From, message.Ref, result)
}
func (a *Actor) ProcessTerminate(reason error) {
	a.behavior.Terminate(reason)
}

//
// default callbacks for ActorBehavior interface
//

// Init
func (a *Actor) Init(args ...any) error {
	return nil
}

func (a *Actor) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	a.Log().Warning("Actor.HandleCall: unhandled request from %s", from)
	return nil, nil
}

func (a *Actor) HandleMessage(from gen.PID, message any) error {
	a.Log().Warning("Actor.HandleMessage: unhandled message from %s", from)
	return nil
}

func (a *Actor) HandleInspect(from gen.PID, item ...string) map[string]string {
	return nil
}

func (a *Actor) HandleLog(message gen.MessageLog) error {
	a.Log().Warning("Actor.HandleLog: unhandled log message %#v", message)
	return nil
}

func (a *Actor) HandleEvent(message gen.MessageEvent) error {
	a.Log().Warning("Actor.HandleEvent: unhandled event message %#v", message)
	return nil
}

func (a *Actor) Terminate(reason error) {}

func (a *Actor) HandleMessageName(name gen.Atom, from gen.PID, message any) error {
	a.Log().Warning("Actor.HandleMessageName %s: unhandled message from %s", a.Name(), from)
	return nil
}

func (a *Actor) HandleMessageAlias(alias gen.Alias, from gen.PID, message any) error {
	a.Log().Warning("Actor.HandleMessageAlias %s: unhandled message from %s", alias, from)
	return nil
}

func (a *Actor) HandleCallName(name gen.Atom, from gen.PID, ref gen.Ref, request any) (any, error) {
	a.Log().Warning("Actor.HandleCallName %s: unhandled request from %s", a.Name(), from)
	return nil, nil
}

func (a *Actor) HandleCallAlias(alias gen.Alias, from gen.PID, ref gen.Ref, request any) (any, error) {
	a.Log().Warning("Actor.HandleCallAlias %s: unhandled request from %s", alias, from)
	return nil, nil
}
