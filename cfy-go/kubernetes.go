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
Kubernetes

kubernetes - related commands:

	init - Return json in kubernetes format for use as init script responce

		cfy-go kubernetes init

	mount - Return json in kubernetes format for use as mount script responce

		cfy-go kubernetes mount /tmp/someunxists '{"kubernetes.io/fsType":"ext4",... "volumegroup":"kube_vg"}' -deployment slave -instance kubenetes_slave_*

	unmount - Return json in kubernetes format for use as unmount script responce

		cfy-go kubernetes unmount /tmp/someunxists -deployment slave -instance kubenetes_slave_*
*/
package main

import (
	"fmt"
	kubernetes "github.com/cloudify-incubator/cloudify-rest-go-client/kubernetes"
)

func kubernetesOptions(args, options []string) int {
	defaultError := "init/mount/unmount subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("kubernetes")

	var deployment string
	operFlagSet.StringVar(&deployment, "deployment", "",
		"The unique identifier for the deployment")

	var instance string
	operFlagSet.StringVar(&instance, "instance", "",
		"The unique identifier for the instance")

	operFlagSet.Parse(options)

	cl := getClient()

	if kubernetes.Run(cl, args[2:], deployment, instance) != 0 {
		fmt.Println(defaultError)
		return 1
	}
	return 0
}
