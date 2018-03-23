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

/*
Package tests - fake classes for api testing.
*/
package tests

// FakeClient - fake clent for tests
type FakeClient struct {
	// get call
	GetURL      string
	GetType     string
	GetResponse []byte
	GetError    error

	// delete call
	DeleteURL      string
	DeleteResponse []byte
	DeleteError    error

	// post call
	PostURL      string
	PostType     string
	PostData     []byte
	PostResponse []byte
	PostError    error

	// put call
	PutURL      string
	PutType     string
	PutData     []byte
	PutResponse []byte
	PutError    error

	// debug
	DebugState bool
}

// Get - mimic to real get
func (cl *FakeClient) Get(url, acceptedContentType string) ([]byte, error) {
	cl.GetURL = url
	cl.GetType = acceptedContentType
	return cl.GetResponse, cl.GetError
}

// Delete - mimic to real delete
func (cl *FakeClient) Delete(url string) ([]byte, error) {
	cl.DeleteURL = url
	return cl.DeleteResponse, cl.DeleteError
}

// Post - mimic to real post
func (cl *FakeClient) Post(url, providedContentType string, data []byte) ([]byte, error) {
	cl.PostURL = url
	cl.PostType = providedContentType
	cl.PostData = data
	return cl.PostResponse, cl.PostError
}

// Put - mimic to real put
func (cl *FakeClient) Put(url, providedContentType string, data []byte) ([]byte, error) {
	cl.PutURL = url
	cl.PutType = providedContentType
	cl.PutData = data
	return cl.PutResponse, cl.PutError
}

// SetDebug - mimic to real set debug
func (cl *FakeClient) SetDebug(state bool) {
	cl.DebugState = state
}

// GetDebug - mimic to real get debug
func (cl *FakeClient) GetDebug() bool {
	return cl.DebugState
}
