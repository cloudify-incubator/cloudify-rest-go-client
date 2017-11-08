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
	"fmt"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	"net/url"
	"os"
	"path/filepath"
)

type Blueprint struct {
	// have id, owner information
	rest.CloudifyResource
	MainFileName string `json:"main_file_name"`
	// TODO describe "plan" struct
}

type BlueprintGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	Blueprint
}

type Blueprints struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []Blueprint           `json:"items"`
}

func (cl *Client) GetBlueprints(params map[string]string) (*Blueprints, error) {
	var blueprints Blueprints

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("blueprints?"+values.Encode(), &blueprints)
	if err != nil {
		return nil, err
	}

	return &blueprints, nil
}

func (cl *Client) DeleteBlueprints(blueprintID string) (*BlueprintGet, error) {
	var blueprint BlueprintGet

	err := cl.Delete("blueprints/"+blueprintID, &blueprint)
	if err != nil {
		return nil, err
	}

	return &blueprint, nil
}

func (cl *Client) DownloadBlueprints(blueprintID string) (string, error) {
	fileName := blueprintID + ".tar.gz"

	_, errFile := os.Stat(fileName)
	if !os.IsNotExist(errFile) {
		return "", fmt.Errorf("file `%s` is exist", fileName)
	}

	err := cl.GetBinary("blueprints/"+blueprintID+"/archive", fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (cl *Client) UploadBlueprint(blueprintID, path string) (*BlueprintGet, error) {

	absPath, errAbs := filepath.Abs(path)
	if errAbs != nil {
		return nil, errAbs
	}

	dirPath, nameFile := filepath.Split(absPath)

	var blueprint BlueprintGet

	err := cl.PutZip("blueprints/"+blueprintID+"?application_file_name="+nameFile, dirPath, &blueprint)
	if err != nil {
		return nil, err
	}

	return &blueprint, nil
}
