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
events - Show events from workflow executions

	delete - Delete deployment events [manager only]. Not Implemented.

	list - List deployments events [manager only]

		cfy-go events list

	Paggination by:
		`-offset`:  the number of resources to skip.
		`-size`: the max size of the result subset to receive.

	Supported filters:
		`blueprint`: The unique identifier for the blueprint
		`deployment`: The unique identifier for the deployment
		`execution`: The unique identifier for the execution
*/
package main

import (
	"fmt"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

func eventsOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("events list")
			var blueprint string
			var deployment string
			var execution string
			operFlagSet.StringVar(&blueprint, "blueprint", "",
				"The unique identifier for the blueprint")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&execution, "execution", "",
				"The unique identifier for the execution")

			params := parsePagination(operFlagSet, options)

			if blueprint != "" {
				params["blueprint_id"] = blueprint
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}
			if execution != "" {
				params["execution_id"] = execution
			}

			cl := getClient()
			events, err := cl.GetEvents(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, len(events.Items))
			for pos, event := range events.Items {
				lines[pos] = make([]string, 5)
				lines[pos][0] = event.Timestamp
				lines[pos][1] = event.DeploymentID
				lines[pos][2] = event.NodeInstanceID
				lines[pos][3] = event.Operation
				lines[pos][4] = event.Message
			}
			utils.PrintTable([]string{
				"Timestamp", "Deployment", "InstanceId", "Operation",
				"Message",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				events.Metadata.Pagination.Offset, len(events.Items),
				events.Metadata.Pagination.Total)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
