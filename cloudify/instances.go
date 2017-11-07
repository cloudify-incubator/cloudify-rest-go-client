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
	"net/url"
)

type CloudifyNodeInstance struct {
	rest.CloudifyIdWithTenant
	Relationships     []interface{}          `json:"relationships,omitempty"`
	RuntimeProperties map[string]interface{} `json:"runtime_properties,omitempty"`
	State             string                 `json:"state,omitempty"`
	Version           int                    `json:"version,omitempty"`
	HostId            string                 `json:"host_id,omitempty"`
	DeploymentId      string                 `json:"deployment_id,omitempty"`
	NodeId            string                 `json:"node_id,omitempty"`
	// TODO describe "scaling_groups" struct
}

func (instance *CloudifyNodeInstance) GetJsonRuntimeProperties() (string, error) {
	jsonData, err := json.Marshal(instance.RuntimeProperties)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type CloudifyNodeInstances struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata  `json:"metadata"`
	Items    []CloudifyNodeInstance `json:"items"`
}

/* Get all node instances */
func (cl *CloudifyClient) GetNodeInstances(params map[string]string) (*CloudifyNodeInstances, error) {
	var instances CloudifyNodeInstances

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("node-instances?"+values.Encode(), &instances)
	if err != nil {
		return nil, err
	}

	return &instances, nil
}

/* Returned list of started node instances with some node type,
 * used mainly for kubernetes */
func (cl *CloudifyClient) GetStartedNodeInstancesWithType(params map[string]string, nodeType string) (*CloudifyNodeInstances, error) {
	nodeInstances, err := cl.GetNodeInstances(params)
	if err != nil {
		return nil, err
	}

	var nodeParams = map[string]string{}
	if val, ok := params["deployment_id"]; ok {
		nodeParams["deployment_id"] = val
	}
	nodes, err := cl.GetNodes(nodeParams)
	if err != nil {
		return nil, err
	}

	instances := []CloudifyNodeInstance{}
	for _, nodeInstance := range nodeInstances.Items {
		var notKubernetesHost bool = true
		for _, node := range nodes.Items {
			if node.Id == nodeInstance.NodeId {
				for _, typeName := range node.TypeHierarchy {
					if typeName == nodeType {
						notKubernetesHost = false
						break
					}
				}
			}
		}

		if notKubernetesHost {
			continue
		}

		if nodeInstance.State != "started" {
			continue
		}

		// add instance to list
		instances = append(instances, nodeInstance)
	}
	var result CloudifyNodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
