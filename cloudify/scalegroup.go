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
)

func (cl *Client) GetDeployment(deploymentID string) (*Deployment, error) {
	var params = map[string]string{}
	params["id"] = deploymentID
	deployments, err := cl.GetDeployments(params)
	if err != nil {
		return nil, err
	}
	if len(deployments.Items) != 1 {
		return nil, fmt.Errorf("Returned wrong count of deployments:%+v", deploymentID)
	}
	return &deployments.Items[0], nil
}

func (cl *Client) GetDeploymentScaleGroup(deploymentID, groupName string) (*ScalingGroup, error) {
	deployment, err := cl.GetDeployment(deploymentID)
	if err != nil {
		return nil, err
	}
	if deployment.ScalingGroups != nil {
		for groupName, scaleGroup := range deployment.ScalingGroups {
			if groupName == groupName {
				return &scaleGroup, nil
			}
		}
	}
	return nil, fmt.Errorf("No such scale group:%+v", groupName)
}

func (cl *Client) GetDeploymentScaleGroupNodes(deploymentID, groupName, nodeType string) (*Nodes, error) {
	// get all nodes
	params := map[string]string{}
	params["deployment_id"] = deploymentID
	cloudNodes, err := cl.GetStartedNodesWithType(params, nodeType)
	if err != nil {
		return nil, err
	}

	// get scale group
	scaleGroup, err := cl.GetDeploymentScaleGroup(deploymentID, groupName)
	if err != nil {
		return nil, err
	}

	// filter by scaling group
	nodes := []Node{}
	for _, node := range cloudNodes.Items {
		for _, nodeID := range scaleGroup.Members {
			if nodeID == node.ID || nodeID == node.HostID {
				nodes = append(nodes, node)
			}
		}
	}
	var result Nodes
	result.Items = nodes
	result.Metadata.Pagination.Total = uint(len(nodes))
	result.Metadata.Pagination.Size = uint(len(nodes))
	result.Metadata.Pagination.Offset = 0
	return &result, nil
}

func (cl *Client) GetDeploymentScaleGroupInstances(deploymentID, groupName, nodeType string) (*NodeInstances, error) {
	// get all instances
	params := map[string]string{}
	params["deployment_id"] = deploymentID
	cloudInstances, err := cl.GetStartedNodeInstancesWithType(params, nodeType)
	if err != nil {
		return nil, err
	}

	// get nodes in scale group (need to get nodes because we need host for each)
	cloudNodes, err := cl.GetDeploymentScaleGroupNodes(deploymentID, groupName, nodeType)
	if err != nil {
		return nil, err
	}

	// filter by scaling group
	instances := []NodeInstance{}
	for _, instance := range cloudInstances.Items {
		for _, node := range cloudNodes.Items {
			if node.ID == instance.NodeID {
				instances = append(instances, instance)
			}
		}
	}
	var result NodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
