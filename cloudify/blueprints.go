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
	"log"
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

func (cl *CloudifyClient) GetBlueprints(params map[string]string) CloudifyBlueprints {
	var blueprints CloudifyBlueprints

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	err := cl.Get("blueprints?"+values.Encode(), &blueprints)
	if err != nil {
		log.Fatal(err)
	}

	return blueprints
}

func (cl *CloudifyClient) DeleteBlueprints(blueprint_id string) CloudifyBlueprintGet {
	var blueprint CloudifyBlueprintGet

	err := cl.Delete("blueprints/"+blueprint_id, &blueprint)
	if err != nil {
		log.Fatal(err)
	}

	return blueprint
}

func (cl *CloudifyClient) DownloadBlueprints(blueprint_id string) string {
	file_name := blueprint_id + ".tar.gz"

	_, err_file := os.Stat(file_name)
	if !os.IsNotExist(err_file) {
		log.Fatal(fmt.Sprintf("File `%s` is exist.", file_name))
	}

	err := cl.GetBinary("blueprints/"+blueprint_id+"/archive", file_name)
	if err != nil {
		log.Fatal(err)
	}

	return file_name
}

func (cl *CloudifyClient) UploadBlueprint(blueprint_id, path string) CloudifyBlueprintGet {

	absPath, err_abs := filepath.Abs(path)
	if err_abs != nil {
		log.Fatal(err_abs)
	}

	dirPath, name_file := filepath.Split(absPath)

	var blueprint CloudifyBlueprintGet

	err := cl.PutZip("blueprints/"+blueprint_id+"?application_file_name="+name_file, dirPath, &blueprint)
	if err != nil {
		log.Fatal(err)
	}

	return blueprint
}
