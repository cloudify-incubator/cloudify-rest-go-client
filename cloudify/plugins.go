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

type CloudifyPluginBase struct {
	DistributionRelease string `json:"distribution_release,omitempty"`
	PackageName         string `json:"package_name,omitempty"`
	PackageVersion      string `json:"package_version,omitempty"`
	DistributionVersion string `json:"distribution_version,omitempty"`
	SupportedPlatform   string `json:"supported_platform,omitempty"`
}

type CloudifyPlugin struct {
	rest.CloudifyIdWithTenant
	CloudifyPluginBase
	SupportedPyVersions []string `json:"supported_py_versions,omitempty"`
	UploadedAt          string   `json:"uploaded_at,omitempty"`
	ArchiveName         string   `json:"archive_name,omitempty"`
	ExcludedWheels      []string `json:"excluded_wheels,omitempty"`
	Distribution        string   `json:"distribution,omitempty"`
	PackageSource       string   `json:"package_source,omitempty"`
	Wheels              []string `json:"wheels,omitempty"`
}

type CloudifyPlugins struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyPlugin      `json:"items"`
}

func (cl *CloudifyClient) GetPlugins(params map[string]string) CloudifyPlugins {
	var plugins CloudifyPlugins

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("plugins?"+values.Encode(), &plugins)
	if err != nil {
		log.Fatal(err)
	}

	return plugins
}
