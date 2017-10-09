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
)

type CloudifyVersion struct {
	rest.CloudifyBaseMessage
	Date    string `json:"date"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Commit  string `json:"commit"`
}

type CloudifyInstanceStatus struct {
	LoadState   string `json:"LoadState"`
	Description string `json:"Description"`
	State       string `json:"state"`
	MainPID     uint   `json:"MainPID"`
	Id          string `json:"Id"`
	ActiveState string `json:"ActiveState"`
	SubState    string `json:"SubState"`
}

type CloudifyInstanceService struct {
	Instances   []CloudifyInstanceStatus `json:"instances"`
	DisplayName string                   `json:"display_name"`
}

func (s CloudifyInstanceService) Status() string {
	var state string = "unknown"

	for _, instance := range s.Instances {
		if state != "failed" {
			state = instance.State
		}
	}

	return state
}

type CloudifyStatus struct {
	rest.CloudifyBaseMessage
	Status   string                    `json:"status"`
	Services []CloudifyInstanceService `json:"services"`
}

func (cl *CloudifyClient) GetVersion() CloudifyVersion {
	var ver CloudifyVersion

	err := cl.Get("version", &ver)
	if err != nil {
		log.Fatal(err)
	}

	return ver
}

func (cl *CloudifyClient) GetStatus() CloudifyStatus {
	var stat CloudifyStatus

	err := cl.Get("status", &stat)
	if err != nil {
		log.Fatal(err)
	}

	return stat
}
