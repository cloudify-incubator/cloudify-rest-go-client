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

/*
Node Instances

node-instances - Handle a deployment's node-instances.

	get: Retrieve node-instance information [manager only]

		cfy-go node-instances list -deployment deployment

	list: List node-instances for a deployment [manager only]

		cfy-go node-instances list -deployment deployment

	started: check started instances in deployment (all, without filter by scaling group)

		cfy-go node-instances started -deployment <deployment_name>

	alive: check created instances in deployment (all, without filter by scaling group)

		cfy-go node-instances alive -deployment <deployment_name>

	loadbalancer: return list of loadbalancers, by default returned free nodes

		cfy-go node-instances alive -deployment <deployment_name>

	host-grouped: list instances grouped by hostID

		cfy-go node-instances host-grouped
*/
package main

import (
	"flag"
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"strings"
	"os"
)

func nodeInstancesPrint(nodeInstances *cloudify.NodeInstances, err error) int {
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	lines := make([][]string, len(nodeInstances.Items))
	for pos, nodeInstance := range nodeInstances.Items {
		var scaleGroups = []string{}
		if nodeInstance.ScalingGroups != nil {
			for _, scaleGroup := range nodeInstance.ScalingGroups {
				scaleGroups = append(scaleGroups, scaleGroup.Name)
			}
		}
		jsonString, err := nodeInstance.GetJSONRuntimeProperties()
		if err != nil {
			log.Printf("Cloudify error: %s\n", err.Error())
			jsonString = ""
		}

		var propertiesString = jsonString
		if len(jsonString) > 40 {
			propertiesString = jsonString[0:37] + "..."
		}

		lines[pos] = make([]string, 9)
		lines[pos][0] = nodeInstance.ID
		lines[pos][1] = nodeInstance.DeploymentID
		lines[pos][2] = nodeInstance.HostID
		lines[pos][3] = nodeInstance.NodeID
		lines[pos][4] = nodeInstance.State
		lines[pos][5] = nodeInstance.Tenant
		lines[pos][6] = nodeInstance.CreatedBy
		lines[pos][7] = strings.Join(scaleGroups, ", ")
		lines[pos][8] = propertiesString
	}
	utils.PrintTable([]string{
		"Id", "Deployment id", "Host id", "Node id", "State", "Tenant",
		"Created by", "Scaling Group", "Properties",
	}, lines)
	return 0
}

func parseInstancesFlags(operFlagSet *flag.FlagSet, options []string) map[string]string {
	var node string
	var deployment string
	var instance string
	var state string
	var hostID string

	operFlagSet.StringVar(&instance, "instance", "",
		"The unique identifier for the instance")
	operFlagSet.StringVar(&node, "node", "",
		"The unique identifier for the node")
	operFlagSet.StringVar(&deployment, "deployment", "",
		"The unique identifier for the deployment")
	operFlagSet.StringVar(&state, "state", "",
		"Filter by  state")
	operFlagSet.StringVar(&hostID, "host-id", "",
		"Filter by hostID")

	operFlagSet.Parse(options)

	params := parsePagination(operFlagSet, options)

	if instance != "" {
		params["id"] = instance
	}
	if node != "" {
		params["node_id"] = node
	}
	if deployment != "" {
		params["deployment_id"] = deployment
	}
	if state != "" {
		params["state"] = state
	}

	if hostID != "" {
		params["host_id"] = hostID
	}
	return params
}

