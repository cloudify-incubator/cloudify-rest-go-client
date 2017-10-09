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

package rest

const ApiVersion = "v3.1"

// Credentials
type CloudifyRestClient struct {
	restURL  string
	user     string
	password string
	tenant   string
	Debug    bool
}

type CloudifyMessageInterface interface {
	ErrorCode() string
	Error() string
	TraceBack() string
}

// We need Cl prefix for make fields public and use in Marshal func
type CloudifyBaseMessage struct {
	CloudifyMessageInterface
	ClMessage         string `json:"message,omitempty"`
	ClErrorCode       string `json:"error_code,omitempty"`
	ClServerTraceback string `json:"server_traceback,omitempty"`
}

func (cm *CloudifyBaseMessage) ErrorCode() string {
	return cm.ClErrorCode
}

// Support reuse as error type
func (cm *CloudifyBaseMessage) Error() string {
	return cm.ClMessage
}

func (cm *CloudifyBaseMessage) TraceBack() string {
	return cm.ClServerTraceback
}

// Common
type CloudifyPagination struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Size   uint `json:"size"`
}

type CloudifyMetadata struct {
	Pagination CloudifyPagination `json:"pagination"`
}

type CloudifyIdWithTenant struct {
	Id              string `json:"id"`
	Tenant          string `json:"tenant_name"`
	CreatedBy       string `json:"created_by"`
	PrivateResource bool   `json:"private_resource"`
}

type CloudifyResource struct {
	CloudifyIdWithTenant
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
