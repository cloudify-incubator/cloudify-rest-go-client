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
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
)

// NodePlugin - information about plugin used by node
type NodePlugin struct {
	PluginBase
	Name     string `json:"name,omitempty"`
	Executor string `json:"executor,omitempty"`
	// TODO describe "install_arguments"
	// TODO describe "source"
	Install bool `json:"install"`
}

// Node - information about cloudify node
type Node struct {
	rest.ObjectIDWithTenant
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

// GetJSONProperties - properties related to node
func (node *Node) GetJSONProperties() (string, error) {
	jsonData, err := json.Marshal(node.Properties)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Nodes - response from manager with nodes list
type Nodes struct {
	rest.BaseMessage
	Metadata rest.Metadata `json:"metadata"`
	Items    []Node        `json:"items"`
}

func (nl *Nodes) GetNodeNamesWithType(nodeType string) []string {
	nodeIDS := []string{}

	for _, node := range nl.Items {
		if utils.InList(node.TypeHierarchy, nodeType) {
			nodeIDS = append(nodeIDS, node.ID)
		}
	}
	return nodeIDS
}

// GetNodes - return nodes filtered by params
func (cl *Client) GetNodes(params map[string]string) (*Nodes, error) {
	var nodes Nodes

	values := cl.stringMapToURLValue(params)

	err := cl.Get("nodes?"+values.Encode(), &nodes)
	if err != nil {
		return nil, err
	}

	return &nodes, nil
}
