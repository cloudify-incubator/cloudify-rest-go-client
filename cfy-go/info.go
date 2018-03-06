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
Status

status - Show manager status [manager only].

	Manager state: Show service list on manager

		cfy-go status state

	Manager version: Show manager version

		cfy-go status version

	Kubernetes: Show diagnostic for current installation

		cfy-go status diag
*/
package main

import (
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

func servicesPrint(stat *cloudify.Status, err error) int {
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}
	fmt.Printf("Manager status: %v\n", stat.Status)
	fmt.Printf("Services:\n")
	lines := make([][]string, len(stat.Services))
	for pos, service := range stat.Services {
		lines[pos] = make([]string, 2)
		lines[pos][0] = service.DisplayName
		lines[pos][1] = service.Status()
	}
	utils.PrintTable([]string{"service", "status"}, lines)
	return 0
}

func versionPrint(ver *cloudify.Version, err error) int {
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	utils.PrintTable([]string{"Version", "Edition"},
		[][]string{{ver.Version, ver.Edition}})
	return 0
}

func instancesChecksPrint(nodeInstances *cloudify.NodeInstances, additional []string) {
	lines := make([][]string, len(nodeInstances.Items))
	for pos, nodeInstance := range nodeInstances.Items {
		lines[pos] = make([]string, 7+len(additional))
		lines[pos][0] = nodeInstance.ID
		lines[pos][1] = nodeInstance.DeploymentID
		lines[pos][2] = nodeInstance.HostID
		lines[pos][3] = nodeInstance.NodeID
		lines[pos][4] = nodeInstance.State
		lines[pos][5] = nodeInstance.GetStringProperty("hostname")

		for col, name := range additional {
			lines[pos][6+col] = nodeInstance.GetStringProperty(name)
		}

		lines[pos][6+len(additional)] = "looks good"
		if len(lines[pos][5]) >= 60 {
			lines[pos][6+len(additional)] = "Possible issues with nodes registration"
		}
	}
	headers := []string{"Id", "Deployment id", "Host id", "Node id",
		"State", "HostName"}

	headers = append(headers, additional...)
	headers = append(headers, "Note")
	utils.PrintTable(headers, lines)
}

func instancesChecks(cl *cloudify.Client, params map[string]string, typeName string, additionalProperties []string) int {
	nodeInstances, err := cl.GetStartedNodeInstancesWithType(params, typeName)
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	if len(nodeInstances.Items) == 0 {
		fmt.Printf("You don't have %v in current deployment.\n", typeName)
		return 0
	}

	instancesChecksPrint(nodeInstances, additionalProperties)
	return 0
}

func nodeWithGroup2CheckLine(node cloudify.NodeWithGroup) []string {
	var line = make([]string, 7)
	line[0] = node.ID
	line[1] = node.DeploymentID
	line[2] = node.HostID
	line[3] = node.Type
	line[4] = node.GroupName
	line[5] = node.ScalingGroupName
	line[6] = "looks good"
	if node.ScalingGroupName == "" || node.GroupName == "" {
		line[6] = "unscalable"
	}
	return line
}

func groupInstancesChecksPrint(nodes *cloudify.NodeWithGroups) {
	lines := [][]string{}
	for _, node := range nodes.Items {

		if node.Type != cloudify.KubernetesNode && node.Type != cloudify.KubernetesLoadBalancer {
			continue
		}

		lines = append(lines, nodeWithGroup2CheckLine(node))
	}
	utils.PrintTable([]string{
		"Id", "Deployment id", "Host id", "Type", "Group", "Scaling Group", "Notes",
	}, lines)
}

func groupInstancesChecks(cl *cloudify.Client, params map[string]string) int {
	nodes, err := cl.GetNodesFull(params)
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}

	if len(nodes.Items) == 0 {
		fmt.Println("You don't have nodes in current deployment.")
		return 0
	}

	groupInstancesChecksPrint(nodes)
	return 0
}

func runChecks(params map[string]string) int {
	cl := getClient()
	var res int

	fmt.Println("* Check manager services status.")
	fmt.Println("  Recheck by 'cfy-go status state'")
	res = servicesPrint(cl.GetStatus())
	if res != 0 {
		return res
	}

	fmt.Println("* Check properties in kubernetes instances.")
	fmt.Println("  Recheck by 'cfy-go node-instances started'")
	res = instancesChecks(cl, params, cloudify.KubernetesNode,
		[]string{"ip", "public_ip"})
	if res != 0 {
		return res
	}

	fmt.Println("* Check properties in kubernetes loadbalancers.")
	fmt.Println("  Recheck by 'cfy-go node-instances loadbalancer'")
	res = instancesChecks(cl, params, cloudify.KubernetesLoadBalancer,
		[]string{"ip", "public_ip", "Cluster", "Namespace", "Service"})
	if res != 0 {
		return res
	}

	fmt.Println("* Check scale group.")
	fmt.Println("  Recheck by 'cfy-go nodes group'")
	res = groupInstancesChecks(cl, params)
	if res != 0 {
		return res
	}

	return 0
}

func optionsToClient(command string, options []string) *cloudify.Client {
	operFlagSet := basicOptions(command)
	operFlagSet.Parse(options)
	cl := getClient()
	return cl
}

func infoOptions(args, options []string) int {
	defaultError := "state/version/diag subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "state":
		{
			cl := optionsToClient("status state", options)
			return servicesPrint(cl.GetStatus())
		}
	case "version":
		{
			cl := optionsToClient("status version", options)
			return versionPrint(cl.GetVersion())
		}
	case "diag":
		{
			var deployment string
			operFlagSet := basicOptions("status diag")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")

			operFlagSet.Parse(options)
			var params = map[string]string{}

			if deployment != "" {
				params["deployment_id"] = deployment
			}

			return runChecks(params)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
