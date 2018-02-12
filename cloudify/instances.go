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
)

// NodeInstanceScalingGroup - short information(ID+Name) about scaling group related to instance
type NodeInstanceScalingGroup struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// NodeInstance - cloudify node instance struct
type NodeInstance struct {
	rest.ObjectIDWithTenant
	Relationships     []interface{}              `json:"relationships,omitempty"`
	RuntimeProperties map[string]interface{}     `json:"runtime_properties,omitempty"`
	State             string                     `json:"state,omitempty"`
	Version           int                        `json:"version,omitempty"`
	HostID            string                     `json:"host_id,omitempty"`
	DeploymentID      string                     `json:"deployment_id,omitempty"`
	NodeID            string                     `json:"node_id,omitempty"`
	ScalingGroups     []NodeInstanceScalingGroup `json:"scaling_groups,omitempty"`
}

// GetJSONRuntimeProperties - instance runtime properties returnd as json string
func (instance *NodeInstance) GetJSONRuntimeProperties() (string, error) {
	jsonData, err := json.Marshal(instance.RuntimeProperties)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// NodeInstances - cloudify manager response with list instances
type NodeInstances struct {
	rest.BaseMessage
	Metadata rest.Metadata  `json:"metadata"`
	Items    []NodeInstance `json:"items"`
}

// GetNodeInstances - Get all node instances
func (cl *Client) GetNodeInstances(params map[string]string) (*NodeInstances, error) {
	var instances NodeInstances

	values := cl.stringMapToURLValue(params)

	err := cl.Get("node-instances?"+values.Encode(), &instances)
	if err != nil {
		return nil, err
	}

	return &instances, nil
}

// AllAreStarted - check that all instances in list are started
func (ni *NodeInstances) AllAreStarted() bool {
	// check that all nodes on same hostID started
	for _, nodeInstance := range ni.Items {
		if nodeInstance.State != "started" {
			return false
		}
	}
	return true
}
