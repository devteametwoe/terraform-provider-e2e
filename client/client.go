package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

//production url  -> "https://api.e2enetworks.com/myaccount/api/v1/nodes/"

//groot url -> "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/"

type Client struct {
	Api_key      string
	Auth_token   string
	Api_endpoint string
	HttpClient   *http.Client
}

func NewClient(api_key string, auth_token string, api_endpoint string) *Client {
	return &Client{
		Api_key:      api_key,
		Auth_token:   auth_token,
		Api_endpoint: api_endpoint,
		HttpClient:   &http.Client{},
	}
}

func (c *Client) NewNode(item *models.NodeCreate, project_id string) (map[string]interface{}, error) {

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlNode := c.Api_endpoint + "nodes/"
	log.Printf("[INFO] CLIENT NEWNODE | BEFORE REQUEST")
	req, err := http.NewRequest("POST", UrlNode, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil

}

func (c *Client) GetNode(nodeId string, project_id string) (map[string]interface{}, error) {

	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | NODE READ")
	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT NODE READ | after response %d", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		log.Printf("[ERROR] CLIENT GET NDE | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) GetNodes(location string, project_id string) (*models.ResponseNodes, error) {

	urlGetNodes := c.Api_endpoint + "nodes/"
	req, err := http.NewRequest("GET", urlGetNodes, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT GET NODES")
	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("User-Agent", "terraform-e2e")

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("GET NODES | INSIDE NO SUCCESS AND ERROR MSG")
			return nil, fmt.Errorf("%v", err)
		}

		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	fmt.Println(response.Body)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseNodes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateNode(nodeId string, action string, Name string, project_id string) (interface{}, error) {

	node_action := models.NodeAction{
		Type: action,
		Name: Name,
	}
	nodeAction, err := json.Marshal(node_action)
	url := c.Api_endpoint + "nodes/" + nodeId + "/actions/"
	log.Printf("[info] %s", url)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(nodeAction))
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	req.URL.RawQuery = params.Encode()
	response, err := c.HttpClient.Do(req)

	if err != nil {

		return nil, err
	}
	log.Printf("[INFO] inside update %s %d", action, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, err
}

func (c *Client) DeleteNode(nodeId string, project_id string) error {

	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}

	return nil
}

func (c *Client) GetSavedImages(location string, project_id string) (*models.ImageListResponse, error) {

	urlImages := c.Api_endpoint + "images/" + "saved-images" + "/"

	req, err := http.NewRequest("GET", urlImages, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	log.Printf("[INFO] inside client saved image before request hit")
	if err != nil {
		log.Printf("[INFO] error inside get image")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ImageListResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (c *Client) GetSecurityGroups() (*models.SecurityGroupsResponse, error) {

	urlSecurityGroups := c.Api_endpoint + "security_group/"
	req, err := http.NewRequest("GET", urlSecurityGroups, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get security groups")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.SecurityGroupsResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get security groups | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetSshKeys() (*models.SshKeyResponse, error) {

	urlSshKeys := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", urlSshKeys, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get ssh keys")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.SshKeyResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetVpcs(location string) (*models.VpcsResponse, error) {

	urlGetVpcs := c.Api_endpoint + "vpc/" + "list/"
	req, err := http.NewRequest("GET", urlGetVpcs, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	res := models.VpcsResponse{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get vpcs | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}
func (c *Client) GetVpc(vpc_id string) (*models.VpcResponse, error) {

	urlGetVpc := c.Api_endpoint + "vpc/" + vpc_id + "/"
	req, err := http.NewRequest("GET", urlGetVpc, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[INFO] client |  error inside get vpc")
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	res := models.VpcResponse{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get vpcs | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateVpc(location string, item *models.VpcCreate) (map[string]interface{}, error) {

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlNode := c.Api_endpoint + "vpc/"
	log.Printf("[INFO] %s", UrlNode)
	req, err := http.NewRequest("POST", UrlNode, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteVpc(vpcId string) (map[string]interface{}, error) {

	urlVpc := c.Api_endpoint + "vpc/" + vpcId + "/"
	log.Printf("[INFO] %s", urlVpc)
	req, err := http.NewRequest("DELETE", urlVpc, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetReservedIps(location string) (*models.ResponseReserveIps, error) {

	urlGetReserveIps := c.Api_endpoint + "reserve_ips/"
	req, err := http.NewRequest("GET", urlGetReserveIps, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[INFO] error inside get vpcs")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseReserveIps{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get vpcs | error while unmarshlling")
		return nil, err
	}
	return &res, nil

}

func (c *Client) GetImage(imageId string, project_id string) (*models.ImageResponse, error) {
	urlGetImage := c.Api_endpoint + "images/" + imageId + "/"
	req, err := http.NewRequest("GET", urlGetImage, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[error]  CLIENT READ IMAGE |  error inside get image")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ImageResponse{}
	err = json.Unmarshal(body, &res)
	log.Printf("[info] CLIENT | GET IMAGE |  %+v", res)
	if err != nil {
		log.Printf("[ERROR] CLIENT  | GET IMAGE | ERROR WHILE UNMARSHALLING")
		return nil, err
	}
	return &res, nil

}
func (c *Client) DeleteImage(imageId string, project_id string) error {
	urlNode := c.Api_endpoint + "images/" + imageId + "/"
	deleteBody := models.ImageDeleteBody{
		Action_type: "delete_image",
	}
	deleteBodyMarshalled, err := json.Marshal(deleteBody)

	req, err := http.NewRequest("DELETE", urlNode, bytes.NewBuffer(deleteBodyMarshalled))
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	log.Printf("CLIENT DELETE IMAGE | STATUS_CODE: %d", response.StatusCode)
	err = CheckResponseStatus(response)
	if err != nil {
		return err
	}
	return nil
}

func SetBasicHeaders(authtoken string, req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+authtoken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
}
func CheckResponseStatus(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	return nil
}