func nodeInstancesOptions(args, options []string) int {
	defaultError := "list/started/host-grouped/node-grouped/by-type/alive/loadbalancer subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}
	switch args[2] {
	case "node-grouped":
		{
			operFlagSet := basicOptions("node-instances node-grouped")

			params := parseInstancesFlags(operFlagSet, options)

			cl := getClient()
			groupedInstances, err := cl.GetDeploymentInstancesNodeGrouped(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			for nodeID, instances := range groupedInstances {
				fmt.Printf("NodeID: %v\n", nodeID)
				if nodeInstancesPrint(&instances, nil) != 0 {
					return 1
				}
			}
			return 0
		}
	case "host-grouped":
		{
			operFlagSet := basicOptions("node-instances host-grouped")

			params := parseInstancesFlags(operFlagSet, options)

			cl := getClient()
			groupedInstances, err := cl.GetDeploymentInstancesHostGrouped(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			for hostID, instances := range groupedInstances {
				fmt.Printf("HostID: %v\n", hostID)
				if nodeInstancesPrint(&instances, nil) != 0 {
					return 1
				}
			}
			return 0
		}
	case "loadbalancer":
		{
			operFlagSet := basicOptions("node-instances loadbalancer")
			var nodeType string
			var loadbalancerName string
			var loadbalancerNamespace string
			var loadbalancerCluster string

			operFlagSet.StringVar(&loadbalancerName, "loadbalancer-name",
				"",
				"Filter by loadbalancer name")
			operFlagSet.StringVar(&loadbalancerNamespace, "loadbalancer-namespace",
				"",
				"Filter by loadbalancer namespace name")
			operFlagSet.StringVar(&loadbalancerCluster, "loadbalancer-cluster",
				"",
				"Filter by loadbalancer cluster name")

			params := parseInstancesFlags(operFlagSet, options)

			nodeType = os.Getenv("CFY_K8S_LOAD_TYPE")
			if nodeType == "" {
				nodeType = cloudify.KubernetesLoadBalancer
			}

			cl := getClient()
			return nodeInstancesPrint(
				cl.GetLoadBalancerInstances(
					params, loadbalancerCluster, loadbalancerNamespace,
					loadbalancerName, nodeType))
		}
	case "by-type":
		{
			operFlagSet := basicOptions("node-instances started")
			params := parseInstancesFlags(operFlagSet, options)

			var nodeType = os.Getenv("CFY_K8S_NODE_TYPE")
			if nodeType == "" {
				nodeType = cloudify.KubernetesNode
			}

			cl := getClient()
			return nodeInstancesPrint(
				cl.GetNodeInstancesWithType(params, nodeType))
		}
	case "started":
		{
			operFlagSet := basicOptions("node-instances started")
			params := parseInstancesFlags(operFlagSet, options)

			var nodeType = os.Getenv("CFY_K8S_NODE_TYPE")
			if nodeType == "" {
				nodeType = cloudify.KubernetesNode
			}

			cl := getClient()
			return nodeInstancesPrint(
				cl.GetStartedNodeInstancesWithType(params, nodeType))
		}
	case "alive":
		{
			operFlagSet := basicOptions("node-instances alive")
			params := parseInstancesFlags(operFlagSet, options)

			var nodeType = os.Getenv("CFY_K8S_NODE_TYPE")
			if nodeType == "" {
				nodeType = cloudify.KubernetesNode
			}

			cl := getClient()
			return nodeInstancesPrint(
				cl.GetAliveNodeInstancesWithType(params, nodeType))
		}
	case "list":
		{
			operFlagSet := basicOptions("node-instances list")

			params := parseInstancesFlags(operFlagSet, options)

			cl := getClient()
			nodeInstances, err := cl.GetNodeInstances(params)
			if nodeInstancesPrint(nodeInstances, err) != 0 {
				return 1
			}
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				nodeInstances.Metadata.Pagination.Offset, len(nodeInstances.Items),
				nodeInstances.Metadata.Pagination.Total)
			if len(nodeInstances.Items) == 1 {
				properties, err := nodeInstances.Items[0].GetJSONRuntimeProperties()
				if err != nil {
					log.Printf("Cloudify error: %s\n", err.Error())
					return 1
				}
				fmt.Printf("Runtime properties: %s\n", properties)
			} else {
				fmt.Printf("Limit to one row if you want to check RuntimeProperties\n")
			}
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
