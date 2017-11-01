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

func (instance *CloudifyNodeInstance) GetJsonRuntimeProperties() (string, error) {
	json_data, err := json.Marshal(instance.RuntimeProperties)
	if err != nil {
		return "", err
	}
	return string(json_data), nil
}

type CloudifyNodeInstances struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata  `json:"metadata"`
	Items    []CloudifyNodeInstance `json:"items"`
}

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

func (cl *CloudifyClient) GetStartedNodeInstances(params map[string]string, node_type string) (*CloudifyNodeInstances, error) {
	nodeInstances, err := cl.GetNodeInstances(params)
	if err != nil {
		return nil, err
	}

	instances := []CloudifyNodeInstance{}
	for _, nodeInstance := range nodeInstances.Items {
		var node_params = map[string]string{}
		node_params["id"] = nodeInstance.NodeId
		nodes, err := cl.GetNodes(node_params)
		if err != nil {
			if cl.restCl.Debug {
				log.Printf("Not found instances: %+v", err)
			}
			continue
		}
		if len(nodes.Items) != 1 {
			if cl.restCl.Debug {
				log.Printf("Found more than one node by nodeId: %+v", nodeInstance.NodeId)
			}
			continue
		}

		var not_kubernetes_host bool = true
		for _, type_name := range nodes.Items[0].TypeHierarchy {
			if type_name == node_type {
				not_kubernetes_host = false
				break
			}
		}

		if not_kubernetes_host {
			continue
		}

		if nodeInstance.State != "started" {
			continue
		}

		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			instances = append(instances, nodeInstance)
		}
	}
	var result CloudifyNodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
