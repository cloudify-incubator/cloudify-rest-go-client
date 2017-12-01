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
cfy-go implements CLI for cloudify client. If we compare to official cfy
command cfy-go has implementation for only external commands.

Usage:

	$ cfy-go -h
	Supported commands:
		blueprints        Handle blueprints on the manager
		deployments       Handle deployments on the Manager
		scaling-groups    Handle scale groups on the Manager
		events            Show events from workflow executions
		executions        Handle workflow executions
		node-instances    Handle a deployment's node-instances
		nodes             Handle a deployment's nodes
		plugins           Handle plugins on the manager
		status            Show manager status
		kubernetes        Additional kubernetes operations
		version           Show client version.

Common parameters for commands required network communication:

	-debug
		Manager debug or CFY_DEBUG in env
	-host string
		Manager host name or CFY_HOST in env (default "localhost")
	-password string
		Manager user password or CFY_PASSWORD in env (default "secret")
	-tenant string
		Manager tenant or CFY_TENANT in env (default "default_tenant")
	-user string
		Manager user name or CFY_USER in env (default "admin")

Example:

	cfy-go status version -host <your manager host> -user admin -password secret -tenant default_tenant

Not implemeted commands:

	agents: Handle a deployment's agents
	install-plugins: Install plugins [locally]
	bootstrap: Bootstrap a manager.
	cluster: Handle the Cloudify Manager cluster
	dev: Run fabric tasks [manager only].
	groups: Handle deployment groups
	init: Initialize a working env
	install: Install an application blueprint [manager only].
	ldap: Set LDAP authenticator.
	logs: Handle manager service logs.
	maintenance-mode: Handle the manager's maintenance-mode.
	profiles: Handle Cloudify CLI profiles Each profile can...
	rollback: Rollback a manager to a previous version.
	secrets: Handle Cloudify secrets (key-value pairs).
	snapshots: Handle manager snapshots.
	ssh: Connect using SSH [manager only].
	teardown: Teardown a manager [manager only]
	tenants: Handle Cloudify tenants (Premium feature)
	uninstall: Uninstall an application blueprint [manager only]
	user-groups: Handle Cloudify user groups (Premium feature)
	users: Handle Cloudify users
	workflows: Handle deployment workflows

Bluprint

For use blueprint related command use cfy-go blueprints, it provide fuctionality
for manage blueprints on the manager:

	create-requirements - Create pip-requirements. Not Implemented.

	delete - Delete a blueprint [manager only]

		cfy-go blueprints delete blueprint

	download - Download a blueprint [manager only]

		cfy-go blueprints download blueprint

	get - Retrieve blueprint information [manager only]

		cfy-go blueprints list -blueprint blueprint

	inputs - Retrieve blueprint inputs [manager only]. Not Implemented.

	list - List blueprints [manager only]

		cfy-go blueprints list

		Paggination by:
			`-offset`:  the number of resources to skip.
			`-size`: the max size of the result subset to receive.

	package - Create a blueprint archive. Not Implemented.

	upload - Upload a blueprint [manager only].
		cfy-go blueprints upload new-blueprint -path src/github.com/cloudify-incubator/cloudify-rest-go-client/examples/blueprint/Minimal.yaml

	validate - Validate a blueprint. Not Implemented.

*/
package main

import (
	"fmt"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

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

			cl := getClient()
			blueprints, err := cl.GetBlueprints(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, len(blueprints.Items))
			for pos, blueprint := range blueprints.Items {
				lines[pos] = make([]string, 7)
				lines[pos][0] = blueprint.ID
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
			var blueprintPath string
			operFlagSet.StringVar(&blueprintPath, "path", "",
				"The blueprint path")
			operFlagSet.Parse(options)

			if len(blueprintPath) < 4 {
				fmt.Println("Blueprint path required")
				return 1
			}
			cl := getClient()
			blueprint, err := cl.UploadBlueprint(args[3], blueprintPath)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, 1)
			lines[0] = make([]string, 7)
			lines[0][0] = blueprint.ID
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

			cl := getClient()
			blueprintPath, err := cl.DownloadBlueprints(args[3])
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
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

			cl := getClient()
			blueprint, err := cl.DeleteBlueprints(args[3])
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, 1)
			lines[0] = make([]string, 7)
			lines[0][0] = blueprint.ID
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
