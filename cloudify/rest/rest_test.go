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

package rest

import (
	"fmt"
	"testing"
)

func TestNewClientDebug(t *testing.T) {
	cl := NewClient("localhost", "admin", "password", "default_tenant")
	if cl.GetDebug() {
		t.Error("Debug must be disabled.")
	}

	cl.SetDebug(true)
	if !cl.GetDebug() {
		t.Error("Debug must be enabled.")
	}

	cl.SetDebug(false)
	if cl.GetDebug() {
		t.Error("Debug must be disabled.")
	}
}

func ExampleNewClient() {
	var cl = NewClient("localhost", "admin", "password", "default_tenant")
	fmt.Printf("Debug: %+v", cl.GetDebug())
	// Output: Debug: false
}
