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
	"log"
	"net/url"
)

type CloudifyEvent struct {
	NodeInstanceId    string `json:"node_instance_id"`
	EventType         string `json:"event_type"`
	Operation         string `json:"operation"`
	BlueprintId       string `json:"blueprint_id"`
	NodeName          string `json:"node_name"`
	WorkflowId        string `json:"workflow_id"`
	ErrorCauses       string `json:"error_causes"`
	ReportedTimestamp string `json:"reported_timestamp"`
	DeploymentId      string `json:"deployment_id"`
	Type              string `json:"type"`
	ExecutionId       string `json:"execution_id"`
	Timestamp         string `json:"timestamp"`
	Message           string `json:"message"`
}

type CloudifyEvents struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyEvent       `json:"items"`
}

func (cl *CloudifyClient) GetEvents(params map[string]string) CloudifyEvents {
	var events CloudifyEvents

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("events?"+values.Encode(), &events)
	if err != nil {
		log.Fatal(err)
	}

	return events
}
