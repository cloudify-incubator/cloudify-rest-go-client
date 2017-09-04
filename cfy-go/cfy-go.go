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

package main

import (
	"flag"
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	utils "github.com/0lvin-cfy/cloudify-rest-go-client/cloudifyutils"
	"os"
)

var host string
var user string
var password string
var tenant string

func basicOptions(name string) *flag.FlagSet {
	var commonFlagSet *flag.FlagSet
	commonFlagSet = flag.NewFlagSet(name, flag.ExitOnError)

	var defaultHost = os.Getenv("CFY_HOST")
	if defaultHost == "" {
		defaultHost = "localhost"
	}
	commonFlagSet.StringVar(&host, "host", defaultHost,
		"Manager host name or CFY_HOST in env")

	var defaultUser = os.Getenv("CFY_USER")
	if defaultUser == "" {
		defaultUser = "admin"
	}
	commonFlagSet.StringVar(&user, "user", defaultUser,
		"Manager user name or CFY_USER in env")

	var defaultPassword = os.Getenv("CFY_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "secret"
	}
	commonFlagSet.StringVar(&password, "password", defaultPassword,
		"Manager user password or CFY_PASSWORD in env")

	var defaultTenant = os.Getenv("CFY_TENANT")
	if defaultTenant == "" {
		defaultTenant = "default_tenant"
	}
	commonFlagSet.StringVar(&tenant, "tenant", defaultTenant,
		"Manager tenant or CFY_TENANT in env")

	return commonFlagSet
}

