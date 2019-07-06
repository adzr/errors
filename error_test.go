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

func TestDescribe(t *testing.T) {
	err := Wrap("error 1", Wrap("error 2", New("error 3")))

	if s := Describe(err, true); len(s) == 0 {
		t.Errorf("error string is empty")
	}
}

func TestMap(t *testing.T) {
	err := Wrap("error 1", Wrap("error 2", New("error 3")))

	if b, _ := json.Marshal(Map(err, true)); len(string(b)) == 0 {
		t.Errorf("error string is empty")
	}
}
