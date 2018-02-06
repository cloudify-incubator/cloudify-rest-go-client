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
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	"net/url"
)

// Tenant - information about cloudify tenant
type Tenant struct {
	Name         string `json:"name"`
	Users        int  `json:"users"`
	Groups       int  `json:"groups"`
}

// Tenants - cloudify response with tenants list
type Tenants struct {
	rest.BaseMessage
	Metadata rest.Metadata `json:"metadata"`
	Items    []Tenant      `json:"items"`
}

// GetTenants - get tenants list filtered by params
func (cl *Client) GetTenants(params map[string]string) (*Tenants, error) {
	var tenants Tenants

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("tenants?"+values.Encode(), &tenants)
	if err != nil {
		return nil, err
	}

	return &tenants, nil
}
