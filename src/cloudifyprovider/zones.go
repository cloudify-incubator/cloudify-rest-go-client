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

package cloudifyprovider

import (
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type CloudifyZones struct {
	client *cloudify.CloudifyClient
}

// GetZone is an implementation of Zones.GetZone
func (r *CloudifyZones) GetZone() (cloudprovider.Zone, error) {
	glog.Infof("GetZone")
	return cloudprovider.Zone{
		FailureDomain: "FailureDomain",
		Region:        "Region",
	}, nil
}

func NewCloudifyZones(client *cloudify.CloudifyClient) *CloudifyZones {
	return &CloudifyZones{
		client: client,
	}
}