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

// GetLoadBalancerInstances - return loadbalancer by name/namespace/cluster
func (cl *Client) GetLoadBalancerInstances(params map[string]string, clusterName, namespace, name, nodeType string) (*NodeInstances, error) {
	nodeInstancesList, err := cl.GetAliveNodeInstancesWithType(params, nodeType)
	if err != nil {
		return nil, err
	}

	instances := []NodeInstance{}
	for _, nodeInstance := range nodeInstancesList.Items {
		// Cluster
		if nodeInstance.GetStringProperty("proxy_cluster") != clusterName {
			continue
		}

		// Name
		if nodeInstance.GetStringProperty("proxy_name") != name {
			continue
		}

		// Namespace
		if nodeInstance.GetStringProperty("proxy_namespace") != namespace {
			continue
		}
		instances = append(instances, nodeInstance)
	}

	return cl.listNodeInstanceToNodeInstances(instances), nil
}
