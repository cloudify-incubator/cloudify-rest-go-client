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
	"encoding/json"
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type baseResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type capabilitiesResponse struct {
	Attach bool `json:"attach"`
}

type initResponse struct {
	baseResponse
	Capabilities capabilitiesResponse `json:"capabilities,omitempty"`
}

type mountResponse struct {
	baseResponse
	Attached bool `json:"attached"`
}

type CloudifyConfig struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Tenant     string `json:"tenant"`
	Deployment string `json:"deployment"`
	Instance   string `json:"intance"`
}

func getConfig() (*CloudifyConfig, error) {
	var config CloudifyConfig
	var configFile = os.Getenv("CFY_CONFIG")
	if configFile == "" {
		configFile = "/etc/cloudify/mount.json"
	}

	log.Printf("Use %s as config.", configFile)

	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err_marshal := json.Unmarshal([]byte(configContent), &config)
	return &config, err_marshal
}

func initFunction() error {
	var response initResponse
	response.Status = "Success"
	response.Capabilities.Attach = false
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func runAction(config *CloudifyConfig, action string, params map[string]interface{}) error {
	cl := cloudify.NewClient(config.Host, config.User, config.Password, config.Tenant)

	log.Printf("Client version %s", cl.GetApiVersion())
	log.Printf("Run %v with %v", action, params)

	var exec cloudify.CloudifyExecutionPost
	exec.WorkflowId = "execute_operation"
	exec.DeploymentId = config.Deployment
	exec.Parameters = map[string]interface{}{}
	exec.Parameters["operation"] = action
	exec.Parameters["node_ids"] = []string{}
	exec.Parameters["type_names"] = []string{}
	exec.Parameters["run_by_dependency_order"] = false
	exec.Parameters["allow_kwargs_override"] = nil
	exec.Parameters["node_instance_ids"] = []string{config.Instance}
	exec.Parameters["operation_kwargs"] = params
	var execution cloudify.CloudifyExecution
	executionGet := cl.PostExecution(exec)
	execution = executionGet.CloudifyExecution
	for execution.Status == "pending" || execution.Status == "started" {
		log.Printf("Check status for %v, last status: %v", execution.Id, execution.Status)

		time.Sleep(15 * time.Second)

		var params = map[string]string{}
		params["id"] = execution.Id
		executions := cl.GetExecutions(params)
		if len(executions.Items) != 1 {
			return fmt.Errorf("Returned wrong count of results.")
		}
		execution = executions.Items[0]
	}

	log.Printf("Final status for %v, last status: %v", execution.Id, execution.Status)

	if execution.Status == "failed" {
		return fmt.Errorf(execution.ErrorMessage)
	}
	return nil
}

func mountFunction(config *CloudifyConfig, path, config_json string) error {
	var in_data_parsed map[string]interface{}
	err := json.Unmarshal([]byte(config_json), &in_data_parsed)
	if err != nil {
		return err
	}

	var params = map[string]interface{}{
		"path":   path,
		"params": in_data_parsed}

	err_action := runAction(config, "maintenance.mount", params)

	if err_action != nil {
		return err_action
	}

	var response mountResponse
	response.Status = "Success"
	response.Attached = true
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func unMountFunction(config *CloudifyConfig, path string) error {
	var params = map[string]interface{}{
		"path": path}

	err_action := runAction(config, "maintenance.unmount", params)

	if err_action != nil {
		return err_action
	}

	var response mountResponse
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

var versionString = "0.1"

func main() {
	var message string = "Unknown"

	f, err := os.OpenFile("/var/log/cfy-mount.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		message = err.Error()
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	log.Printf("Mount version %s, called with %+v", versionString, os.Args)

	config, config_err := getConfig()
	if config_err != nil {
		message = config_err.Error()
	} else if len(os.Args) > 1 {
		command := os.Args[1]
		if len(os.Args) == 2 && command == "init" {
			err := initFunction()
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
		if len(os.Args) == 4 && command == "mount" {
			err := mountFunction(config, os.Args[2], os.Args[3])
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
		if len(os.Args) == 3 && command == "unmount" {
			err := unMountFunction(config, os.Args[2])
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
	}

	log.Printf("Error: %v", message)

	var response baseResponse
	response.Status = "Not supported"
	response.Message = message
	json_data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(json_data))
	os.Exit(0)
}
