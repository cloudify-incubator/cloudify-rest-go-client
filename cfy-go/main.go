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
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
	"os"
)

var host string
var user string
var password string
var tenant string
var cfyDebug bool

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

	commonFlagSet.BoolVar(&cfyDebug, "debug", false,
		"Manager debug or CFY_DEBUG in env")

	return commonFlagSet
}

func getClient() *cloudify.Client {
	cl := cloudify.NewClient(host, user, password, tenant)
	if cfyDebug {
		cl.EnableDebug()
	}
	return cl
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

var versionString = "0.1"

func main() {
	f, err := os.OpenFile("/var/log/cloudify.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Logs outputs to standart output: %s\n", err.Error())
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	args, options := utils.CliArgumentsList(os.Args)
	defaultError := ("Supported commands:\n" +
		"\tblueprints        Handle blueprints on the manager\n" +
		"\tdeployments       Handle deployments on the Manager\n" +
		"\tscaling-groups    Handle scale groups on the Manager\n" +
		"\tevents            Show events from workflow executions\n" +
		"\texecutions        Handle workflow executions\n" +
		"\tnode-instances    Handle a deployment's node-instances\n" +
		"\tnodes             Handle a deployment's nodes\n" +
		"\tplugins           Handle plugins on the manager\n" +
		"\tstatus            Show manager status\n" +
		"\tkubernetes        Additional kubernetes operations\n" +
		"\tversion           Show client version.\n")

	if len(args) < 2 {
		fmt.Println(defaultError)
		return
	}

	switch args[1] {
	case "version":
		{
			fmt.Printf("CFY Go client: %s\n", versionString)
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
	case "scaling-groups":
		{
			os.Exit(scalingGroupsOptions(args, options))
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
	case "kubernetes":
		{
			os.Exit(kubernetesOptions(args, options))
		}
	default:
		{
			fmt.Println(defaultError)
			os.Exit(1)
		}
	}
}
