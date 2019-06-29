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
	"encoding/json"
	"errors"
	"regexp"
	"testing"
)

func TestUnwrap(t *testing.T) {
	if err := Unwrap(Wrap("error 1", New("error 2"))); err.Error() != "error 2" {
		t.Errorf("error cause not matching the expected value")
	}

	if err := Unwrap(New("error 1")); err != nil {
		t.Errorf("error expected cause must be nil")
	}

	if err := Unwrap(errors.New("error 1")); err != nil {
		t.Errorf("error expected cause must be nil")
	}
}

func TestTrace(t *testing.T) {
	if trc := StackTrace(Wrap("error 1", New("error 2"))); trc == nil {
		t.Errorf("error trace must not be nil")
	}

	if trc := StackTrace(New("error 1")); trc == nil {
		t.Errorf("error trace must not be nil")
	}

	if trc := StackTrace(errors.New("error 1")); trc != nil {
		t.Errorf("error trace must be nil")
	}
}

func TestError(t *testing.T) {
	if err := Wrap("error 1", New("error 2")); err.Error() != "error 1" {
		t.Errorf("error string not matching the expected value")
	}

	if err := New("error 1"); err.Error() != "error 1" {
		t.Errorf("error string not matching the expected value")
	}
}

func TestString(t *testing.T) {
	err := Wrap("error 1", Wrap("error 2", New("error 3")))

	if matches, _ := regexp.MatchString("\\A((\n \\-\\>)??.+?(\\n\\t\\^ .+?\\(.*?\\:\\d+?\\))+?){3}\\z", String(err, true)); !matches {
		t.Errorf("error string not matching the expected value")
	}
}

func TestMap(t *testing.T) {
	err := Wrap("error 1", Wrap("error 2", New("error 3")))

	s := `{"cause":{"cause":{"cause":"null","msg":"error 3","trace":[{"file":"/home/adzr/Documents/code/foss/errors/error_test.go","func":"github.com/adzr/errors.TestMap","line":73},{"file":"/home/adzr/Tools/packages/go/src/testing/testing.go","func":"testing.tRunner","line":865},{"file":"/home/adzr/Tools/packages/go/src/runtime/asm_amd64.s","func":"runtime.goexit","line":1337}]},"msg":"error 2","trace":[{"file":"/home/adzr/Documents/code/foss/errors/error_test.go","func":"github.com/adzr/errors.TestMap","line":73},{"file":"/home/adzr/Tools/packages/go/src/testing/testing.go","func":"testing.tRunner","line":865},{"file":"/home/adzr/Tools/packages/go/src/runtime/asm_amd64.s","func":"runtime.goexit","line":1337}]},"msg":"error 1","trace":[{"file":"/home/adzr/Documents/code/foss/errors/error_test.go","func":"github.com/adzr/errors.TestMap","line":73},{"file":"/home/adzr/Tools/packages/go/src/testing/testing.go","func":"testing.tRunner","line":865},{"file":"/home/adzr/Tools/packages/go/src/runtime/asm_amd64.s","func":"runtime.goexit","line":1337}]}`

	if b, _ := json.Marshal(Map(err, true)); string(b) != s {
		t.Errorf("error string '%v' not matching the expected value '%v'", string(b), s)
	}
}
