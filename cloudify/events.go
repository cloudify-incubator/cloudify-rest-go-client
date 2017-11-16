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
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	"net/url"
)

// Event - infromation about cloudify event
type Event struct {
	NodeInstanceID    string `json:"node_instance_id"`
	EventType         string `json:"event_type"`
	Operation         string `json:"operation"`
	BlueprintID       string `json:"blueprint_id"`
	NodeName          string `json:"node_name"`
	WorkflowID        string `json:"workflow_id"`
	ErrorCauses       string `json:"error_causes"`
	ReportedTimestamp string `json:"reported_timestamp"`
	DeploymentID      string `json:"deployment_id"`
	Type              string `json:"type"`
	ExecutionID       string `json:"execution_id"`
	Timestamp         string `json:"timestamp"`
	Message           string `json:"message"`
}

// Events - cloudify response with events list
type Events struct {
	rest.BaseMessage
	Metadata rest.Metadata `json:"metadata"`
	Items    []Event       `json:"items"`
}

// GetEvents - get events list filtered by params
func (cl *Client) GetEvents(params map[string]string) (*Events, error) {
	var events Events

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("events?"+values.Encode(), &events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}
