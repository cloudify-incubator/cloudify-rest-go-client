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
Executions

executions - Handle workflow executions

	cancel: Cancel a workflow execution [manager only]. Not Implemented.

	get: Retrieve execution information [manager only]. Not Implemented.

	list: List deployment executions [manager only].

		cfy-go executions list
		cfy-go executions list -deployment deployment

	Paggination by:
		`-offset`:  the number of resources to skip.
		`-size`: the max size of the result subset to receive.

	start: Execute a workflow [manager only]. Partially implemented, you can set params only as json string.

		cfy-go executions start uninstall -deployment deployment
*/
package main

import (
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

func executionsOptions(args, options []string) int {
	defaultError := "list/start subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("executions list")

			var deployment string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.Parse(options)

			params := parsePagination(operFlagSet, options)

			if deployment != "" {
				params["deployment_id"] = deployment
			}

			cl := getClient()
			executions, err := cl.GetExecutions(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, len(executions.Items))
			for pos, execution := range executions.Items {
				lines[pos] = make([]string, 8)
				lines[pos][0] = execution.ID
				lines[pos][1] = execution.WorkflowID
				lines[pos][2] = execution.Status
				lines[pos][3] = execution.DeploymentID
				lines[pos][4] = execution.CreatedAt
				lines[pos][5] = execution.ErrorMessage
				lines[pos][6] = execution.Tenant
				lines[pos][7] = execution.CreatedBy
			}
			utils.PrintTable([]string{
				"id", "workflow_id", "status", "deployment_id", "created_at",
				"error", "tenant_name", "created_by",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				executions.Metadata.Pagination.Offset, len(executions.Items),
				executions.Metadata.Pagination.Total)
		}
	case "start":
		{
			operFlagSet := basicOptions("executions start <workflow id>")
			if len(args) < 4 {
				fmt.Println("Workflow Id required")
				return 1
			}

			var deployment string
			var jsonParams string
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")
			operFlagSet.StringVar(&jsonParams, "params", "{}",
				"The json params string")
			operFlagSet.Parse(options)

			var exec cloudify.ExecutionPost
			exec.WorkflowID = args[3]
			exec.DeploymentID = deployment
			exec.SetJSONParameters(jsonParams)

			cl := getClient()
			execution, err := cl.PostExecution(exec)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}

			lines := make([][]string, 1)
			lines[0] = make([]string, 8)
			lines[0][0] = execution.ID
			lines[0][1] = execution.WorkflowID
			lines[0][2] = execution.Status
			lines[0][3] = execution.DeploymentID
			lines[0][4] = execution.CreatedAt
			lines[0][5] = execution.ErrorMessage
			lines[0][6] = execution.Tenant
			lines[0][7] = execution.CreatedBy
			utils.PrintTable([]string{
				"id", "workflow_id", "status", "deployment_id", "created_at",
				"error", "tenant_name", "created_by",
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
