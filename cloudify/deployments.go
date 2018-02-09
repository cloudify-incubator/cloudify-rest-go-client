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
)

// Workflow - information about workflow
type Workflow struct {
	CreatedAt  string                 `json:"created_at"`
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// DeploymentPost - create deployment struct
type DeploymentPost struct {
	BlueprintID string                 `json:"blueprint_id"`
	Inputs      map[string]interface{} `json:"inputs"`
}

// SetJSONInputs - set inputs from json string
func (depl *DeploymentPost) SetJSONInputs(inputs string) error {
	if len(inputs) == 0 {
		depl.Inputs = map[string]interface{}{}
		return nil
	}

	return json.Unmarshal([]byte(inputs), &depl.Inputs)
}

// GetJSONInputs - get inputs as json string
func (depl *DeploymentPost) GetJSONInputs() (string, error) {
	jsonData, err := json.Marshal(depl.Inputs)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ScalingGroupProperties - scaling group properties struct
type ScalingGroupProperties struct {
	MinInstances     int `json:"min_instances"`
	PlannedInstances int `json:"planned_instances"`
	DefaultInstances int `json:"default_instances"`
	MaxInstances     int `json:"max_instances"`
	CurrentInstances int `json:"current_instances"`
}

// ScalingGroup - Scaling group struct
type ScalingGroup struct {
	Properties ScalingGroupProperties `json:"properties"`
	Members    []string               `json:"members"`
}

// NodeGroup - Node group struct
type NodeGroup struct {
	Members []string `json:"members"`
	// TODO use correct type for "policies" struct
	Policies map[string]interface{} `json:"policies"`
}

// Deployment - deployment struct
type Deployment struct {
	// have id, owner information
	rest.Resource
	// contain information from post
	DeploymentPost
	Permalink     string                  `json:"permalink"`
	Workflows     []Workflow              `json:"workflows"`
	Outputs       map[string]interface{}  `json:"outputs"`
	ScalingGroups map[string]ScalingGroup `json:"scaling_groups"`
	Groups        map[string]NodeGroup    `json:"groups"`
	// TODO use correct type for "policy_types" struct
	PolicyTypes map[string]interface{} `json:"policy_types"`
	// TODO use correct type for "policy_triggers" struct
	PolicyTriggers map[string]interface{} `json:"policy_triggers"`
}

// GetJSONOutputs - get deployments outputs as json string
func (depl *Deployment) GetJSONOutputs() (string, error) {
	jsonData, err := json.Marshal(depl.Outputs)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GetJSONInputs - get deployments inputs as json string
func (depl *Deployment) GetJSONInputs() (string, error) {
	jsonData, err := json.Marshal(depl.Inputs)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// DeploymentGet - information about deployment on server
type DeploymentGet struct {
	// can be response from api
	rest.BaseMessage
	Deployment
}

// Deployments - response with list deployments
type Deployments struct {
	rest.BaseMessage
	Metadata rest.Metadata `json:"metadata"`
	Items    []Deployment  `json:"items"`
}

// GetDeployments - get deployments list from server filtered by params
func (cl *Client) GetDeployments(params map[string]string) (*Deployments, error) {
	var deployments Deployments

	values := cl.stringMapToURLValue(params)

	err := cl.Get("deployments?"+values.Encode(), &deployments)
	if err != nil {
		return nil, err
	}

	return &deployments, nil
}

// DeleteDeployments - delete deployment by ID
func (cl *Client) DeleteDeployments(deploymentID string) (*DeploymentGet, error) {
	var deployment DeploymentGet

	err := cl.Delete("deployments/"+deploymentID, &deployment)
	if err != nil {
		return nil, err
	}

	return &deployment, nil
}

// CreateDeployments - create deployment
func (cl *Client) CreateDeployments(deploymentID string, depl DeploymentPost) (*DeploymentGet, error) {
	var deployment DeploymentGet

	err := cl.Put("deployments/"+deploymentID, depl, &deployment)
	if err != nil {
		return nil, err
	}

	return &deployment, nil
}
