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
Plugins

plugins - Handle plugins on the manager

	delete: Delete a plugin [manager only]. Not Implemented.

	download: Download a plugin [manager only]. Not Implemented.

	get: Retrieve plugin information [manager only]. Not Implemented.

	list: List plugins [manager only]

		cfy-go plugins list

	upload: Upload a plugin [manager only]. Not Implemented.

	validate: Validate a plugin. Not Implemented.

*/
package main

import (
	"fmt"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

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

			cl := getClient()
			plugins, err := cl.GetPlugins(params)
			if err != nil {
				log.Printf("Cloudify error: %s", err.Error())
				return 1
			}
			lines := make([][]string, len(plugins.Items))
			for pos, plugin := range plugins.Items {
				lines[pos] = make([]string, 9)
				lines[pos][0] = plugin.ID
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
