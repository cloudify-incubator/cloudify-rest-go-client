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
			if cl.Debug {
				log.Printf("Can't load config: %s\n", err.Error())
			}
		} else {
			err = json.Unmarshal(configJSON, &agentConfig)
			if err != nil {
				if cl.Debug {
					log.Printf("Can't parse config: %s\n", err.Error())
				}
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

func (cl *Client) printf(format string, v ...interface{}) {
	if cl.Debug {
		log.Printf(format, v...)
	}
}
