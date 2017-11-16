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

// APIVersion - currently supported version of Cloudify API
const APIVersion = "v3.1"

// MessageInterface - Interface for any cloudify error resoponse
type MessageInterface interface {
	ErrorCode() string
	Error() string
	TraceBack() string
}

// BaseMessage - common part of any result from cloudify
// Note: We need Cl prefix for make fields public and use in Marshal func
// Check https://blog.golang.org/json-and-go for more info about json marshaling.
type BaseMessage struct {
	MessageInterface
	ClMessage         string `json:"message,omitempty"`
	ClErrorCode       string `json:"error_code,omitempty"`
	ClServerTraceback string `json:"server_traceback,omitempty"`
}

// ErrorCode - current error code if any
func (cm *BaseMessage) ErrorCode() string {
	return cm.ClErrorCode
}

// Error - Support reuse BaseMessage as error type
func (cm *BaseMessage) Error() string {
	return cm.ClMessage
}

// TraceBack - traceback from response
func (cm *BaseMessage) TraceBack() string {
	return cm.ClServerTraceback
}

// Pagination - common struct of any result with pagination
type Pagination struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Size   uint `json:"size"`
}

// Metadata - common struct of any result sevaral items in response
type Metadata struct {
	Pagination Pagination `json:"pagination"`
}

// ObjectIDWithTenant - common struct for any response with object id and tenant
type ObjectIDWithTenant struct {
	ID              string `json:"id"`
	Tenant          string `json:"tenant_name"`
	CreatedBy       string `json:"created_by"`
	PrivateResource bool   `json:"private_resource"`
}

// Resource - common struct for any object from cloudify with description
type Resource struct {
	ObjectIDWithTenant
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ConnectionOperationsInterface - mandatory methods for any cloudify connection
// For now implemented only http/https version
type ConnectionOperationsInterface interface {
	Get(url, acceptedContentType string) ([]byte, error)
	Delete(url string) ([]byte, error)
	Post(url string, data []byte) ([]byte, error)
	Put(url, providedContentType string, data []byte) ([]byte, error)
	SetDebug(bool)
	GetDebug() bool
}