func infoOptions(args, options []string) int {
	defaultError := "state/version subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "state":
		{
			operFlagSet := basicOptions("status state")
			operFlagSet.Parse(options)

			cl := cloudify.NewClient(host, user, password, tenant)
			stat := cl.GetStatus()
			fmt.Printf("Retrieving manager services status... [ip=%v]\n", host)
			fmt.Printf("Manager status: %v\n", stat.Status)
			fmt.Printf("Services:\n")
			var lines [][]string = make([][]string, len(stat.Services))
			for pos, service := range stat.Services {
				lines[pos] = make([]string, 2)
				lines[pos][0] = service.DisplayName
				lines[pos][1] = service.Status()
			}
			utils.PrintTable([]string{"service", "status"}, lines)
		}
	case "version":
		{
			operFlagSet := basicOptions("status version")
			operFlagSet.Parse(options)

			cl := cloudify.NewClient(host, user, password, tenant)
			ver := cl.GetVersion()
			fmt.Printf("Retrieving manager services version... [ip=%v]\n", host)
			utils.PrintTable([]string{"Version", "Edition", "Api Version"},
				[][]string{{ver.Version, ver.Edition, cloudify.ApiVersion}})
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func blueprintsOptions(args, options []string) int {
	defaultError := "list/delete/download/upload subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}
	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("blueprints list")
			var blueprint string
			operFlagSet.StringVar(&blueprint, "blueprint", "",
				"The unique identifier for the blueprint")

			params := parsePagination(operFlagSet, options)

			if blueprint != "" {
				params["id"] = blueprint
			}

			cl := cloudify.NewClient(host, user, password, tenant)
			blueprints := cl.GetBlueprints(params)
			var lines [][]string = make([][]string, len(blueprints.Items))
			for pos, blueprint := range blueprints.Items {
				lines[pos] = make([]string, 7)
				lines[pos][0] = blueprint.Id
				lines[pos][1] = blueprint.Description
				lines[pos][2] = blueprint.MainFileName
				lines[pos][3] = blueprint.CreatedAt
				lines[pos][4] = blueprint.UpdatedAt
				lines[pos][5] = blueprint.Tenant
				lines[pos][6] = blueprint.CreatedBy
			}
			utils.PrintTable([]string{
				"id", "description", "main_file_name", "created_at",
				"updated_at", "tenant_name", "created_by",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				blueprints.Metadata.Pagination.Offset, len(blueprints.Items),
				blueprints.Metadata.Pagination.Total)
		}
	case "upload":
		{
			operFlagSet := basicOptions("blueprints upload")
			if len(args) < 4 {
				fmt.Println("Blueprint Id required")
				return 1
			}
			var blueprint_path string
			operFlagSet.StringVar(&blueprint_path, "path", "",
				"The blueprint path")
			operFlagSet.Parse(options)

			if len(blueprint_path) < 4 {
				fmt.Println("Blueprint path required")
				return 1
			}
			cl := cloudify.NewClient(host, user, password, tenant)
			blueprint := cl.UploadBlueprint(args[3], blueprint_path)
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 7)
			lines[0][0] = blueprint.Id
			lines[0][1] = blueprint.Description
			lines[0][2] = blueprint.MainFileName
			lines[0][3] = blueprint.CreatedAt
			lines[0][4] = blueprint.UpdatedAt
			lines[0][5] = blueprint.Tenant
			lines[0][6] = blueprint.CreatedBy
			utils.PrintTable([]string{
				"id", "description", "main_file_name", "created_at",
				"updated_at", "tenant_name", "created_by",
			}, lines)
		}
	case "download":
		{
			operFlagSet := basicOptions("blueprints download")
			if len(args) < 4 {
				fmt.Println("Blueprint Id required")
				return 1
			}
			operFlagSet.Parse(options)

			cl := cloudify.NewClient(host, user, password, tenant)
			blueprintPath := cl.DownloadBlueprints(args[3])
			fmt.Printf("Blueprint saved to %s\n", blueprintPath)
		}
	case "delete":
		{
			operFlagSet := basicOptions("blueprints delete")
			if len(args) < 4 {
				fmt.Println("Blueprint Id required")
				return 1
			}
			operFlagSet.Parse(options)

			cl := cloudify.NewClient(host, user, password, tenant)
			blueprint := cl.DeleteBlueprints(args[3])
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 7)
			lines[0][0] = blueprint.Id
			lines[0][1] = blueprint.Description
			lines[0][2] = blueprint.MainFileName
			lines[0][3] = blueprint.CreatedAt
			lines[0][4] = blueprint.UpdatedAt
			lines[0][5] = blueprint.Tenant
			lines[0][6] = blueprint.CreatedBy
			utils.PrintTable([]string{
				"id", "description", "main_file_name", "created_at",
				"updated_at", "tenant_name", "created_by",
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

func parsePagination(operFlagSet *flag.FlagSet, options []string) map[string]string {
	var pageSize int
	var pageOffset int
	operFlagSet.IntVar(&pageSize, "size", 100, "Page size.")
	operFlagSet.IntVar(&pageOffset, "offset", 0, "Page offset.")
	operFlagSet.Parse(options)

	var params = map[string]string{}
	params["_size"] = fmt.Sprintf("%d", pageSize)
	params["_offset"] = fmt.Sprintf("%d", pageOffset)

	return params
}

func deploymentsFilter(operFlagSet *flag.FlagSet, options []string) cloudify.CloudifyDeployments {
	var deployment string
	operFlagSet.StringVar(&deployment, "deployment", "",
		"The unique identifier for the deployment")

	params := parsePagination(operFlagSet, options)

	if deployment != "" {
		params["id"] = deployment
	}

	cl := cloudify.NewClient(host, user, password, tenant)
	return cl.GetDeployments(params)
}

func deploymentsOptions(args, options []string) int {
	defaultError := "list/create/delete/inputs/outputs subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("deployments list")
			deployments := deploymentsFilter(operFlagSet, options)
			var lines [][]string = make([][]string, len(deployments.Items))
			for pos, deployment := range deployments.Items {
				lines[pos] = make([]string, 6)
				lines[pos][0] = deployment.Id
				lines[pos][1] = deployment.BlueprintId
				lines[pos][2] = deployment.CreatedAt
				lines[pos][3] = deployment.UpdatedAt
				lines[pos][4] = deployment.Tenant
				lines[pos][5] = deployment.CreatedBy
			}
			utils.PrintTable([]string{
				"id", "blueprint_id", "created_at", "updated_at",
				"tenant_name", "created_by",
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

			var depl cloudify.CloudifyDeploymentPost
			depl.BlueprintId = blueprint
			depl.SetJsonInputs(jsonInputs)

			cl := cloudify.NewClient(host, user, password, tenant)
			deployment := cl.CreateDeployments(args[3], depl)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
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
			deployments := deploymentsFilter(operFlagSet, options)
			if len(deployments.Items) != 1 {
				fmt.Println("Please recheck list of deployments")
				return 1
			}

			fmt.Printf("Deployment outputs: %+v\n", deployments.Items[0].GetJsonOutputs())
		}
	case "inputs":
		{
			operFlagSet := basicOptions("deployments inputs")
			deployments := deploymentsFilter(operFlagSet, options)
			if len(deployments.Items) != 1 {
				fmt.Println("Please recheck list of deployments")
				return 1
			}

			fmt.Printf("Deployment inputs: %+v\n", deployments.Items[0].GetJsonInputs())
		}
	case "delete":
		{
			operFlagSet := basicOptions("deployments delete <deployment id>")
			if len(args) < 4 {
				fmt.Println("Deployment Id required")
				return 1
			}

			operFlagSet.Parse(options)

			cl := cloudify.NewClient(host, user, password, tenant)
			deployment := cl.DeleteDeployments(args[3])
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
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

			cl := cloudify.NewClient(host, user, password, tenant)
			executions := cl.GetExecutions(params)
			var lines [][]string = make([][]string, len(executions.Items))
			for pos, execution := range executions.Items {
				lines[pos] = make([]string, 8)
				lines[pos][0] = execution.Id
				lines[pos][1] = execution.WorkflowId
				lines[pos][2] = execution.Status
				lines[pos][3] = execution.DeploymentId
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

			var exec cloudify.CloudifyExecutionPost
			exec.WorkflowId = args[3]
			exec.DeploymentId = deployment
			exec.SetJsonParameters(jsonParams)

			cl := cloudify.NewClient(host, user, password, tenant)
			execution := cl.PostExecution(exec)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 8)
			lines[0][0] = execution.Id
			lines[0][1] = execution.WorkflowId
			lines[0][2] = execution.Status
			lines[0][3] = execution.DeploymentId
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

func nodesOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("nodes list")
			var node string
			var deployment string
			operFlagSet.StringVar(&node, "node", "",
				"The unique identifier for the node")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")

			params := parsePagination(operFlagSet, options)

			if node != "" {
				params["id"] = node
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}

			cl := cloudify.NewClient(host, user, password, tenant)
			nodes := cl.GetNodes(params)
			var lines [][]string = make([][]string, len(nodes.Items))
			for pos, node := range nodes.Items {
				lines[pos] = make([]string, 9)
				lines[pos][0] = node.Id
				lines[pos][1] = node.DeploymentId
				lines[pos][2] = node.BlueprintId
				lines[pos][3] = node.HostId
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
				fmt.Printf("Properties: %s\n", nodes.Items[0].GetJsonProperties())
			} else {
				fmt.Printf("Limit to one row if you want to check Properties\n")
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

func nodeInstancesOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("node-instances list")
			var node string
			var deployment string
			operFlagSet.StringVar(&node, "node", "",
				"The unique identifier for the node")
			operFlagSet.StringVar(&deployment, "deployment", "",
				"The unique identifier for the deployment")

			params := parsePagination(operFlagSet, options)

			if node != "" {
				params["node_id"] = node
			}
			if deployment != "" {
				params["deployment_id"] = deployment
			}

			cl := cloudify.NewClient(host, user, password, tenant)
			nodeInstances := cl.GetNodeInstances(params)
			var lines [][]string = make([][]string, len(nodeInstances.Items))
			for pos, nodeInstance := range nodeInstances.Items {
				lines[pos] = make([]string, 7)
				lines[pos][0] = nodeInstance.Id
				lines[pos][1] = nodeInstance.DeploymentId
				lines[pos][2] = nodeInstance.HostId
				lines[pos][3] = nodeInstance.NodeId
				lines[pos][4] = nodeInstance.State
				lines[pos][5] = nodeInstance.Tenant
				lines[pos][6] = nodeInstance.CreatedBy
			}
			utils.PrintTable([]string{
				"Id", "Deployment id", "Host id", "Node id", "State", "Tenant",
				"created_by",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				nodeInstances.Metadata.Pagination.Offset, len(nodeInstances.Items),
				nodeInstances.Metadata.Pagination.Total)
			if len(nodeInstances.Items) == 1 {
				fmt.Printf("Runtime properties: %s\n", nodeInstances.Items[0].GetJsonRuntimeProperties())
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

			cl := cloudify.NewClient(host, user, password, tenant)
			events := cl.GetEvents(params)
			var lines [][]string = make([][]string, len(events.Items))
			for pos, event := range events.Items {
				lines[pos] = make([]string, 5)
				lines[pos][0] = event.Timestamp
				lines[pos][1] = event.DeploymentId
				lines[pos][2] = event.NodeInstanceId
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

func pluginsOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("plugins list")
			params := parsePagination(operFlagSet, options)

			cl := cloudify.NewClient(host, user, password, tenant)
			plugins := cl.GetPlugins(params)
			var lines [][]string = make([][]string, len(plugins.Items))
			for pos, plugin := range plugins.Items {
				lines[pos] = make([]string, 9)
				lines[pos][0] = plugin.Id
				lines[pos][1] = plugin.PackageName
				lines[pos][2] = plugin.PackageVersion
				lines[pos][3] = plugin.Distribution
				lines[pos][4] = plugin.SupportedPlatform
				lines[pos][5] = plugin.DistributionRelease
				lines[pos][6] = plugin.UploadedAt
				lines[pos][7] = plugin.Tenant
				lines[pos][8] = plugin.CreatedBy
			}
			utils.PrintTable([]string{
				"Id", "Package name", "Package version", "Distribution",
				"Supported platform", "Distribution release", "Uploaded at",
				"Tenant", "Created by",
			}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				plugins.Metadata.Pagination.Offset, len(plugins.Items),
				plugins.Metadata.Pagination.Total)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

var versionString = "0.1"

func main() {

	args, options := utils.CliArgumentsList(os.Args)
	var defaultError string = ("Supported commands:\n" +
		"\tblueprints        Handle blueprints on the manager\n" +
		"\tdeployments       Handle deployments on the Manager\n" +
		"\tevents            Show events from workflow executions\n" +
		"\texecutions        Handle workflow executions\n" +
		"\tnode-instances    Handle a deployment's node-instances\n" +
		"\tnodes             Handle a deployment's nodes\n" +
		"\tplugins           Handle plugins on the manager\n" +
		"\tstatus            Show manager status\n" +
		"\tversion           Show client version.\n")

	if len(args) < 2 {
		fmt.Println(defaultError)
		return
	}

	switch args[1] {
	case "version":
		{
			fmt.Printf("CFY Go client: %s/%s\n", cloudify.ApiVersion, versionString)
			os.Exit(0)
		}
	case "status":
		{
			os.Exit(infoOptions(args, options))
		}
	case "blueprints":
		{
			os.Exit(blueprintsOptions(args, options))
		}
	case "deployments":
		{
			os.Exit(deploymentsOptions(args, options))
		}
	case "executions":
		{
			os.Exit(executionsOptions(args, options))
		}
	case "plugins":
		{
			os.Exit(pluginsOptions(args, options))
		}
	case "events":
		{
			os.Exit(eventsOptions(args, options))
		}
	case "nodes":
		{
			os.Exit(nodesOptions(args, options))
		}
	case "node-instances":
		{
			os.Exit(nodeInstancesOptions(args, options))
		}
	default:
		{
			fmt.Println(defaultError)
			os.Exit(1)
		}
	}
}
