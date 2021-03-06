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
Package rest - implementation http/https connection to cloudify manager.
*/
package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// JSONContentType - type used in communication with manager
const JSONContentType = "application/json"

// DataContentType - binary only data, like archives
const DataContentType = "application/octet-stream"

// HTTPClient - Credentials for cloudify
type HTTPClient struct {
	restURL  string
	user     string
	password string
	tenant   string
	debug    bool
}

func (r *HTTPClient) debugLogf(format string, v ...interface{}) {
	if r.debug {
		log.Printf(format, v...)
	}
}

// getRequest - create new request by params
func (r *HTTPClient) getRequest(url, method string, body io.Reader) (*http.Request, error) {
	r.debugLogf("Use: %v %v:%v@%v#%s\n", method, r.user, r.password, r.restURL+url, r.tenant)

	var authString string
	authString = r.user + ":" + r.password
	req, err := http.NewRequest(method, r.restURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))
	if len(r.tenant) > 0 {
		req.Header.Add("Tenant", r.tenant)
	}

	return req, nil
}

// Get - http(s) get request
func (r *HTTPClient) Get(url, acceptedContentType string) ([]byte, error) {
	req, err := r.getRequest(url, "GET", nil)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if len(contentType) < len(acceptedContentType) || contentType[:len(acceptedContentType)] != acceptedContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if acceptedContentType == JSONContentType {
		r.debugLogf("Response %s\n", string(body))
	} else {
		r.debugLogf("Binary response length: %d\n", len(body))
	}

	return body, nil
}

// Delete - http(s) delete request
func (r *HTTPClient) Delete(url, providedContentType string, data []byte) ([]byte, error) {
	var input io.Reader
	if len(data) > 0 {
		input = bytes.NewBuffer(data)
	}
	req, err := r.getRequest(url, "DELETE", input)
	if err != nil {
		return []byte{}, err
	}
	if input != nil {
		req.Header.Set("Content-Type", providedContentType)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	contentType := resp.Header.Get("Content-Type")

	if len(contentType) < len(JSONContentType) || contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	r.debugLogf("Response %s\n", string(body))

	return body, nil
}

// Post - http(s) post request
func (r *HTTPClient) Post(url, providedContentType string, data []byte) ([]byte, error) {
	req, err := r.getRequest(url, "POST", bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", providedContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if len(contentType) < len(JSONContentType) || contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r.debugLogf("Response %s\n", string(body))

	return body, nil
}

// Put - http(s) put request
func (r *HTTPClient) Put(url, providedContentType string, data []byte) ([]byte, error) {
	req, err := r.getRequest(url, "PUT", bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", providedContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if len(contentType) < len(JSONContentType) || contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	r.debugLogf("Response %s\n", string(body))

	return body, nil
}

// GetDebug - get current debug state
func (r *HTTPClient) GetDebug() bool {
	return r.debug
}

// SetDebug - change current debug state
func (r *HTTPClient) SetDebug(state bool) {
	r.debug = state
}

// NewClient - create new http(s) client
func NewClient(host, user, password, tenant string) ConnectionOperationsInterface {
	var restCl HTTPClient
	if len(host) >= len("http://") && (host[:len("https://")] == "https://" ||
		host[:len("http://")] == "http://") {
		restCl.restURL = host + "/api/" + APIVersion + "/"
	} else {
		restCl.restURL = "http://" + host + "/api/" + APIVersion + "/"
	}
	restCl.user = user
	restCl.password = password
	restCl.tenant = tenant
	restCl.debug = false
	return &restCl
}
