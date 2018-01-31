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
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
)

// GetDeployment - return deployment by ID
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

// GetDeploymentInstancesHostGrouped - return instances grouped by host
func (cl *Client) GetDeploymentInstancesHostGrouped(params map[string]string) (map[string]NodeInstances, error) {
	var result = map[string]NodeInstances{}

	nodeInstances, err := cl.GetNodeInstances(params)
	if err != nil {
		return result, err
	}

	for _, nodeInstance := range nodeInstances.Items {
		if nodeInstance.HostID != "" {
			// add instance list if is not existed
			if _, ok := result[nodeInstance.HostID]; ok == false {
				result[nodeInstance.HostID] = NodeInstances{}
			}

			nodeHostInstance := result[nodeInstance.HostID]

			nodeHostInstance.Items = append(
				nodeHostInstance.Items, nodeInstance,
			)

			nodeHostInstance.Metadata.Pagination.Total++
			nodeHostInstance.Metadata.Pagination.Size++

			result[nodeInstance.HostID] = nodeHostInstance
		}
	}
	return result, nil
}

// GetDeploymentInstancesNodeGrouped - return instances grouped by node
func (cl *Client) GetDeploymentInstancesNodeGrouped(params map[string]string) (map[string]NodeInstances, error) {
	var result = map[string]NodeInstances{}

	nodeInstances, err := cl.GetNodeInstances(params)
	if err != nil {
		return result, err
	}

	for _, nodeInstance := range nodeInstances.Items {
		if nodeInstance.NodeID != "" {
			// add instance list if is not existed
			if _, ok := result[nodeInstance.NodeID]; ok == false {
				result[nodeInstance.NodeID] = NodeInstances{}
			}

			nodeHostInstance := result[nodeInstance.NodeID]

			nodeHostInstance.Items = append(
				nodeHostInstance.Items, nodeInstance,
			)

			nodeHostInstance.Metadata.Pagination.Total++
			nodeHostInstance.Metadata.Pagination.Size++

			result[nodeInstance.NodeID] = nodeHostInstance
		}
	}
	return result, nil
}

// GetNodeInstancesWithType - Returned list of started node instances with some node type,
// used mainly for kubernetes, also check that all instances related to same hostId started
func (cl *Client) GetNodeInstancesWithType(params map[string]string, nodeType string) (*NodeInstances, error) {
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

	instances := []NodeInstance{}
	for _, nodeInstance := range nodeInstances.Items {
		notKubernetesHost := true
		for _, node := range nodes.Items {
			if node.ID == nodeInstance.NodeID {
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

		// add instance to list
		instances = append(instances, nodeInstance)
	}
	var result NodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}

// GetAliveNodeInstancesWithType - Returned list of alive node instances with some node type,
// used mainly for kubernetes, need to get instances that can be joined to cluster
// Useful for cloudprovider logic only.
func (cl *Client) GetAliveNodeInstancesWithType(params map[string]string, nodeType string) (*NodeInstances, error) {
	nodeInstances, err := cl.GetNodeInstancesWithType(params, nodeType)
	if err != nil {
		return nil, err
	}

	// starting only because we restart kubelet after join
	aliveStates := []string{
		// "initializing", "creating", // workflow started for instance
		// "created", "configuring", // create action, had ip
		"configured", "starting", // configure action, joined to cluster
		"started", // everything done
	}
	instances := []NodeInstance{}
	for _, instance := range nodeInstances.Items {
		if utils.InList(aliveStates, instance.State) {
			instances = append(instances, instance)
		}
	}
	var result NodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}

// GetStartedNodeInstancesWithType - Returned list of started node instances with some node type,
// used mainly for kubernetes, also check that all instances related to same hostId started
// Useful for scale only.
func (cl *Client) GetStartedNodeInstancesWithType(params map[string]string, nodeType string) (*NodeInstances, error) {
	nodeInstancesGrouped, err := cl.GetDeploymentInstancesHostGrouped(params)
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

	instances := []NodeInstance{}
	for _, nodeInstances := range nodeInstancesGrouped {
		// check that all nodes on same hostID started
		allStarted := true
		for _, nodeInstance := range nodeInstances.Items {
			if nodeInstance.State != "started" {
				allStarted = false
				break
			}
		}

		if !allStarted {
			continue
		}

		// check type
		for _, nodeInstance := range nodeInstances.Items {
			notKubernetesHost := true
			for _, node := range nodes.Items {
				if node.ID == nodeInstance.NodeID {
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

			// add instance to list
			instances = append(instances, nodeInstance)
		}
	}
	var result NodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}

// GetDeploymentScaleGroup - return scaling group by name and deployment
func (cl *Client) GetDeploymentScaleGroup(deploymentID, scaleGroupName string) (*ScalingGroup, error) {
	deployment, err := cl.GetDeployment(deploymentID)
	if err != nil {
		return nil, err
	}
	if deployment.ScalingGroups != nil {
		for groupName, scaleGroup := range deployment.ScalingGroups {
			if scaleGroupName == groupName {
				return &scaleGroup, nil
			}
		}
	}
	return nil, fmt.Errorf("No such scale group:%+v", scaleGroupName)
}

// GetDeploymentScaleGroupNodes - return nodes related to scaling group
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

// GetDeploymentScaleGroupInstances - return instances related to scaling group
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

// GetDeploymentInstancesScaleGrouped - return instances grouped by scaleing group
func (cl *Client) GetDeploymentInstancesScaleGrouped(deploymentID, nodeType string) (map[string]NodeInstances, error) {
	var result = map[string]NodeInstances{}

	deployment, err := cl.GetDeployment(deploymentID)
	if err != nil {
		return result, err
	}

	var params = map[string]string{}
	params["deployment_id"] = deploymentID
	nodes, err := cl.GetStartedNodesWithType(params, nodeType)
	if err != nil {
		return result, err
	}

	cloudInstances, err := cl.GetStartedNodeInstancesWithType(params, nodeType)
	if err != nil {
		return result, err
	}

	if deployment.ScalingGroups != nil {
		// check what types we have in members
		for groupName, scaleGroup := range deployment.ScalingGroups {
			var resultedInstances = []NodeInstance{}
			var supportedMembers = []string{}
			for _, member := range scaleGroup.Members {
				supportedMembers = append(supportedMembers, member)
				for _, node := range nodes.Items {
					if node.HostID == member {
						if !utils.InList(supportedMembers, node.ID) {
							supportedMembers = append(supportedMembers, node.ID)
						}
					}
				}
			}

			// search instance
			for _, cloudInstance := range cloudInstances.Items {
				if utils.InList(supportedMembers, cloudInstance.NodeID) {
					resultedInstances = append(resultedInstances, cloudInstance)
				}
			}
			var resultInstance NodeInstances
			resultInstance.Items = resultedInstances
			resultInstance.Metadata.Pagination.Total = uint(len(resultedInstances))
			resultInstance.Metadata.Pagination.Size = uint(len(resultedInstances))
			resultInstance.Metadata.Pagination.Offset = 0
			result[groupName] = resultInstance
		}
	}
	return result, nil
}
