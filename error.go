/*
Copyright 2018 Ahmed Zaher

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

const (
	// MaxStackLength is the maximum count of stacked errors that can be
	// processed by the libarary's functions.
	MaxStackLength int = 32
)

// A Wrapper is an error implementation
// wrapping context around another error.
type Wrapper interface {
	// Unwrap returns the next error in the error chain.
	// If there is no next error, Unwrap returns nil.
	Unwrap() error
}

func trace(skip int) []uintptr {
	var callers [MaxStackLength]uintptr
	n := runtime.Callers(skip, callers[:])
	return callers[0:n]
}

type errorWrapper struct {
	msg     string
	wrapped error
	callers []uintptr
}

func (w *errorWrapper) Error() string {
	return w.msg
}

func (w *errorWrapper) Unwrap() error {
	return w.wrapped
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorWrapper{msg: text, wrapped: nil, callers: trace(3)}
}

// Wrap returns an error that formats as the given text encapsulating a cause.
func Wrap(text string, cause error) error {
	return &errorWrapper{msg: text, wrapped: cause, callers: trace(3)}
}

// Unwrap returns the wrapped error if there is one, otherwise nil.
func Unwrap(e error) error {

	if e != nil {
		if err, ok := e.(*errorWrapper); ok {
			return err.Unwrap()
		}
	}

	return nil
}

// StackTrace returns the runtime frames of the specified error if there is one, otherwise nil.
func StackTrace(err error) *runtime.Frames {
	if err != nil {
		if e, ok := err.(*errorWrapper); ok {
			return runtime.CallersFrames(e.callers)
		}
	}
	return nil
}

// String returns a full string of the specified error stack if there is one, otherwise empty string.
// It also can include stack trace for each error in the stack.
func String(err error, stackTrace bool) string {
	var buffer bytes.Buffer
	var e = err

	for i := 0; i < MaxStackLength; i++ {
		if e == nil {
			break
		}

		if i > 0 {
			buffer.WriteString("\n -> ")
		}

		buffer.WriteString(e.Error())

		if stackTrace {
			if frames := StackTrace(e); frames != nil {
				for {
					frame, more := frames.Next()
					buffer.WriteString(fmt.Sprintf("\n\t^ %v(%v:%v)", frame.Function, frame.File, frame.Line))
					if !more {
						break
					}
				}
			}
		}

		e = Unwrap(e)
	}

	return buffer.String()
}

// Map converts the error stack into a nested map
// so it can be easily marshalled to JSON.
func Map(err error, stackTrace bool) interface{} {

	if err == nil {
		return "null"
	}

	obj := map[string]interface{}{
		"msg":   err.Error(),
		"cause": Map(Unwrap(err), stackTrace),
	}

	if stackTrace {
		t := make([]map[string]interface{}, 0)

		if frames := StackTrace(err); frames != nil {
			for {
				frame, more := frames.Next()
				t = append(t, map[string]interface{}{"func": frame.Function, "file": frame.File, "line": frame.Line})
				if !more {
					break
				}
			}
		}

		obj["trace"] = t
	}

	return obj
}
