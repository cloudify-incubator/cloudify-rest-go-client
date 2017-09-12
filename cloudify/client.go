package cloudify

import (
	"encoding/json"
	rest "github.com/0lvin-cfy/cloudify-rest-go-client/cloudifyrest"
	utils "github.com/0lvin-cfy/cloudify-rest-go-client/cloudifyutils"
	"io/ioutil"
)

const ApiVersion = "v3.1"

type CloudifyClient struct {
	RestCl rest.CloudifyRestClient
}

func NewClient(host, user, password, tenant string) *CloudifyClient {
	var cliCl CloudifyClient
	if host[:len("https://")] == "https://" || host[:len("http://")] == "http://" {
		cliCl.RestCl.RestURL = host + "/api/" + ApiVersion + "/"
	} else {
		cliCl.RestCl.RestURL = "http://" + host + "/api/" + ApiVersion + "/"
	}
	cliCl.RestCl.User = user
	cliCl.RestCl.Password = password
	cliCl.RestCl.Tenant = tenant
	cliCl.RestCl.Debug = false
	return &cliCl
}

func (cl *CloudifyClient) EnableDebug() {
	cl.RestCl.Debug = true
}

func (cl *CloudifyClient) Get(url string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Get(url, rest.JsonContentType)

	err := json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func (cl *CloudifyClient) GetBinary(url, output_path string) error {
	body := cl.RestCl.Get(url, rest.DataContentType)

	err := ioutil.WriteFile(output_path, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func binaryPut(cl *CloudifyClient, url string, input []byte, input_type string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Put(url, input_type, input)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
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

	body := cl.RestCl.Post(url, json_data)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func (cl *CloudifyClient) Delete(url string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Delete(url)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}
