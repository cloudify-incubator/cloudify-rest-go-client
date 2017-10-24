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

package cloudify

import (
	"encoding/json"
	"fmt"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	"log"
	"net/url"
	"time"
)

type CloudifyExecutionPost struct {
	WorkflowId   string                 `json:"workflow_id"`
	DeploymentId string                 `json:"deployment_id"`
	Parameters   map[string]interface{} `json:"parameters"`
}

func (exec *CloudifyExecutionPost) SetJsonParameters(parameters string) {
	if len(parameters) == 0 {
		exec.Parameters = map[string]interface{}{}
		return
	}

	err := json.Unmarshal([]byte(parameters), &exec.Parameters)
	if err != nil {
		log.Fatal(err)
	}
}

type CloudifyExecution struct {
	// have id, owner information
	rest.CloudifyResource
	// contain information from post
	CloudifyExecutionPost
	IsSystemWorkflow bool   `json:"is_system_workflow"`
	ErrorMessage     string `json:"error"`
	BlueprintId      string `json:"blueprint_id"`
	Status           string `json:"status"`
}

type CloudifyExecutionGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	CloudifyExecution
}

type CloudifyExecutions struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyExecution   `json:"items"`
}

// change params type if you want use non uniq values in params
func (cl *CloudifyClient) GetExecutions(params map[string]string) CloudifyExecutions {
	var executions CloudifyExecutions

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("executions?"+values.Encode(), &executions)
	if err != nil {
		log.Fatal(err)
	}

	return executions
}

// run executions without waiting
func (cl *CloudifyClient) PostExecution(exec CloudifyExecutionPost) CloudifyExecutionGet {
	var execution CloudifyExecutionGet

	var err error

	err = cl.Post("executions", exec, &execution)
	if err != nil {
		log.Fatal(err)
	}

	return execution
}

/* Check that all executions finished */
func (cl *CloudifyClient) WaitBeforeRunExecution(deploymentID string) error {
	for true {
		var params = map[string]string{}
		params["deployment_id"] = deploymentID
		executions := cl.GetExecutions(params)
		for _, execution := range executions.Items {
			if execution.WorkflowId == "create_deployment_environment" && execution.Status == "failed" {
				return fmt.Errorf(execution.ErrorMessage)
			}
			if execution.Status == "pending" || execution.Status == "started" || execution.Status == "cancelling" {
				log.Printf("Check status for %v, last status: %v", execution.Id, execution.Status)
				time.Sleep(15 * time.Second)
				break
			}
			return nil
		}
	}
	return nil
}

/* Run executions and wait results
 * exec_post: executions description for run
 * full_finish: wait to full finish
 */
func (cl *CloudifyClient) RunExecution(exec_post CloudifyExecutionPost, full_finish bool) CloudifyExecution {
	var execution CloudifyExecution
	executionGet := cl.PostExecution(exec_post)
	execution = executionGet.CloudifyExecution
	for execution.Status == "pending" || (execution.Status == "started" && full_finish) {
		log.Printf("Check status for %v, last status: %v", execution.Id, execution.Status)

		time.Sleep(15 * time.Second)

		var params = map[string]string{}
		params["id"] = execution.Id
		executions := cl.GetExecutions(params)
		if len(executions.Items) != 1 {
			log.Fatal(fmt.Errorf("Returned wrong count of results."))
		}
		execution = executions.Items[0]
	}
	return execution
}
