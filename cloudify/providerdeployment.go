/*
Copyright (c) 2018 GigaSpaces Technologies Ltd. All rights reserved

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
	"os"
)

/*
Example for deployments info Which need to be used by cloudify provider
{
   "deployments":[
      {
         "id":"dep-1",
         "deployment_type":"node",
         "node_data_data_type":"cloudify.nodes.Kubernetes.Node"
      },
      {
         "id":"dep-2",
         "deployment_type":"load",
         "node_data_type":"cloudify.nodes.ApplicationServer.kubernetes.LoadBalancer"
      }
   ]
}
*/

// DeploymentsInfo - all deployments used on kubernetes cloudify provider
type DeploymentsInfo struct {
	Deployments []interface{} `json:"deployments,omitempty"`
}

//ParseDeploymentFile - Get deployments provider info needed to be used by kubernetes cloudify provider
func ParseDeploymentFile(deploymentFile string) (*DeploymentsInfo, error) {
	var deploymentInfo DeploymentsInfo

	raw, err := ioutil.ReadFile(deploymentFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &deploymentInfo)
	if err != nil {
		return nil, err
	}

	return &deploymentInfo, nil
}
