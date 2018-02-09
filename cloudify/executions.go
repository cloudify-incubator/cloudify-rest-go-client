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
	"time"
)

// ExecutionPost - information for create new execution
type ExecutionPost struct {
	WorkflowID   string                 `json:"workflow_id"`
	DeploymentID string                 `json:"deployment_id"`
	Parameters   map[string]interface{} `json:"parameters"`
}

// SetJSONParameters - set parameters for execution (use before send)
func (exec *ExecutionPost) SetJSONParameters(parameters string) error {
	if len(parameters) == 0 {
		exec.Parameters = map[string]interface{}{}
		return nil
	}

	err := json.Unmarshal([]byte(parameters), &exec.Parameters)
	if err != nil {
		return err
	}
	return nil
}

// Execution - information about execution on manager
type Execution struct {
	// have id, owner information
	rest.Resource
	// contain information from post
	ExecutionPost
	IsSystemWorkflow bool   `json:"is_system_workflow"`
	ErrorMessage     string `json:"error"`
	BlueprintID      string `json:"blueprint_id"`
	Status           string `json:"status"`
}

// ExecutionGet - response from manager about selected execution
type ExecutionGet struct {
	// can be response from api
	rest.BaseMessage
	Execution
}

// Executions - response from manager about several executions
type Executions struct {
	rest.BaseMessage
	Metadata rest.Metadata `json:"metadata"`
	Items    []Execution   `json:"items"`
}

// GetExecutions - return list of execution on manager
// NOTE: change params type if you want use non uniq values in params
func (cl *Client) GetExecutions(params map[string]string) (*Executions, error) {
	var executions Executions

	values := cl.stringMapToURLValue(params)

	err := cl.Get("executions?"+values.Encode(), &executions)
	if err != nil {
		return nil, err
	}

	return &executions, nil
}

// PostExecution - run executions without waiting
func (cl *Client) PostExecution(exec ExecutionPost) (*ExecutionGet, error) {
	var execution ExecutionGet

	var err error

	err = cl.Post("executions", exec, &execution)
	if err != nil {
		return nil, err
	}

	return &execution, nil
}

// WaitBeforeRunExecution - wait while all other executions will be finished
func (cl *Client) WaitBeforeRunExecution(deploymentID string) error {
	for true {
		var params = map[string]string{}
		params["deployment_id"] = deploymentID
		executions, err := cl.GetExecutions(params)
		if err != nil {
			return err
		}
		haveUnfinished := false
		for _, execution := range executions.Items {
			if execution.WorkflowID == "create_deployment_environment" && execution.Status == "failed" {
				return fmt.Errorf(execution.ErrorMessage)
			}
			if execution.Status == "pending" || execution.Status == "started" || execution.Status == "cancelling" {
				if cl.restCl().GetDebug() {
					log.Printf("Check status for %v, last status: %v", execution.ID, execution.Status)
				}
				time.Sleep(15 * time.Second)
				haveUnfinished = true
				break
			}
		}
		if !haveUnfinished {
			return nil
		}
	}
	return nil
}

// RunExecution - Run executions and wait results
// execPost: executions description for run
// fullFinish: wait to full finish
func (cl *Client) RunExecution(execPost ExecutionPost, fullFinish bool) (*Execution, error) {
	var execution Execution
	executionGet, err := cl.PostExecution(execPost)
	if err != nil {
		return nil, err
	}
	execution = executionGet.Execution
	for execution.Status == "pending" || (execution.Status == "started" && fullFinish) {
		if cl.restCl().GetDebug() {
			log.Printf("Check status for %v, last status: %v", execution.ID, execution.Status)
		}

		time.Sleep(15 * time.Second)

		var params = map[string]string{}
		params["id"] = execution.ID
		executions, err := cl.GetExecutions(params)
		if err != nil {
			return nil, err
		}
		if len(executions.Items) != 1 {
			return nil, fmt.Errorf("returned wrong count of results")
		}
		execution = executions.Items[0]
	}
	return &execution, nil
}
