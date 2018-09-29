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

type errorStack struct {
	msg     string
	cause   error
	callers []uintptr
}

func (e *errorStack) Error() string {
	return e.msg
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorStack{msg: text, cause: nil, callers: trace(3)}
}

// NewWithCause returns an error that formats as the given text encapsulating a cause.
func NewWithCause(text string, cause error) error {
	return &errorStack{msg: text, cause: cause, callers: trace(3)}
}

// Cause simply returns the cause of the specified error if there is one, otherwise nil.
func Cause(err error) error {
	if err != nil {
		if e, ok := err.(*errorStack); ok {
			return e.cause
		}
	}
	return nil
}

// Trace returns the runtime frames of the specified error if there is one, otherwise nil.
func Trace(err error) *runtime.Frames {
	if err != nil {
		if e, ok := err.(*errorStack); ok {
			return runtime.CallersFrames(e.callers)
		}
	}
	return nil
}

// String returns a full string of the specified error stack if there is one, otherwise nil.
// It also can include stack trace for each error in the stack.
func String(err error, withTrace bool) string {
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

		if withTrace {
			if frames := Trace(e); frames != nil {
				for {
					frame, more := frames.Next()
					buffer.WriteString(fmt.Sprintf("\n\t^ %v(%v:%v)", frame.Function, frame.File, frame.Line))
					if !more {
						break
					}
				}
			}
		}

		e = Cause(e)
	}

	return buffer.String()
}

func trace(skip int) []uintptr {
	var callers [MaxStackLength]uintptr
	n := runtime.Callers(skip, callers[:])
	return callers[0:n]
}
