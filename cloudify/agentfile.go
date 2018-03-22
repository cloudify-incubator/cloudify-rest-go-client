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
	"fmt"
	"io/ioutil"
	"log"
)

// CFYAgentConfig - useful(not all) fields for cloudify agent config
type CFYAgentConfig struct {
	RestHost string `json:"rest_host"`
	RestPort string `json:"rest_port"`
}

func (cl *Client) updateHostFromAgent() {
	if cl.AgentFile != "" {
		var agentConfig CFYAgentConfig
		if configJSON, err := ioutil.ReadFile(cl.AgentFile); err != nil {
			cl.debugLogf("Can't load config: %s\n", err.Error())
		} else {
			err = json.Unmarshal(configJSON, &agentConfig)
			if err != nil {
				cl.debugLogf("Can't parse config: %s\n", err.Error())
			} else {
				if agentConfig.RestPort != "" {
					cl.Host = "https://" + agentConfig.RestHost + ":" + agentConfig.RestPort
				} else {
					cl.Host = agentConfig.RestHost
				}
			}
		}
	}
}

func (cl *Client) debugLogf(format string, v ...interface{}) {
	if cl.Debug {
		log.Printf(format, v...)
	}
}

// ValidateBaseConnection - check configuration params (without tenant)
func ValidateBaseConnection(cloudConfig ClientConfig) error {
	if len(cloudConfig.Host) == 0 && len(cloudConfig.AgentFile) == 0 {
		return fmt.Errorf("You have empty host")
	}

	if len(cloudConfig.User) == 0 {
		return fmt.Errorf("You have empty user")
	}

	if len(cloudConfig.Password) == 0 {
		return fmt.Errorf("You have empty password")
	}

	return nil
}

// ValidateConnectionTenant - check configuration params
func ValidateConnectionTenant(cloudConfig ClientConfig) error {
	err := ValidateBaseConnection(cloudConfig)
	if err != nil {
		return nil
	}

	if len(cloudConfig.Tenant) == 0 {
		return fmt.Errorf("You have empty tenant")
	}

	if len(cloudConfig.DeploymentsFile) == 0 {
		return fmt.Errorf("You have empty deployments")
	}

	return nil
}
