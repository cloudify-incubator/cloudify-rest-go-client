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
