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

// NodeWithGroup - full information about cloudify node
type NodeWithGroup struct {
	ID               string `json:"id"`
	DeploymentID     string `json:"deployment_id,omitempty"`
	Type             string `json:"type,omitempty"`
	HostID           string `json:"host_id,omitempty"`
	ScalingGroupName string `json:"scaling_group"`
	GroupName        string `json:"group"`
}

// NodeWithGroups - response from manager with nodes list
type NodeWithGroups struct {
	rest.BaseMessage
	Metadata rest.Metadata   `json:"metadata"`
	Items    []NodeWithGroup `json:"items"`
}

// GetNodes - return nodes filtered by params
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

// GetNodesFull - return nodes filtered by params
func (cl *Client) GetNodesFull(params map[string]string) (*NodeWithGroups, error) {
	var nodes Nodes
	var NodeWithGroups NodeWithGroups

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("nodes?"+values.Encode(), &nodes)
	if err != nil {
		return nil, err
	}

	infoNodes := []NodeWithGroup{}
	for _, node := range nodes.Items {
		fullInfo := NodeWithGroup{}
		fullInfo.ID = node.ID
		fullInfo.DeploymentID = node.DeploymentID
		fullInfo.Type = node.Type
		fullInfo.HostID = node.HostID
		fullInfo.ScalingGroupName = ""
		fullInfo.GroupName = ""
		infoNodes = append(infoNodes, fullInfo)
	}

	NodeWithGroups.Items = infoNodes
	NodeWithGroups.Metadata = nodes.Metadata

	return &NodeWithGroups, nil
}

// GetStartedNodesWithType - return nodes specified type with more than zero instances
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
