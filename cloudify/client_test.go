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

package cloudify

import (
	"fmt"
	tests "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/tests"
	"testing"
)

// TestGetVersion - check GetVersion
func TestGetAPIVersion(t *testing.T) {
	var conn tests.FakeClient
	conn.GetResponse = []byte(versionResponce)
	conn.GetError = nil
	cl := ClientFromConnection(&conn)
	version := cl.GetAPIVersion()
	if version != "v3.1" {
		t.Errorf("Recheck unmarshal for 'version' field '%s'", version)
	}
}


func ExampleClient() {
	cl := NewClient("localhost", "admin", "password", "default_tenant")
	fmt.Printf("Version: %s", cl.GetAPIVersion())
	// Output: Version: v3.1
}
