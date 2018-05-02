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
Container

container - Run command in container

	Run: Run command in container

		cfy-go container run -base container-place/base -- cfy profile use local

*/
package main

import (
	"flag"
	"fmt"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	container "github.com/cloudify-incubator/cloudify-rest-go-client/container"
	"os"
	"path"
)

func containerRunCall(operFlagSet *flag.FlagSet, args, options []string) int {
	argsCalls, commandList := utils.CliSubArgumentsList(options)
	var baseDir string
	var dataDir string
	var workDir string
	operFlagSet.StringVar(&baseDir, "base", "",
		"Base dir with cloudify container")
	operFlagSet.StringVar(&dataDir, "data", "",
		"Data dir for same changes in container")
	operFlagSet.StringVar(&workDir, "work", "",
		"Work dir for temporary save information")

	operFlagSet.Parse(argsCalls)

	if baseDir == "" {
		fmt.Println("We need base dir with cloudify container data")
		return 1
	}

	if len(commandList) == 0 {
		commandList = []string{"/bin/sh"}
	}

	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			fmt.Printf("Cloudify error: %s\n", err.Error())
			return 1
		}
	}

	if dataDir == "" {
		dataDir = path.Join(workDir, "data")
	}

	return container.Run(baseDir, dataDir, workDir, commandList)
}

func containerOptions(args, options []string) int {
	var pluginsCalls = []CommandInfo{{
		CommandName: "run",
		Callback:    containerRunCall,
	}}

	return ParseCalls(pluginsCalls, 3, args, options)
}
