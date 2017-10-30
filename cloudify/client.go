package cloudify

import (
	"encoding/json"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"io/ioutil"
)

type CloudifyClient struct {
	restCl *rest.CloudifyRestClient
}

func NewClient(host, user, password, tenant string) *CloudifyClient {
	var cliCl CloudifyClient
	cliCl.restCl = rest.NewClient(host, user, password, tenant)
	return &cliCl
}

func (cl *CloudifyClient) EnableDebug() {
	cl.restCl.Debug = true
}

func (cl *CloudifyClient) GetApiVersion() string {
	return rest.ApiVersion
}

func (cl *CloudifyClient) Get(url string, output rest.CloudifyMessageInterface) error {
	body, err := cl.restCl.Get(url, rest.JsonContentType)
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

func (cl *CloudifyClient) GetBinary(url, output_path string) error {
	body, err := cl.restCl.Get(url, rest.DataContentType)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output_path, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func binaryPut(cl *CloudifyClient, url string, input []byte, input_type string, output rest.CloudifyMessageInterface) error {
	body, err := cl.restCl.Put(url, input_type, input)
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

func (cl *CloudifyClient) PutBinary(url string, data []byte, output rest.CloudifyMessageInterface) error {
	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *CloudifyClient) PutZip(url, path string, output rest.CloudifyMessageInterface) error {
	data, err := utils.DirZipArchive(path)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *CloudifyClient) Put(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	json_data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, json_data, rest.JsonContentType, output)
}

func (cl *CloudifyClient) Post(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	json_data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	body, err := cl.restCl.Post(url, json_data)
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

func (cl *CloudifyClient) Delete(url string, output rest.CloudifyMessageInterface) error {
	body, err := cl.restCl.Delete(url)
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
