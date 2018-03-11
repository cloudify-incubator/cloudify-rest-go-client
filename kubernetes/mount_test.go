/*
Copyright (c) 2018 GigaSpaces Technologies Ltd. All rights reserved

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

/*
Package kubernetes - Flex Volume Driver.
Driver implementation Flex Volume for kubernetes. Has implemetation for init,
mount, unmount calls.
*/
package kubernetes

import (
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
)

func ExampleRun() {
	cl := cloudify.NewClient(cloudify.ClientConfig{
		Host:     "localhost",
		User:     "user",
		Password: "password",
		Tenant:   "tenant"})
	Run(cl, []string{"init"}, "some-deployment", "some-instance")
	// Output: {"status":"Success","capabilities":{"attach":false}}
}
