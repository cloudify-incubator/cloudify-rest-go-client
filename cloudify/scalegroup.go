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

func (cl *CloudifyClient) GetDeployment(deploymentID string) (*CloudifyDeployment, error) {
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

func (cl *CloudifyClient) GetDeploymentScaleGroup(deploymentID, groupName string) (*ScalingGroup, error) {
	deployment, err := cl.GetDeployment(deploymentID)
	if err != nil {
		return nil, err
	}
	if deployment.ScalingGroups != nil {
		for group_name, scale_group := range deployment.ScalingGroups {
			if group_name == groupName {
				return &scale_group, nil
			}
		}
	}
	return nil, fmt.Errorf("No such scale group:%+v", groupName)
}

func (cl *CloudifyClient) GetDeploymentScaleGroupNodes(deploymentID, groupName, node_type string) (*CloudifyNodes, error) {
	// get all nodes
	params := map[string]string{}
	params["deployment_id"] = deploymentID
	cloud_nodes, err := cl.GetStartedNodesWithType(params, node_type)
	if err != nil {
		return nil, err
	}

	// get scale group
	scale_group, err := cl.GetDeploymentScaleGroup(deploymentID, groupName)
	if err != nil {
		return nil, err
	}

	// filter by scaling group
	nodes := []CloudifyNode{}
	for _, node := range cloud_nodes.Items {
		for _, nodeId := range scale_group.Members {
			if nodeId == node.Id || nodeId == node.HostId {
				nodes = append(nodes, node)
			}
		}
	}
	var result CloudifyNodes
	result.Items = nodes
	result.Metadata.Pagination.Total = uint(len(nodes))
	result.Metadata.Pagination.Size = uint(len(nodes))
	result.Metadata.Pagination.Offset = 0
	return &result, nil
}

func (cl *CloudifyClient) GetDeploymentScaleGroupInstances(deploymentID, groupName, node_type string) (*CloudifyNodeInstances, error) {
	// get all instances
	params := map[string]string{}
	params["deployment_id"] = deploymentID
	cloud_instances, err := cl.GetStartedNodeInstancesWithType(params, node_type)
	if err != nil {
		return nil, err
	}

	// get nodes in scale group (need to get nodes because we need host for each)
	cloud_nodes, err := cl.GetDeploymentScaleGroupNodes(deploymentID, groupName, node_type)
	if err != nil {
		return nil, err
	}

	// filter by scaling group
	instances := []CloudifyNodeInstance{}
	for _, instance := range cloud_instances.Items {
		for _, node := range cloud_nodes.Items {
			if node.Id == instance.NodeId {
				instances = append(instances, instance)
			}
		}
	}
	var result CloudifyNodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
