/*
Copyright (c) 2018 GigaSpaces Technologies Ltd. All rights reserved

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
Tenants

tenants - List tenants in this instance of cloudify manager [manager only].

		- cfy-go tenants list
*/

package main

import (
	"fmt"
	"github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"strconv"
)

func tenantsOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}
	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("tenants list")
			params := parsePagination(operFlagSet, options)
			cl := getClient()
			tenants, err := cl.GetTenants(params)
			if err != nil {
				log.Printf("Cloudify error: %s\n", err.Error())
				return 1
			}
			lines := make([][]string, len(tenants.Items))
			for pos, tenant := range tenants.Items {
				lines[pos] = make([]string, 3)
				lines[pos][0] = tenant.Name
				lines[pos][1] = strconv.Itoa(tenant.Users)
				lines[pos][2] = strconv.Itoa(tenant.Groups)
			}
			utils.PrintTable([]string{"name", "users", "groups",}, lines)
			fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
				tenants.Metadata.Pagination.Offset, len(tenants.Items),
				tenants.Metadata.Pagination.Total)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}
