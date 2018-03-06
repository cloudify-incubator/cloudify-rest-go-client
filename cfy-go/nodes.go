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
Nodes

nodes - Handle a deployment's nodes

	get: Retrieve node information [manager only]

		cfy-go nodes list -node server -deployment deployment


	list: List nodes for a deployment [manager only]

		cfy-go nodes list

	group: List nodes for a deployment [manager only], with groups names

		cfy-go nodes group

	started - check started nodes in deployment (all, without filter by scaling group)

		cfy-go nodes started -deployment deployment

*/
package main

import (
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

func nodesGroupPrint(nodes *cloudify.NodeWithGroups, err error) int {
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	lines := make([][]string, len(nodes.Items))
	for pos, node := range nodes.Items {
		lines[pos] = make([]string, 6)
		lines[pos][0] = node.ID
		lines[pos][1] = node.DeploymentID
		lines[pos][2] = node.HostID
		lines[pos][3] = node.Type
		lines[pos][4] = node.GroupName
		lines[pos][5] = node.ScalingGroupName
	}
	utils.PrintTable([]string{
		"Id", "Deployment id", "Host id", "Type", "Group", "Scaling Group",
	}, lines)
	fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
		nodes.Metadata.Pagination.Offset, len(nodes.Items),
		nodes.Metadata.Pagination.Total)
	return 0
}

func nodesPrint(nodes *cloudify.Nodes, err error) int {
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	lines := make([][]string, len(nodes.Items))
	for pos, node := range nodes.Items {
		lines[pos] = make([]string, 9)
		lines[pos][0] = node.ID
		lines[pos][1] = node.DeploymentID
		lines[pos][2] = node.BlueprintID
		lines[pos][3] = node.HostID
		lines[pos][4] = node.Type
		lines[pos][5] = fmt.Sprintf("%d", node.NumberOfInstances)
		lines[pos][6] = fmt.Sprintf("%d", node.PlannedNumberOfInstances)
		lines[pos][7] = node.Tenant
		lines[pos][8] = node.CreatedBy
	}
	utils.PrintTable([]string{
		"Id", "Deployment id", "Blueprint id", "Host id", "Type",
		"Number of instances", "Planned number of instances",
		"Tenant", "created_by",
	}, lines)
	fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
		nodes.Metadata.Pagination.Offset, len(nodes.Items),
		nodes.Metadata.Pagination.Total)
	if len(nodes.Items) == 1 {
		properties, err := nodes.Items[0].GetJSONProperties()
		if err != nil {
			log.Printf("Cloudify error: %s\n", err.Error())
			return 1
		}
		fmt.Printf("Properties: %s\n", properties)
	} else {
		fmt.Printf("Limit to one row if you want to check Properties\n")
	}
	return 0
}

func nodesOptions(args, options []string) int {
	defaultError := "list/group/started subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "started":
		{
			operFlagSet := basicOptions("nodes started")
			var node string
			var deployment string
			var nodeType string
			var hostID string
			operFlagSet.StringVar(&node, "node", "",
				"The unique identifier for the node")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&nodeType, "node-type",
				cloudify.KubernetesNode, "Filter by node type")
			operFlagSet.StringVar(&hostID, "host-id", "",
				"Filter by hostID")

			operFlagSet.Parse(options)

			var params = map[string]string{}

			if node != "" {
				params["id"] = node
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}
			if hostID != "" {
				params["host_id"] = hostID
			}

			cl := getClient()
			return nodesPrint(cl.GetStartedNodesWithType(params, nodeType))
		}
	case "group":
		{
			operFlagSet := basicOptions("nodes group")
			var node string
			var deployment string
			var nodeType string
			operFlagSet.StringVar(&node, "node", "",
				"The unique identifier for the node")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&nodeType, "node-type",
				"",
				"Filter by node type")

			params := parsePagination(operFlagSet, options)

			if node != "" {
				params["id"] = node
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}
			if nodeType != "" {
				params["type"] = nodeType
			}

			cl := getClient()
			return nodesGroupPrint(cl.GetNodesFull(params))
		}
	case "list":
		{
			operFlagSet := basicOptions("nodes list")
			var node string
			var deployment string
			var nodeType string
			operFlagSet.StringVar(&node, "node", "",
				"The unique identifier for the node")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&nodeType, "node-type",
				"",
				"Filter by node type")

			params := parsePagination(operFlagSet, options)

			if node != "" {
				params["id"] = node
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}
			if nodeType != "" {
				params["type"] = nodeType
			}

			cl := getClient()
			return nodesPrint(cl.GetNodes(params))
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
