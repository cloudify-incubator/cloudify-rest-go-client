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
	"encoding/json"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	"log"
	"net/url"
)

type CloudifyNodeInstance struct {
	rest.CloudifyIdWithTenant
	Relationships     []interface{}          `json:"relationships,omitempty"`
	RuntimeProperties map[string]interface{} `json:"runtime_properties,omitempty"`
	State             string                 `json:"state,omitempty"`
	Version           int                    `json:"host_id,version"`
	HostId            string                 `json:"host_id,omitempty"`
	DeploymentId      string                 `json:"deployment_id,omitempty"`
	NodeId            string                 `json:"node_id,omitempty"`
	// TODO describe "scaling_groups" struct
}

func (instance *CloudifyNodeInstance) GetJsonRuntimeProperties() string {
	json_data, err := json.Marshal(instance.RuntimeProperties)
	if err != nil {
		log.Fatal(err)
	}
	return string(json_data)
}

type CloudifyNodeInstances struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata  `json:"metadata"`
	Items    []CloudifyNodeInstance `json:"items"`
}

func (cl *CloudifyClient) GetNodeInstances(params map[string]string) CloudifyNodeInstances {
	var instances CloudifyNodeInstances

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("node-instances?"+values.Encode(), &instances)
	if err != nil {
		log.Fatal(err)
	}

	return instances
}
