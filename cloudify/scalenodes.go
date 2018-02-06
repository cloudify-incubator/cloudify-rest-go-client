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
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
)

// NodeWithGroup - full information about cloudify node
type NodeWithGroup struct {
	Node
	ScalingGroupName string `json:"scaling_group"`
	GroupName        string `json:"group"`
}

// NodeWithGroups - response from manager with nodes list
type NodeWithGroups struct {
	rest.BaseMessage
	Metadata rest.Metadata   `json:"metadata"`
	Items    []NodeWithGroup `json:"items"`
}

// GetNodesFull - return nodes filtered by params
func (cl *Client) GetNodesFull(params map[string]string) (*NodeWithGroups, error) {
	var nodeWithGroups NodeWithGroups

	deploymentParams := map[string]string{}

	nodes, err := cl.GetNodes(params)
	if err != nil {
		return nil, err
	}

	if value, ok := params["deployment_id"]; ok == true {
		deploymentParams["deployment_id"] = value
	}

	deployments, err := cl.GetDeployments(deploymentParams)
	if err != nil {
		return nil, err
	}

	infoNodes := []NodeWithGroup{}
	for _, node := range nodes.Items {
		// copy original properties
		fullInfo := NodeWithGroup{}
		fullInfo.Node = node
		fullInfo.ScalingGroupName = ""
		fullInfo.GroupName = ""
		// update scaling group by deployments
		for _, deployment := range deployments.Items {
			if deployment.ID == node.DeploymentID {
				for scaleGroupName, scaleGroup := range deployment.ScalingGroups {
					if utils.InList(scaleGroup.Members, node.ID) {
						fullInfo.ScalingGroupName = scaleGroupName
					}
				}
				for groupName, group := range deployment.Groups {
					if utils.InList(group.Members, node.ID) {
						fullInfo.GroupName = groupName
					}
				}
			}
		}
		infoNodes = append(infoNodes, fullInfo)
	}

	// update group names from parent
	for childInd, child := range infoNodes {
		if child.HostID != child.ID {
			if child.GroupName == "" || child.ScalingGroupName == "" {
				for _, host := range infoNodes {
					if child.HostID == host.ID && child.DeploymentID == host.DeploymentID {
						if child.GroupName == "" {
							infoNodes[childInd].GroupName = host.GroupName
						}
						if child.ScalingGroupName == "" {
							infoNodes[childInd].ScalingGroupName = host.ScalingGroupName
						}
					}
				}
			}
		}
	}
	nodeWithGroups.Items = infoNodes
	nodeWithGroups.Metadata = nodes.Metadata

	return &nodeWithGroups, nil
}

// GetStartedNodesWithType - return nodes specified type with more than zero instances
func (cl *Client) GetStartedNodesWithType(params map[string]string, nodeType string) (*Nodes, error) {
	cloudNodes, err := cl.GetNodes(params)
	if err != nil {
		return nil, err
	}

	nodes := []Node{}
	for _, node := range cloudNodes.Items {

		if !utils.InList(node.TypeHierarchy, nodeType) {
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
