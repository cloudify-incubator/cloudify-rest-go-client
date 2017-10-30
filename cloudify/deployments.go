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

// Check https://blog.golang.org/json-and-go for more info about json marshaling.
type CloudifyWorkflow struct {
	CreatedAt  string                 `json:"created_at"`
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type CloudifyDeploymentPost struct {
	BlueprintId string                 `json:"blueprint_id"`
	Inputs      map[string]interface{} `json:"inputs"`
}

func (depl *CloudifyDeploymentPost) SetJsonInputs(inputs string) error {
	if len(inputs) == 0 {
		depl.Inputs = map[string]interface{}{}
		return nil
	}

	return json.Unmarshal([]byte(inputs), &depl.Inputs)
}

func (depl *CloudifyDeploymentPost) GetJsonInputs() (string, error) {
	json_data, err := json.Marshal(depl.Inputs)
	if err != nil {
		return "", err
	}
	return string(json_data), nil
}

type ScalingGroupProperties struct {
	MinInstances     int `json:"min_instances"`
	PlannedInstances int `json:"planned_instances"`
	DefaultInstances int `json:"default_instances"`
	MaxInstances     int `json:"max_instances"`
	CurrentInstances int `json:"current_instances"`
}

type ScalingGroup struct {
	Properties ScalingGroupProperties `json:"properties"`
	Members    []string               `json:"members"`
}
type CloudifyDeployment struct {
	// have id, owner information
	rest.CloudifyResource
	// contain information from post
	CloudifyDeploymentPost
	Permalink     string                  `json:"permalink"`
	Workflows     []CloudifyWorkflow      `json:"workflows"`
	Outputs       map[string]interface{}  `json:"outputs"`
	ScalingGroups map[string]ScalingGroup `json:"scaling_groups"`
	// TODO describe "policy_types" struct
	// TODO describe "policy_triggers" struct
	// TODO describe "groups" struct
	// TODO describe "scaling_groups" struct
}

func (depl *CloudifyDeployment) GetJsonOutputs() (string, error) {
	json_data, err := json.Marshal(depl.Outputs)
	if err != nil {
		return "", err
	}
	return string(json_data), nil
}

func (depl *CloudifyDeployment) GetJsonInputs() (string, error) {
	json_data, err := json.Marshal(depl.Inputs)
	if err != nil {
		return "", err
	}
	return string(json_data), nil
}

type CloudifyDeploymentGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	CloudifyDeployment
}

type CloudifyDeployments struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyDeployment  `json:"items"`
}

func (cl *CloudifyClient) GetDeployments(params map[string]string) (*CloudifyDeployments, error) {
	var deployments CloudifyDeployments

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("deployments?"+values.Encode(), &deployments)
	if err != nil {
		return nil, err
	}

	return &deployments, nil
}

func (cl *CloudifyClient) DeleteDeployments(deployment_id string) (*CloudifyDeploymentGet, error) {
	var deployment CloudifyDeploymentGet

	err := cl.Delete("deployments/"+deployment_id, &deployment)
	if err != nil {
		return nil, err
	}

	return &deployment, nil
}

func (cl *CloudifyClient) CreateDeployments(deployment_id string, depl CloudifyDeploymentPost) (*CloudifyDeploymentGet, error) {
	var deployment CloudifyDeploymentGet

	err := cl.Put("deployments/"+deployment_id, depl, &deployment)
	if err != nil {
		return nil, err
	}

	return &deployment, nil
}
