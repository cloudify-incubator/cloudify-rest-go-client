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

type CloudifyBlueprint struct {
	// have id, owner information
	rest.CloudifyResource
	MainFileName string `json:"main_file_name"`
	// TODO describe "plan" struct
}

type CloudifyBlueprintGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	CloudifyBlueprint
}

type CloudifyBlueprints struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyBlueprint   `json:"items"`
}

func (cl *CloudifyClient) GetBlueprints(params map[string]string) (*CloudifyBlueprints, error) {
	var blueprints CloudifyBlueprints

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

func (cl *CloudifyClient) DeleteBlueprints(blueprintId string) (*CloudifyBlueprintGet, error) {
	var blueprint CloudifyBlueprintGet

	err := cl.Delete("blueprints/"+blueprintId, &blueprint)
	if err != nil {
		return nil, err
	}

	return &blueprint, nil
}

func (cl *CloudifyClient) DownloadBlueprints(blueprintId string) (string, error) {
	fileName := blueprintId + ".tar.gz"

	_, errFile := os.Stat(fileName)
	if !os.IsNotExist(errFile) {
		return "", fmt.Errorf("File `%s` is exist.", fileName)
	}

	err := cl.GetBinary("blueprints/"+blueprintId+"/archive", fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (cl *CloudifyClient) UploadBlueprint(blueprintId, path string) (*CloudifyBlueprintGet, error) {

	absPath, errAbs := filepath.Abs(path)
	if errAbs != nil {
		return nil, errAbs
	}

	dirPath, nameFile := filepath.Split(absPath)

	var blueprint CloudifyBlueprintGet

	err := cl.PutZip("blueprints/"+blueprintId+"?application_file_name="+nameFile, dirPath, &blueprint)
	if err != nil {
		return nil, err
	}

	return &blueprint, nil
}
