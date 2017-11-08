package cloudify

import (
	"encoding/json"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"io/ioutil"
)

type Client struct {
	restCl *rest.CloudifyRestClient
}

func NewClient(host, user, password, tenant string) *Client {
	var cliCl Client
	cliCl.restCl = rest.NewClient(host, user, password, tenant)
	return &cliCl
}

func (cl *Client) EnableDebug() {
	cl.restCl.Debug = true
}

func (cl *Client) GetAPIVersion() string {
	return rest.APIVersion
}

func (cl *Client) Get(url string, output rest.CloudifyMessageInterface) error {
	body, err := cl.restCl.Get(url, rest.JSONContentType)
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

func (cl *Client) GetBinary(url, outputPath string) error {
	body, err := cl.restCl.Get(url, rest.DataContentType)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func binaryPut(cl *Client, url string, input []byte, inputType string, output rest.CloudifyMessageInterface) error {
	body, err := cl.restCl.Put(url, inputType, input)
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

func (cl *Client) PutBinary(url string, data []byte, output rest.CloudifyMessageInterface) error {
	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *Client) PutZip(url, path string, output rest.CloudifyMessageInterface) error {
	data, err := utils.DirZipArchive(path)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *Client) Put(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, jsonData, rest.JSONContentType, output)
}

func (cl *Client) Post(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	body, err := cl.restCl.Post(url, jsonData)
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

func (cl *Client) Delete(url string, output rest.CloudifyMessageInterface) error {
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
