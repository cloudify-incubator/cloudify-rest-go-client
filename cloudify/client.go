package cloudify

import (
	"encoding/json"
	rest "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/rest"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"io/ioutil"
)

//Client - struct with connection settings for connect to manager
type Client struct {
	restCl rest.ConnectionOperationsInterface
}

//ClientFromConnection - return new client with internally use provided connection
func ClientFromConnection(conn rest.ConnectionOperationsInterface) *Client {
	var cliCl Client
	cliCl.restCl = conn
	return &cliCl
}

//NewClient - return new connection with params
func NewClient(host, user, password, tenant string) *Client {
	return ClientFromConnection(rest.NewClient(host, user, password, tenant))
}

//EnableDebug - Enable debug on current connection
func (cl *Client) EnableDebug() {
	cl.restCl.SetDebug(true)
}

//GetAPIVersion - return currently supported api version
func (cl *Client) GetAPIVersion() string {
	return rest.APIVersion
}

//Get - get cloudify object from server
func (cl *Client) Get(url string, output rest.MessageInterface) error {
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

//GetBinary - get binary object from manager without any kind of unmarshaling
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

//binaryPut - store/send object to manger without marshaling, response will be unmarshaled
func binaryPut(cl *Client, url string, input []byte, inputType string, output rest.MessageInterface) error {
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

//PutBinary - store/send binary object to manger without marshaling, response will be unmarshaled
func (cl *Client) PutBinary(url string, data []byte, output rest.MessageInterface) error {
	return binaryPut(cl, url, data, rest.DataContentType, output)
}

//PutZip - store/send path as archive to manger without marshaling, response will be unmarshaled
func (cl *Client) PutZip(url, path string, output rest.MessageInterface) error {
	data, err := utils.DirZipArchive(path)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, data, rest.DataContentType, output)
}

//Put - send object to manager(mainly replece old one)
func (cl *Client) Put(url string, input interface{}, output rest.MessageInterface) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, jsonData, rest.JSONContentType, output)
}

//Post - send cloudify object to manager
func (cl *Client) Post(url string, input interface{}, output rest.MessageInterface) error {
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

//Delete - delete cloudify object on manager
func (cl *Client) Delete(url string, output rest.MessageInterface) error {
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
