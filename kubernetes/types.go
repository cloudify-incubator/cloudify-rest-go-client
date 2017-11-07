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

package kubernetes

/*
BaseResponse - base type for all responses from mount operation
*/
type BaseResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

/*
CapabilitiesResponse - list of supported capabilities,
return `attach` cappability for now
*/
type CapabilitiesResponse struct {
	Attach bool `json:"attach"`
}

/*
InitResponse - describe result of init operation with list of supported capabilities
*/
type InitResponse struct {
	BaseResponse
	Capabilities CapabilitiesResponse `json:"capabilities,omitempty"`
}

/*
MountResponse - describe result of mount action with final state for device
*/
type MountResponse struct {
	BaseResponse
	Attached bool `json:"attached"`
}
