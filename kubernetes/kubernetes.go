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
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
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

func runAction(cl *cloudify.CloudifyClient, action string, params map[string]interface{}, deployment, instance string) error {
	log.Printf("Client version %s", cl.GetApiVersion())
	log.Printf("Run %v with %v", action, params)

	err := cl.WaitBeforeRunExecution(deployment)
	if err != nil {
		return err
	}
	var exec cloudify.CloudifyExecutionPost
	exec.WorkflowId = "execute_operation"
	exec.DeploymentId = deployment
	exec.Parameters = map[string]interface{}{}
	exec.Parameters["operation"] = action
	exec.Parameters["node_ids"] = []string{}
	exec.Parameters["type_names"] = []string{}
	exec.Parameters["run_by_dependency_order"] = false
	exec.Parameters["allow_kwargs_override"] = nil
	exec.Parameters["node_instance_ids"] = []string{instance}
	exec.Parameters["operation_kwargs"] = params
	execution := cl.RunExecution(exec, true)

	log.Printf("Final status for %v, last status: %v", execution.Id, execution.Status)

	if execution.Status == "failed" {
		return fmt.Errorf(execution.ErrorMessage)
	}
	return nil
}

func mountFunction(cl *cloudify.CloudifyClient, path, config_json, deployment, instance string) error {
	var in_data_parsed map[string]interface{}
	err := json.Unmarshal([]byte(config_json), &in_data_parsed)
	if err != nil {
		return err
	}

	var params = map[string]interface{}{
		"path":   path,
		"params": in_data_parsed}

	err_action := runAction(cl, "maintenance.mount", params, deployment, instance)

	if err_action != nil {
		return err_action
	}

	var response MountResponse
	response.Status = "Success"
	response.Attached = true
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func unMountFunction(cl *cloudify.CloudifyClient, path, deployment, instance string) error {
	var params = map[string]interface{}{
		"path": path}

	err_action := runAction(cl, "maintenance.unmount", params, deployment, instance)

	if err_action != nil {
		return err_action
	}

	var response MountResponse
	response.Status = "Success"
	response.Attached = false
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	} else {
		fmt.Println(string(json_data))
	}
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
		if len(args) == 3 && command == "mount" {
			err := mountFunction(cl, args[1], args[2], deployment, instance)
			if err != nil {
				message = err.Error()
			} else {
				return 0
			}
		}
		if len(args) == 2 && command == "unmount" {
			err := unMountFunction(cl, args[1], deployment, instance)
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
