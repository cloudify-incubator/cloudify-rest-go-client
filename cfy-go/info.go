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

	Kubernetes: Show diagnostic for current installation [deployment-id is optional]
		Show diagnostic for all current installation
			cfy-go status diag [-deployment deployment-id]
			cfy-go status diag -all [-deployment deployment-id]

		Show diagnostic only for Kubernetes nodes
			cfy-go status diag -node [-deployment deployment-id]

		Show diagnostic only for Kubernetes load balancer
			cfy-go status diag -load [-deployment deployment-id]


*/
package main

import (
	"flag"
	"fmt"
	"github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	"github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"os"
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

func groupInstancesChecksPrint(nodes *cloudify.NodeWithGroups) {
	lines := [][]string{}
	for _, node := range nodes.Items {

		if node.Type != cloudify.KubernetesNode && node.Type != cloudify.KubernetesLoadBalancer {
			continue
		}

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

		lines = append(lines, line)
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
	res := runChecksStatus() != 1 && runChecksNode(params) != 1 && runChecksLoad(params) != 1 && runChecksScaleGroup(params) != 1
	if res {
		return 0
	}
	return 1
}

func runChecksNode(params map[string]string) int {
	cl := getClient()
	var res int
	var nodeType = os.Getenv("CFY_K8S_NODE_TYPE")
	if nodeType == "" {
		nodeType = cloudify.KubernetesNode
	}

	fmt.Println("* Check properties in kubernetes instances.")
	fmt.Println("  Recheck by 'cfy-go node-instances started'")
	res = instancesChecks(cl, params, nodeType,
		[]string{"ip", "public_ip"})
	if res != 0 {
		return res
	}
	return 0
}

func runChecksLoad(params map[string]string) int {
	cl := getClient()
	var res int
	var loadType = os.Getenv("CFY_K8S_LOAD_TYPE")
	if loadType == "" {
		loadType = cloudify.KubernetesLoadBalancer
	}

	fmt.Println("* Check properties in kubernetes loadbalancers.")
	fmt.Println("  Recheck by 'cfy-go node-instances loadbalancer'")
	res = instancesChecks(cl, params, loadType,
		[]string{"ip", "public_ip", "proxy_cluster", "proxy_namespace", "proxy_name"})
	if res != 0 {
		return res
	}
	return 0
}

func runChecksStatus() int {
	cl := getClient()
	var res int

	fmt.Println("* Check manager services status.")
	fmt.Println("  Recheck by 'cfy-go status state'")
	res = servicesPrint(cl.GetStatus())
	if res != 0 {
		return res
	}
	return 0
}

func runChecksScaleGroup(params map[string]string) int {
	cl := getClient()
	var res int

	fmt.Println("* Check scale group.")
	fmt.Println("  Recheck by 'cfy-go nodes group'")
	res = groupInstancesChecks(cl, params)
	if res != 0 {
		return res
	}
	return 0
}

func optionsToClient(operFlagSet *flag.FlagSet, options []string) *cloudify.Client {
	operFlagSet.Parse(options)
	cl := getClient()
	return cl
}

func stateInfoCall(operFlagSet *flag.FlagSet, args, options []string) int {
	cl := optionsToClient(operFlagSet, options)
	return servicesPrint(cl.GetStatus())
}

func versionInfoCall(operFlagSet *flag.FlagSet, args, options []string) int {
	cl := optionsToClient(operFlagSet, options)
	return versionPrint(cl.GetVersion())
}

func diagInfoCall(operFlagSet *flag.FlagSet, args, options []string) int {
	var deployment string
	var diagAll bool
	var diagNode bool
	var diagLoad bool

	operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
	operFlagSet.BoolVar(&diagAll, "all", false, "Flag to check if need to diagnose all nodes (node + load) types")
	operFlagSet.BoolVar(&diagNode, "node", false, "Flag to check if need to diagnose only nodes types")
	operFlagSet.BoolVar(&diagLoad, "load", false, "Flag to check if need to diagnose only load node types")

	operFlagSet.Parse(options)
	var params = map[string]string{}

	if deployment != "" {
		params["deployment_id"] = deployment
	}

	if diagAll {
		return runChecks(params)
	} else if diagNode {
		return runChecksNode(params)
	} else if diagLoad {
		return runChecksLoad(params)
	}
	return runChecks(params)
}

func infoOptions(args, options []string) int {
	var pluginsCalls = []CommandInfo{{
		CommandName: "state",
		Callback:    stateInfoCall,
	}, {
		CommandName: "version",
		Callback:    versionInfoCall,
	}, {
		CommandName: "diag",
		Callback:    diagInfoCall,
	}}

	return ParseCalls(pluginsCalls, 3, args, options)
}
