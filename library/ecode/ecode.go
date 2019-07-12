package ecode

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

var (
	_messages = make(map[int]string)
	_codes    = map[int]struct{}{}
)

// New new a ecode.Codes by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(code int, message string) Code {
	if code <= 0 {
		panic("business ecode must greater than zero")
	}
	return add(code, message)
}

func add(code int, message string) Code {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", code))
	}
	_codes[code] = struct{}{}
	_messages[code] = message
	return Int(code)
}

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
}

// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	if msg, ok := _messages[e.Code()]; ok {
		return msg
	}
	return e.Error()
}

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}
