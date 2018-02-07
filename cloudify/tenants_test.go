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

const tenantsResponce = `{
  "items": [
    {
      "name": "default_tenant",
      "groups": 0,
      "users": 1
    },
    {
      "name": "examples_tenant",
      "groups": 0,
      "users": 0
    },
    {
      "name": "examples",
      "groups": 0,
      "users": 0
    }
  ],
  "metadata": {
    "pagination": {
      "total": 3,
      "offset": 0,
      "size": 100
    }
  }
}`

// TestGetTenants - check GetTenants
func TestGetTenants(t *testing.T) {
	var conn tests.FakeClient
	conn.GetResponse = []byte(tenantsResponce)
	conn.GetError = nil
	cl := ClientFromConnection(&conn)

	var params = map[string]string{}
	params["_size"] = "100"
	params["_offset"] = "0"

	fmt.Print()
	tenants, err := cl.GetTenants(params)
	if err != nil {
		t.Error("Recheck error reporting")
	}

	lines := len(tenants.Items)
	total := int(tenants.Metadata.Pagination.Total)
	offset := int(tenants.Metadata.Pagination.Offset)
	size := int(tenants.Metadata.Pagination.Size)

	tests.AssertEqual(t, 3, lines, "The number of tenants should be 3")
	tests.AssertEqual(t, 3, total, "The number of tenants should be 3")
	tests.AssertEqual(t, 0, offset, "The offset of tenants page should be 0")
	tests.AssertEqual(t, 100, size, "The size of of tenants page should be 100")

}
