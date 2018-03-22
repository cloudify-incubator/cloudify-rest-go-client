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
	"io"
	"log"
	"os"
)

// ServiceConfig - settings for connect to cloudify
type ServiceConfig struct {
	ClientConfig
	// TODO Add support for dynamic type from proxy deployment (load + node)
	//NodeType string `json:"nodeType,omitempty"`
}

// ServiceClientInit - common functionality for load config for service
func ServiceClientInit(config io.Reader) (*ServiceConfig, error) {
	var cloudConfig ServiceConfig
	cloudConfig.Host = os.Getenv("CFY_HOST")
	cloudConfig.User = os.Getenv("CFY_USER")
	cloudConfig.Password = os.Getenv("CFY_PASSWORD")
	cloudConfig.Tenant = os.Getenv("CFY_TENANT")
	cloudConfig.AgentFile = os.Getenv("CFY_AGENT")
	cloudConfig.DeploymentsFile = os.Getenv("CFY_DEPLOYMENTS")

	// TODO Add support
	//cloudConfig.NodeType = os.Getenv("CFY_NODE_TYPE")

	if config != nil {
		err := json.NewDecoder(config).Decode(&cloudConfig)
		if err != nil {
			return nil, err
		}
	}

	configErr := ValidateConnectionTenant(cloudConfig.ClientConfig)
	if configErr != nil {
		return nil, configErr
	}

	log.Printf("Config %+v", cloudConfig)

	return &cloudConfig, nil
}
