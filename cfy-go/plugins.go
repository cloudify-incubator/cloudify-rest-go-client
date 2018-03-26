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

	delete: Delete a plugin [manager only].

		cfy-go plugins delete -plugin-id <plugin-id>

	download: Download a plugin [manager only]. Not Implemented.

	get: Retrieve plugin information [manager only]. Not Implemented.

	list: List plugins [manager only]

		cfy-go plugins list

	upload: Upload a plugin [manager only].

		cfy-go plugins upload -host 172.16.168.176 -plugin-path <plugin-path>.wgn -yaml-path <yaml-path>.yaml

	validate: Validate a plugin. Not Implemented.

*/
package main

import (
	"flag"
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"log"
)

func printPlugins(plugins []cloudify.Plugin) {
	lines := make([][]string, len(plugins))
	for pos, plugin := range plugins {
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

}

func uploadPluginsCall(operFlagSet *flag.FlagSet, args, options []string) int {
	var pluginPath string
	operFlagSet.StringVar(&pluginPath, "plugin-path", "",
		"The plugin path")
	var yamlPath string
	operFlagSet.StringVar(&yamlPath, "yaml-path", "",
		"The plugin yaml path")
	var visibility string
	operFlagSet.StringVar(&visibility, "visibility", "tenant",
		"The plugin visibility")
	operFlagSet.Parse(options)
	if len(pluginPath) < 4 {
		fmt.Println("Plugin path required")
		return 1
	}
	if len(yamlPath) < 4 {
		fmt.Println("Plugin yaml file required")
		return 1
	}

	var params = map[string]string{}
	params["visibility"] = visibility
	cl := getClient()
	plugin, err := cl.UploadPlugin(params, pluginPath, yamlPath)
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}
	printPlugins([]cloudify.Plugin{plugin.Plugin})
	return 0
}

func deletePluginsCall(operFlagSet *flag.FlagSet, args, options []string) int {
	var pluginID string
	var forceParams cloudify.CallWithForce
	operFlagSet.StringVar(&pluginID, "plugin-id", "",
		"The unique identifier for the plugin")
	operFlagSet.BoolVar(&forceParams.Force, "force", false,
		"Specifies whether to force plugin deletion even if there are deployments that currently use it.")

	operFlagSet.Parse(options)

	if pluginID == "" {
		fmt.Println("Plugin Id required")
		return 1
	}

	cl := getClient()
	plugin, err := cl.DeletePlugins(pluginID, forceParams)
	if err != nil {
		log.Printf("Cloudify error: %s\n", err.Error())
		return 1
	}
	printPlugins([]cloudify.Plugin{plugin.Plugin})
	return 0
}

func listPluginsCall(operFlagSet *flag.FlagSet, args, options []string) int {
	var pluginID string
	operFlagSet.StringVar(&pluginID, "plugin-id", "",
		"The unique identifier for the plugin")

	params := parsePagination(operFlagSet, options)

	if pluginID != "" {
		params["id"] = pluginID
	}

	cl := getClient()
	plugins, err := cl.GetPlugins(params)
	if err != nil {
		log.Printf("Cloudify error: %s", err.Error())
		return 1
	}
	printPlugins(plugins.Items)
	fmt.Printf("Showed %d+%d/%d results. Use offset/size for get more.\n",
		plugins.Metadata.Pagination.Offset, len(plugins.Items),
		plugins.Metadata.Pagination.Total)

	return 0
}

func pluginsOptions(args, options []string) int {
	var pluginsCalls = []CommandInfo{{
		CommandName: "list",
		Callback:    listPluginsCall,
	}, {
		CommandName: "upload",
		Callback:    uploadPluginsCall,
	}, {
		CommandName: "delete",
		Callback:    deletePluginsCall,
	}}

	return ParseCalls(pluginsCalls, 3, args, options)
}
