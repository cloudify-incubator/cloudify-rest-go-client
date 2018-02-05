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
Deployments

deployments - Handle deployments on the Manager

	create - Create a deployment [manager only]. Partially implemented, you can set inputs only as json string.

		cfy-go deployments create deployment  -blueprint blueprint --inputs '{"ip": "b"}'

	delete - Delete a deployment [manager only]
		cfy-go deployments delete  deployment


	inputs - Show deployment inputs [manager only]. Not Implemented.

	list - List deployments [manager only].
		cfy-go deployments list

	Paggination by:
		`-offset`:  the number of resources to skip.
		`-size`: the max size of the result subset to receive.

	outputs - Show deployment outputs [manager only]

		cfy-go deployments inputs -deployment deployment

	update - Update a deployment [manager only]. Not Implemented.

	scaling-groups - check limits for scaling group

		cfy-go deployments scaling-groups -deployment <deployment_name>

	groups - list of node groups

		cfy-go deployments groups -deployment <deployment_name>
*/
package main

import (
	"flag"
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"strings"
)

func deploymentsFilter(operFlagSet *flag.FlagSet, options []string) (*cloudify.Deployments, error) {
	var deployment string
	operFlagSet.StringVar(&deployment, "deployment", "",
		"The unique identifier for the deployment")

	params := parsePagination(operFlagSet, options)

	if deployment != "" {
		params["id"] = deployment
	}

	cl := getClient()
	return cl.GetDeployments(params)
}

func groupPrint(deploymentScalingGroups map[string]cloudify.NodeGroup) int {
	lines := make([][]string, len(deploymentScalingGroups))
	var pos int
	if deploymentScalingGroups != nil {
		for groupName, nodeGroup := range deploymentScalingGroups {
			lines[pos] = make([]string, 2)
			lines[pos][0] = groupName
			lines[pos][1] = strings.Join(nodeGroup.Members, ", ")
			pos++
		}
	}
	utils.PrintTable([]string{
		"Group name", "Members",
	}, lines)
	return 0
}

func deploymentsOptions(args, options []string) int {
	defaultError := "list/create/delete/inputs/outputs/groups/scaling-groups subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "scaling-groups":
		{
			operFlagSet := basicOptions("deployments scale-groups")
			deployments, err := deploymentsFilter(operFlagSet, options)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			for _, deployment := range deployments.Items {
				fmt.Printf("Scale group in: %v\n", deployment.ID)
				scaleGroupPrint(deployment.ScalingGroups)
			}
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				deployments.Metadata.Pagination.Offset, len(deployments.Items),
				deployments.Metadata.Pagination.Total)
		}
	case "groups":
		{
			operFlagSet := basicOptions("deployments groups")
			deployments, err := deploymentsFilter(operFlagSet, options)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			for _, deployment := range deployments.Items {
				fmt.Printf("Node Group in: %v\n", deployment.ID)
				groupPrint(deployment.Groups)
			}
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				deployments.Metadata.Pagination.Offset, len(deployments.Items),
				deployments.Metadata.Pagination.Total)
		}
	case "list":
		{
			operFlagSet := basicOptions("deployments list")
			deployments, err := deploymentsFilter(operFlagSet, options)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, len(deployments.Items))
			for pos, deployment := range deployments.Items {
				var scaleGroups = []string{}
				if deployment.ScalingGroups != nil {
					for groupName := range deployment.ScalingGroups {
						scaleGroups = append(scaleGroups, groupName)
					}
				}
				lines[pos] = make([]string, 7)
				lines[pos][0] = deployment.ID
				lines[pos][1] = deployment.BlueprintID
				lines[pos][2] = deployment.CreatedAt
				lines[pos][3] = deployment.UpdatedAt
				lines[pos][4] = deployment.Tenant
				lines[pos][5] = deployment.CreatedBy
				lines[pos][6] = strings.Join(scaleGroups, ", ")
			}
			utils.PrintTable([]string{
				"id", "blueprint_id", "created_at", "updated_at",
				"tenant_name", "created_by", "scale_groups",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				deployments.Metadata.Pagination.Offset, len(deployments.Items),
				deployments.Metadata.Pagination.Total)
		}
	case "create":
		{
			operFlagSet := basicOptions("deployments list <deployment id>")
			if len(args) < 4 {
				fmt.Println("Deployment Id required")
				return 1
			}

			var blueprint string
			var jsonInputs string
			operFlagSet.StringVar(&blueprint, "blueprint", "",
				"The unique identifier for the blueprint")
			operFlagSet.StringVar(&jsonInputs, "inputs", "{}",
				"The json input string")
			operFlagSet.Parse(options)

			var depl cloudify.DeploymentPost
			depl.BlueprintID = blueprint
			depl.SetJSONInputs(jsonInputs)

			cl := getClient()
			deployment, err := cl.CreateDeployments(args[3], depl)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}

			lines := make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.ID
			lines[0][1] = deployment.BlueprintID
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			utils.PrintTable([]string{
				"id", "blueprint_id", "created_at", "updated_at",
				"tenant_name", "created_by",
			}, lines)
		}
	case "outputs":
		{
			operFlagSet := basicOptions("deployments outputs")
			deployments, err := deploymentsFilter(operFlagSet, options)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			if len(deployments.Items) != 1 {
				fmt.Println("Please recheck list of deployments")
				return 1
			}
			jsonOutputs, err := deployments.Items[0].GetJSONOutputs()
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			fmt.Printf("Deployment outputs: %+v\n", jsonOutputs)
		}
	case "inputs":
		{
			operFlagSet := basicOptions("deployments inputs")
			deployments, err := deploymentsFilter(operFlagSet, options)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			if len(deployments.Items) != 1 {
				fmt.Println("Please recheck list of deployments")
				return 1
			}
			jsonInputs, err := deployments.Items[0].GetJSONInputs()
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			fmt.Printf("Deployment inputs: %+v\n", jsonInputs)
		}
	case "delete":
		{
			operFlagSet := basicOptions("deployments delete <deployment id>")
			if len(args) < 4 {
				fmt.Println("Deployment Id required")
				return 1
			}

			operFlagSet.Parse(options)

			cl := getClient()
			deployment, err := cl.DeleteDeployments(args[3])
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.ID
			lines[0][1] = deployment.BlueprintID
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			utils.PrintTable([]string{
				"id", "blueprint_id", "created_at", "updated_at",
				"tenant_name", "created_by",
			}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
