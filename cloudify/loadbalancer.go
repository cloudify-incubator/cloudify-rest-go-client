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
	"reflect"
)

// GetLoadBalancerInstances - return loadbalancer by name/namespace/cluster
func (cl *Client) GetLoadBalancerInstances(params map[string]string, clusterName, namespace, name, nodeType string) (*NodeInstances, error) {
	nodeInstancesList, err := cl.GetAliveNodeInstancesWithType(params, nodeType)
	if err != nil {
		return nil, err
	}

	instances := []NodeInstance{}
	for _, nodeInstance := range nodeInstancesList.Items {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			// cluster
			if v, ok := nodeInstance.RuntimeProperties["proxy_cluster"]; ok == true {
				fmt.Printf("'%+v'(%+v) == '%+v'\n", v, reflect.TypeOf(v), clusterName)
				switch v.(type) {
				case string:
					{
						if v.(string) != clusterName {
							// node with different cluster
							continue
						}
					}
				}
			} else {
				// node without cluster
				if len(clusterName) != 0 {
					continue
				}
			}

			// name
			if v, ok := nodeInstance.RuntimeProperties["proxy_name"]; ok == true {
				fmt.Printf("'%+v'(%+v) == '%+v'\n", v, reflect.TypeOf(v), name)
				switch v.(type) {
				case string:
					{
						fmt.Printf("'%+v' == '%+v'", v.(string), name)
						if v.(string) != name {
							// node with different name
							continue
						}
					}
				}
			} else {
				// node without name
				if len(name) != 0 {
					continue
				}
			}

			// name space
			if v, ok := nodeInstance.RuntimeProperties["proxy_namespace"]; ok == true {
				fmt.Printf("'%+v'(%+v) == '%+v'\n", v, reflect.TypeOf(v), namespace)
				switch v.(type) {
				case string:
					{
						if v.(string) != namespace {
							// node with different name
							continue
						}
					}
				}
			} else {
				// node without namespace
				if len(namespace) != 0 {
					continue
				}
			}
			instances = append(instances, nodeInstance)
		} else if len(namespace) == 0 && len(name) == 0 && len(clusterName) == 0 {
			// special case for search first empty
			instances = append(instances, nodeInstance)
		}
	}
	var result NodeInstances
	result.Items = instances
	result.Metadata.Pagination.Total = uint(len(instances))
	result.Metadata.Pagination.Size = uint(len(instances))
	result.Metadata.Pagination.Offset = 0

	return &result, nil
}
