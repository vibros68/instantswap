package errors

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

type error interface {
	Error() string
}

// Kind describes the class of error.
type Kind int

// Error kinds.
const (
	Other               Kind = iota // Unclassified error -- does not appear in error strings 	0
	Bug                             // Error is known to be a result of our bug 				1
	Invalid                         // Invalid operation 										2
	IO                              // I/O error 												3
	Exist                           // Item already exists 										4
	NotExist                        // Item does not exist 										5
	Encoding                        // Invalid encoding 										6
	InsufficientBalance             // Insufficient balance 									7
	Timeout                         // Timeout error 											8
	Connection                      // Connection error 										9
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "unclassified error"
	case Bug:
		return "internal bug"
	case Invalid:
		return "invalid operation"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case Encoding:
		return "invalid encoding"
	case InsufficientBalance:
		return "insufficient balance"
	case Timeout:
		return ""
	case Connection:
		return "connection error"
	default:
		return "unknown error kind"
	}
}

type Error struct {
	Kind  Kind
	Err   error
	Sleep int64 //this is used for sleep intervals in scripts (sec, mins, hours, etc)
}

var Separator = ":\n\t"

func (e *Error) Error() string {
	var b strings.Builder

	// Record the last added fields to the string to avoid duplication.
	var last Error

	for {
		pad := false // whether to pad/separate next field
		if e.Kind != 0 && e.Kind != last.Kind {
			if pad {
				b.WriteString(": ")
			}
			b.WriteString(e.Kind.String())
			pad = true
			last.Kind = e.Kind
		}
		if e.Err == nil {
			break
		}
		if err, ok := e.Err.(*Error); ok {
			if pad {
				b.WriteString(Separator)
			}
			e = err
			continue
		}
		if pad {
			b.WriteString(": ")
		}
		b.WriteString(e.Err.Error())
		break
	}

	s := b.String()
	if s == "" {
		return Other.String()
	}
	return s
}

// New creates a simple error from a string.  New is identical to "errors".New
// from the standard library.
func New(text string) error {
	return errors.New(text)
}

// Errorf creates a simple error from a format string and arguments.  Errorf is
// identical to "fmt".Errorf from the standard library.
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Is returns whether err is of type *Error and has a matching kind in err or
// any nested errors.  Does not match against the Other kind.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	return Is(kind, e.Err)
}
func HandleError(errStr string, errIn error) (errOut Error) {
	switch errIn := errIn.(type) {
	case *url.Error, net.Error:
		if errStr != "" {
			errStr = strings.Replace(errStr, "%", "%%", -1)
			errOut.Err = errors.New(errStr)
		} else {
			errOut.Err = errIn
		}
		errOut.Kind = 9 //Connection Error
		return errOut
	case *Error:
		if errStr != "" {
			errStr = strings.Replace(errStr, "%", "%%", -1)
			errOut.Err = errors.New(errStr)
		} else {
			errOut.Err = errIn
		}
		errOut.Kind = errIn.Kind
		return
	default:
		if errStr != "" {
			errStr = strings.Replace(errStr, "%", "%%", -1)
			errOut.Err = errors.New(errStr)
		} else {
			errOut.Err = errIn
		}
		return errOut
	}
}
