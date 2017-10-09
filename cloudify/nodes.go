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
	"log"
	"net/url"
)

type CloudifyNodePlugin struct {
	CloudifyPluginBase
	Name     string `json:"name,omitempty"`
	Executor string `json:"executor,omitempty"`
	// TODO describe "install_arguments"
	// TODO describe "source"
	Install bool `json:"install"`
}

type CloudifyNode struct {
	rest.CloudifyIdWithTenant
	Operations               map[string]interface{} `json:"operations,omitempty"`
	Relationships            []interface{}          `json:"relationships,omitempty"`
	DeployNumberOfInstances  int                    `json:"deploy_number_of_instances"`
	TypeHierarchy            []string               `json:"type_hierarchy,omitempty"`
	BlueprintId              string                 `json:"blueprint_id,omitempty"`
	NumberOfInstances        int                    `json:"number_of_instances"`
	DeploymentId             string                 `json:"deployment_id,omitempty"`
	Properties               map[string]interface{} `json:"properties,omitempty"`
	PlannedNumberOfInstances int                    `json:"planned_number_of_instances"`
	Plugins                  []CloudifyNodePlugin   `json:"plugins,omitempty"`
	MaxNumberOfInstances     int                    `json:"max_number_of_instances"`
	HostId                   string                 `json:"host_id,omitempty"`
	MinNumberOfInstances     int                    `json:"min_number_of_instances"`
	Type                     string                 `json:"type,omitempty"`
	PluginsToInstall         []interface{}          `json:"plugins_to_install,omitempty"`
}

func (node *CloudifyNode) GetJsonProperties() string {
	json_data, err := json.Marshal(node.Properties)
	if err != nil {
		log.Fatal(err)
	}
	return string(json_data)
}

type CloudifyNodes struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyNode        `json:"items"`
}

func (cl *CloudifyClient) GetNodes(params map[string]string) CloudifyNodes {
	var nodes CloudifyNodes

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("nodes?"+values.Encode(), &nodes)
	if err != nil {
		log.Fatal(err)
	}

	return nodes
}
