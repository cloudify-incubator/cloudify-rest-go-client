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

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/*
JSONContentType - type used in communication with manager
*/
const JSONContentType = "application/json"

/*
DataContentType - binary only data, like archives
*/
const DataContentType = "application/octet-stream"

func (r *CloudifyRestClient) GetRequest(url, method string, body io.Reader) (*http.Request, error) {
	if r.Debug {
		log.Printf("Use: %v:%v@%v#%s\n", r.user, r.password, r.restURL+url, r.tenant)
	}

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

func (r *CloudifyRestClient) Get(url, acceptedContentType string) ([]byte, error) {
	req, err := r.GetRequest(url, "GET", nil)
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

	if r.Debug {
		if acceptedContentType == JSONContentType {
			log.Printf("Response %s\n", string(body))
		} else {
			log.Printf("Binary response length: %d\n", len(body))
		}
	}

	return body, nil
}

func (r *CloudifyRestClient) Delete(url string) ([]byte, error) {
	req, err := r.GetRequest(url, "DELETE", nil)
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

	if contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if r.Debug {
		log.Printf("Response %s\n", string(body))
	}

	return body, nil
}

func (r *CloudifyRestClient) Post(url string, data []byte) ([]byte, error) {
	req, err := r.GetRequest(url, "POST", bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", JSONContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if r.Debug {
		log.Printf("Response %s\n", string(body))
	}

	return body, nil
}

func (r *CloudifyRestClient) Put(url, providedContentType string, data []byte) ([]byte, error) {
	req, err := r.GetRequest(url, "PUT", bytes.NewBuffer(data))
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

	if contentType[:len(JSONContentType)] != JSONContentType {
		return []byte{}, fmt.Errorf("Wrong content type: %+v", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if r.Debug {
		log.Printf("Response %s\n", string(body))
	}

	return body, nil
}

func NewClient(host, user, password, tenant string) *CloudifyRestClient {
	var restCl CloudifyRestClient
	if (host[:len("https://")] == "https://" ||
		host[:len("http://")] == "http://") && (len(host) >= len("http://")) {
		restCl.restURL = host + "/api/" + APIVersion + "/"
	} else {
		restCl.restURL = "http://" + host + "/api/" + APIVersion + "/"
	}
	restCl.user = user
	restCl.password = password
	restCl.tenant = tenant
	restCl.Debug = false
	return &restCl
}
