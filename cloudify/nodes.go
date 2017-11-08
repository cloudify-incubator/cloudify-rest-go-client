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

type NodePlugin struct {
	PluginBase
	Name     string `json:"name,omitempty"`
	Executor string `json:"executor,omitempty"`
	// TODO describe "install_arguments"
	// TODO describe "source"
	Install bool `json:"install"`
}

type Node struct {
	rest.CloudifyIDWithTenant
	Operations               map[string]interface{} `json:"operations,omitempty"`
	Relationships            []interface{}          `json:"relationships,omitempty"`
	DeployNumberOfInstances  int                    `json:"deploy_number_of_instances"`
	TypeHierarchy            []string               `json:"type_hierarchy,omitempty"`
	BlueprintID              string                 `json:"blueprint_id,omitempty"`
	NumberOfInstances        int                    `json:"number_of_instances"`
	DeploymentID             string                 `json:"deployment_id,omitempty"`
	Properties               map[string]interface{} `json:"properties,omitempty"`
	PlannedNumberOfInstances int                    `json:"planned_number_of_instances"`
	Plugins                  []NodePlugin           `json:"plugins,omitempty"`
	MaxNumberOfInstances     int                    `json:"max_number_of_instances"`
	HostID                   string                 `json:"host_id,omitempty"`
	MinNumberOfInstances     int                    `json:"min_number_of_instances"`
	Type                     string                 `json:"type,omitempty"`
	PluginsToInstall         []interface{}          `json:"plugins_to_install,omitempty"`
}

func (node *Node) GetJSONProperties() (string, error) {
	jsonData, err := json.Marshal(node.Properties)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type Nodes struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []Node                `json:"items"`
}

func (cl *Client) GetNodes(params map[string]string) (*Nodes, error) {
	var nodes Nodes

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("nodes?"+values.Encode(), &nodes)
	if err != nil {
		return nil, err
	}

	return &nodes, nil
}

func (cl *Client) GetStartedNodesWithType(params map[string]string, nodeType string) (*Nodes, error) {
	cloudNodes, err := cl.GetNodes(params)
	if err != nil {
		return nil, err
	}

	nodes := []Node{}
	for _, node := range cloudNodes.Items {

		notKubernetesHost := true
		for _, typeName := range node.TypeHierarchy {
			if typeName == nodeType {
				notKubernetesHost = false
				break
			}
		}

		if notKubernetesHost {
			continue
		}

		if node.NumberOfInstances <= 0 {
			continue
		}

		// add node to list
		nodes = append(nodes, node)
	}
	var result Nodes
	result.Items = nodes
	result.Metadata.Pagination.Total = uint(len(nodes))
	result.Metadata.Pagination.Size = uint(len(nodes))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
