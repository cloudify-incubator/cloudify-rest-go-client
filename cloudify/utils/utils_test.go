/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

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

package utils

import (
	"fmt"
	"reflect"
	"testing"
)

type cliTestpair struct {
	input  []string
	output [2][]string
}

var cliTests = []cliTestpair{
	{[]string{"cfy-go", "status"}, [2][]string{{"cfy-go", "status"}, {}}},
	{[]string{"cfy-go", "status", "-user", "admin"}, [2][]string{{"cfy-go", "status"}, {"-user", "admin"}}},
}

func TestCliArgumentsList(t *testing.T) {
	for _, pair := range cliTests {
		args, options := CliArgumentsList(pair.input)

		if !reflect.DeepEqual(args, pair.output[0]) {
			t.Error(
				"For", pair.input,
				"expected", pair.output[0],
				"got", args,
			)
		}

		if !reflect.DeepEqual(options, pair.output[1]) {
			t.Error(
				"For", pair.input,
				"expected", pair.output[1],
				"got", options,
			)
		}

	}
}

func ExampleCliArgumentsList() {
	args, options := CliArgumentsList(
		[]string{"cfy-go", "status", "-user", "admin"})
	fmt.Printf("Args: %v, Options: %v.", args, options)
	// Output: Args: [cfy-go status], Options: [-user admin].
}

func ExampleCliSubArgumentsList() {
	args, options := CliSubArgumentsList([]string{
		"-data", "container-place/data", "-work-dir", "container-place/",
		"--", "cfy", "profile", "use", "local"})
	fmt.Printf("Args: %v, Options: %v.", args, options)
	// Output: Args: [-data container-place/data -work-dir container-place/], Options: [cfy profile use local].
}

func TestInList(t *testing.T) {
	if InList([]string{"a", "b"}, "c") {
		t.Error("Recheck 'c' not in ['a', 'b'] ")
	}
	if !InList([]string{"a", "b"}, "b") {
		t.Error("Recheck 'b' in ['a', 'b'] ")
	}
}

func ExampleInList() {
	if !InList([]string{"a", "b"}, "c") {
		fmt.Print("'c' not in ['a', 'b'] list.")
	}
	// Output: 'c' not in ['a', 'b'] list.
}

func ExamplePrintTable() {
	PrintTable([]string{"1", "2", "3"}, [][]string{{"a", "-b", "c"}, {"d", "e", "-f-"}})
	// Output: +---+----+-----+
	// | 1 | 2  | 3   |
	// +---+----+-----+
	// | a | -b | c   |
	// | d | e  | -f- |
	// +---+----+-----+
}
