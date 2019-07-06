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

/*
Package errors provides drop in replacement for the standard errors package with stack trace feature.

Brief

This library is a drop in replacement for the default Go error implementation with a feature to nest
errors and provide call stack traces.

Usage

	$ go get -u github.com/adzr/errors

Then, import the package:

  import (
    "github.com/adzr/errors"
  )

Example

  // This is a normal error creation, same as the standard library.
  err := errors.New("this is an error")

  // And here is how to get the error details.
  println(err.Error()) // this will print 'this is an error'
  println(errors.Unwrap(err)) // this will print <nil>

  // This code block...
  frames := errors.StackTrace(err)
	for {
		frame, more := frames.Next()
		fmt.Printf("%v(%v:%v)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
  }
  // will simply print something like this:
  // example.TestErrors($GOPATH/src/github.com/adzr/errors/error_test.go:75)
  // testing.tRunner($GOROOT/src/testing/testing.go:777)
  // runtime.goexit($GOROOT/src/runtime/asm_amd64.s:2361)

  println(errors.Describe(err, false)) // this will print 'this is an error'

  // However, this...
  println(errors.Describe(err, true))
  // will print something like this:
  // this is an error
	//    example.TestErrors($GOPATH/src/github.com/adzr/errors/error_test.go:75)
	//    testing.tRunner($GOROOT/src/testing/testing.go:777)
	//    runtime.goexit($GOROOT/src/runtime/asm_amd64.s:2361)

  // Same functions can be used with errWithUnwrap.
  errors.Describe(errors.Wrap("this is an error", errors.New("this is the cause")), true)
*/
package errors
