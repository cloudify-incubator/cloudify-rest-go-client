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
scaling-groups - operations related to Scaling Groups

	groups: check nodes in group - recheck code used in get scaling group by instance(hostname) in autoscale

		cfy-go scaling-groups groups -deployment <deployment_name>

	nodes: check nodes in group in autoscale, check that we have node in scaling group

		cfy-go scaling-groups nodes -deployment <deployment_name> -scalegroup <scale_group_name>

	instances: check instances in group in autoscale

		cfy-go scaling-groups instances -deployment <deployment_name> -scalegroup <scale_group_name>
*/
package main

import (
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"strings"
)

func scaleGroupPrint(deploymentScalingGroups map[string]cloudify.ScalingGroup) int {
	lines := make([][]string, len(deploymentScalingGroups))
	var pos int
	if deploymentScalingGroups != nil {
		for groupName, scaleGroup := range deploymentScalingGroups {
			lines[pos] = make([]string, 7)
			lines[pos][0] = groupName
			lines[pos][1] = strings.Join(scaleGroup.Members, ", ")
			lines[pos][2] = fmt.Sprintf("%d", scaleGroup.Properties.MinInstances)
			lines[pos][3] = fmt.Sprintf("%d", scaleGroup.Properties.PlannedInstances)
			lines[pos][4] = fmt.Sprintf("%d", scaleGroup.Properties.DefaultInstances)
			lines[pos][5] = fmt.Sprintf("%d", scaleGroup.Properties.MaxInstances)
			lines[pos][6] = fmt.Sprintf("%d", scaleGroup.Properties.CurrentInstances)
			pos++
		}
	}
	utils.PrintTable([]string{
		"Group name", "Members", "Min Instances", "Planned Instances",
		"Default Instances", "Max Instances", "Current Instances",
	}, lines)
	return 0
}

func scalingGroupsOptions(args, options []string) int {
	defaultError := "info/nodes/instances/groups subcommand with deployment and scalegroup params is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "groups":
		{
			operFlagSet := basicOptions("scaling-groups instances")
			var deployment string
			var nodeType string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&nodeType, "node-type",
				"cloudify.nodes.ApplicationServer.kubernetes.Node",
				"Filter by node type")

			operFlagSet.Parse(options)

			if deployment == "" {
				fmt.Println("Please provide deployment")
				return 1
			}

			cl := getClient()
			groupedInstances, err := cl.GetDeploymentInstancesScaleGrouped(deployment, nodeType)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			for groupName, instances := range groupedInstances {
				fmt.Printf("Scale group: %v\n", groupName)
				if nodeInstancesPrint(&instances) != 0 {
					return 1
				}
			}
			return 0
		}
	case "instances":
		{
			operFlagSet := basicOptions("scaling-groups instances")
			var deployment string
			var scalegroup string
			var nodeType string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&scalegroup, "scalegroup", "",
				"The unique identifier for the scalegroup")
			operFlagSet.StringVar(&nodeType, "node-type",
				"cloudify.nodes.ApplicationServer.kubernetes.Node",
				"Filter by node type")

			operFlagSet.Parse(options)

			if deployment == "" {
				fmt.Println("Please provide deployment")
				return 1
			}
			if scalegroup == "" {
				fmt.Println("Please provide scalegroup")
				return 1
			}

			cl := getClient()
			instances, err := cl.GetDeploymentScaleGroupInstances(deployment, scalegroup, nodeType)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			return nodeInstancesPrint(instances)
		}
	case "nodes":
		{
			operFlagSet := basicOptions("scalegroups nodes")
			var deployment string
			var scalegroup string
			var nodeType string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&scalegroup, "scalegroup", "",
				"The unique identifier for the scalegroup")
			operFlagSet.StringVar(&nodeType, "node-type",
				"cloudify.nodes.ApplicationServer.kubernetes.Node",
				"Filter by node type")

			operFlagSet.Parse(options)

			if deployment == "" {
				fmt.Println("Please provide deployment")
				return 1
			}
			if scalegroup == "" {
				fmt.Println("Please provide scalegroup")
				return 1
			}

			cl := getClient()
			nodes, err := cl.GetDeploymentScaleGroupNodes(deployment, scalegroup, nodeType)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			return nodesPrint(nodes)
		}
	case "info":
		{
			operFlagSet := basicOptions("scalegroups info")
			var deployment string
			var scalegroup string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&scalegroup, "scalegroup", "",
				"The unique identifier for the scalegroup")

			operFlagSet.Parse(options)

			if deployment == "" {
				fmt.Println("Please provide deployment")
				return 1
			}
			if scalegroup == "" {
				fmt.Println("Please provide scalegroup")
				return 1
			}

			cl := getClient()
			scaleGroupObj, err := cl.GetDeploymentScaleGroup(deployment, scalegroup)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			var scaleGroups = map[string]cloudify.ScalingGroup{}
			scaleGroups[scalegroup] = *scaleGroupObj
			return scaleGroupPrint(scaleGroups)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
