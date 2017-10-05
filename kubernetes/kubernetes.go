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

package kubernetes

import (
	"encoding/json"
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"log"
)

func initFunction() error {
	var response InitResponse
	response.Status = "Success"
	response.Capabilities.Attach = false
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func Run(cl *cloudify.CloudifyClient, args []string, deployment, instance string) int {
	var message string = "Unknown"

	log.Printf("Kubernetes mount called with %+v", args)

	if len(args) > 0 {
		command := args[0]
		if len(args) == 1 && command == "init" {
			err := initFunction()
			if err != nil {
				message = err.Error()
			} else {
				return 0
			}
		}
	}
	log.Printf("Error: %v", message)

	var response BaseResponse
	response.Status = "Not supported"
	response.Message = message
	json_data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Println(string(json_data))
	return 0
}
