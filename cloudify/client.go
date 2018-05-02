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
	"encoding/json"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"io/ioutil"
	"net/url"
)

// ClientConfig - all configuration fields for connection
type ClientConfig struct {
	Host            string `json:"host,omitempty"`
	User            string `json:"user,omitempty"`
	Password        string `json:"password,omitempty"`
	Tenant          string `json:"tenant,omitempty"`
	AgentFile       string `json:"agent,omitempty"`
	DeploymentsFile string `json:"deployment,omitempty"`
	Debug           bool   `json:"debug,omitempty"`
}

//Client - struct with connection settings for connect to manager
type Client struct {
	ClientConfig
	restClCache rest.ConnectionOperationsInterface
}

//restCl - return client connection
func (cl *Client) restCl() rest.ConnectionOperationsInterface {
	var conn rest.ConnectionOperationsInterface
	if cl.restClCache != nil {
		conn = cl.restClCache
	} else {
		cl.updateHostFromAgent()
		conn = rest.NewClient(cl.Host, cl.User, cl.Password, cl.Tenant)
	}
	conn.SetDebug(cl.Debug)
	return conn
}

//CacheConnection - lock connection, don't reread agent file
func (cl *Client) CacheConnection() {
	// already cached
	if cl.restClCache != nil {
		return
	}
	// cache connection
	cl.restClCache = cl.restCl()
}

//ResetConnection - reset cached connection settings, need to recreate connection
// if you have used ClientFromConnection
func (cl *Client) ResetConnection() {
	cl.restClCache = nil
}

//ClientFromConnection - return new client with internally use provided connection
func ClientFromConnection(conn rest.ConnectionOperationsInterface) *Client {
	var cliCl Client
	cliCl.restClCache = conn
	return &cliCl
}

//NewClient - return new connection with params
func NewClient(cloudConfig ClientConfig) *Client {
	var cliCl Client
	cliCl.ClientConfig = cloudConfig
	return &cliCl
}

//EnableDebug - Enable debug on current connection
func (cl *Client) EnableDebug() {
	cl.Debug = true
}

//GetAPIVersion - return currently supported api version
func (cl *Client) GetAPIVersion() string {
	return rest.APIVersion
}

//Get - get cloudify object from server
func (cl *Client) Get(url string, output rest.MessageInterface) error {
	body, err := cl.restCl().Get(url, rest.JSONContentType)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

//GetBinary - get binary object from manager without any kind of unmarshaling
func (cl *Client) GetBinary(url, outputPath string) error {
	body, err := cl.restCl().Get(url, rest.DataContentType)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

//binarySend - store/send object to manger without marshaling, response will be unmarshaled
func binarySend(cl *Client, usePut bool, url string, input []byte, inputType string, output rest.MessageInterface) error {
	var body []byte
	var err error
	if usePut {
		body, err = cl.restCl().Put(url, inputType, input)
	} else {
		body, err = cl.restCl().Post(url, inputType, input)
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

//PutBinary - store/send binary object to manger without marshaling, response will be unmarshaled
func (cl *Client) PutBinary(url string, data []byte, output rest.MessageInterface) error {
	return binarySend(cl, true, url, data, rest.DataContentType, output)
}

//PutZip - store/send path as archive to manger without marshaling, response will be unmarshaled
func (cl *Client) PutZip(url string, paths []string, output rest.MessageInterface) error {
	data, err := utils.DirZipArchive(paths)
	if err != nil {
		return err
	}

	return binarySend(cl, true, url, data, rest.DataContentType, output)
}

//PostZip - store/send path as archive to manger without marshaling, response will be unmarshaled
func (cl *Client) PostZip(url string, paths []string, output rest.MessageInterface) error {
	data, err := utils.DirZipArchive(paths)
	if err != nil {
		return err
	}

	return binarySend(cl, false, url, data, rest.DataContentType, output)
}

//Put - send object to manager(mainly replece old one)
func (cl *Client) Put(url string, input interface{}, output rest.MessageInterface) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return binarySend(cl, true, url, jsonData, rest.JSONContentType, output)
}

//Post - send cloudify object to manager
func (cl *Client) Post(url string, input interface{}, output rest.MessageInterface) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	body, err := cl.restCl().Post(url, rest.JSONContentType, jsonData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

//Delete - delete cloudify object on manager
func (cl *Client) Delete(url string, input interface{}, output rest.MessageInterface) error {
	var jsonData = []byte{}
	var err error
	if input != nil {
		jsonData, err = json.Marshal(input)
		if err != nil {
			return err
		}
	}
	body, err := cl.restCl().Delete(url, rest.JSONContentType, jsonData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

// stringMapToURLValue - convert map[string]string -> url.Values
func (cl *Client) stringMapToURLValue(params map[string]string) url.Values {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return values
}
