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
)

// Version - information about manager version
type Version struct {
	rest.BaseMessage
	Date    string `json:"date"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Commit  string `json:"commit"`
}

// InstanceStatus - system service status
type InstanceStatus struct {
	LoadState   string `json:"LoadState"`
	Description string `json:"Description"`
	State       string `json:"state"`
	MainPID     uint   `json:"MainPID"`
	ID          string `json:"Id"`
	ActiveState string `json:"ActiveState"`
	SubState    string `json:"SubState"`
}

// InstanceService - information about system service started on manager
type InstanceService struct {
	Instances   []InstanceStatus `json:"instances"`
	DisplayName string           `json:"display_name"`
}

// Status - current status for service
func (s InstanceService) Status() string {
	state := "unknown"

	for _, instance := range s.Instances {
		if state != "failed" {
			state = instance.State
		}
	}

	return state
}

// Status - response from server about current status
type Status struct {
	rest.StrStatusMessage
	Services []InstanceService `json:"services"`
}

// GetVersion - manager version
func (cl *Client) GetVersion() (*Version, error) {
	var ver Version

	err := cl.Get("version", &ver)
	if err != nil {
		return nil, err
	}

	return &ver, nil
}

// GetStatus - manager status
func (cl *Client) GetStatus() (*Status, error) {
	var stat Status

	err := cl.Get("status", &stat)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}
