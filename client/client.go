package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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

func (c *Client) NewNode(item *models.NodeCreate, project_id string, location string) (map[string]interface{}, error) {
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
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	log.Printf("inside new Nodes req = %+v", req)
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
	params.Add("contact_person_id", "null")
	params.Add("project_id", project_id)
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

func (c *Client) UpdateNode(nodeId string, action string, Name string, project_id string, location string) (interface{}, error) {

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
	params.Add("location", location)
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	req.URL.RawQuery = params.Encode()
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] INSIDE NODE UPDATE %s %+v %+v", action, req, response)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		log.Printf("[INFO] INSIDE NODE UPDATE WRONG_STATUS %s %+v", action, respBody.String())
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

func (c *Client) UpdateNodeSSH(nodeId string, action string, ssh_keys []interface{}, project_id string, location string) (interface{}, error) {

	ssh_keys_map := generateSSHKeyMap(ssh_keys)
	if len(ssh_keys_map) == 0 {
		ssh_keys_map = make([]map[string]interface{}, 0)
	}
	log.Printf("[INFO] inside update ssh | ssh_keys_map = %+v", ssh_keys_map)
	node_action := models.NodeActionSSH{
		Type:     action,
		SSH_KEYS: ssh_keys_map,
	}
	nodeAction, _ := json.Marshal(node_action)
	url := c.Api_endpoint + "nodes/" + nodeId + "/actions/"
	log.Printf("[info] %s", url)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(nodeAction))
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	req.URL.RawQuery = params.Encode()
	log.Printf("[INFO] inside update ssh req = %+v", req)
	// return nil, err
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
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, err
}
func (c *Client) UpgradeNodePlan(nodeId string, plan string, image string, project_id string, location string) (interface{}, error) {
	node_action := models.NodePlanUpgradeAction{
		Plan:  plan,
		Image: image,
	}
	nodeAction, _ := json.Marshal(node_action)

	url := c.Api_endpoint + "nodes/upgrade/" + nodeId
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(nodeAction))
	if err != nil {
		log.Printf("[INFO] error inside upgrade node plan")
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	log.Printf("CLIENT UPGRADE NODE PLAN | request = %+v", req)
	log.Printf("CLIENT UPGRADE NODE PLAN | STATUS_CODE: %d, response = %+v", response.StatusCode, response)
	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[INFO] error inside upgrade node plan")
		return nil, err
	}
	return response, err
}

func (c *Client) DeleteNode(nodeId string, project_id string, location string) error {

	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	log.Printf("inside delete node req.URL = %s", req.URL)
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
func (c *Client) GetSecurityGroups(project_id int, location string) (*models.SecurityGroupsResponse, error) {

	urlSecurityGroups := c.Api_endpoint + "security_group/"
	req, err := http.NewRequest("GET", urlSecurityGroups, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
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

func (c *Client) GetVpcs(location string, project_id string) (*models.VpcsResponse, error) {

	urlGetVpcs := c.Api_endpoint + "vpc/" + "list/"
	req, err := http.NewRequest("GET", urlGetVpcs, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
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
func (c *Client) GetVpc(vpc_id string, project_id string, location string) (*models.VpcResponse, error) {

	urlGetVpc := c.Api_endpoint + "vpc/" + vpc_id + "/"
	req, err := http.NewRequest("GET", urlGetVpc, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
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

func (c *Client) CreateVpc(location string, item *models.VpcCreate, project_id string) (map[string]interface{}, error) {

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
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)

	log.Printf("inside create vpc req = %+v, res = %+v, Error = %+v", req, response, err)
	if err != nil {
		return nil, err
	}
	err = CheckResponseCreatedStatus(response)
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

func (c *Client) DeleteVpc(vpcId string, project_id string, location string) (map[string]interface{}, error) {

	urlVpc := c.Api_endpoint + "vpc/" + vpcId + "/"
	log.Printf("[INFO] %s", urlVpc)
	req, err := http.NewRequest("DELETE", urlVpc, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
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
func (c *Client) NewReservedIp(project_id string, location string) (map[string]interface{}, error) {

	UrlReservedIp := c.Api_endpoint + "reserve_ips/"
	log.Printf("[INFO] Url = %s", UrlReservedIp)
	req, err := http.NewRequest("POST", UrlReservedIp, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	log.Printf("\n\n[INFO] CLIENT NEW RESERVED IP | STATUS_CODE: %+v ==================***************\n\n", response)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unauthorized | status %v | The provided api_token or api_key or project_id seem to be incorrect. Please revise them accordingly", response.StatusCode)
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

func (c *Client) DeleteReserveIP(ip_address string, project_id string, location string) error {
	urlNode := c.Api_endpoint + "reserve_ips/" + ip_address + "/actions/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		log.Printf("[INFO] error inside delete reserve ip")
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside delete reserve ip")
		return err
	}
	log.Printf("CLIENT DELETE NODE | STATUS_CODE: %d", response.StatusCode)
	return nil

}

func (c *Client) GetReservedIp(ip_address string, project_id string, location string) (*models.ResponseReserveIps, error) {

	urlNode := c.Api_endpoint + "reserve_ips/" + ip_address + "/actions/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get reserve ip")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseReserveIps{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get reserve ip | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetReservedIps(project_id string, location string) (*models.ResponseReserveIps, error) {

	urlGetReserveIps := c.Api_endpoint + "reserve_ips/"
	req, err := http.NewRequest("GET", urlGetReserveIps, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[INFO] error inside GetReservedIps")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseReserveIps{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside GetReservedIps | error while unmarshlling")
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

func (c *Client) DetachSecurityGroup(item *models.UpdateSecurityGroups, vm_id int, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUP DETACH | BEFORE REQUEST")
	if err != nil {
		return nil, err
	}
	vmIDInString := strconv.Itoa(vm_id)
	urlNode := c.Api_endpoint + "security_group/" + vmIDInString + "/detach/"
	req, err := http.NewRequest("POST", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	projectIDInt, error := strconv.Atoi(project_id)
	if error != nil {
		return nil, error
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
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

func (c *Client) AttachSecurityGroup(item *models.UpdateSecurityGroups, vm_id int, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUP ATTACH | BEFORE REQUEST")
	if err != nil {
		return nil, err
	}
	vmIDInString := strconv.Itoa(vm_id)
	urlNode := c.Api_endpoint + "security_group/" + vmIDInString + "/attach/"
	req, err := http.NewRequest("POST", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	projectIDInt, error := strconv.Atoi(project_id)
	if error != nil {
		return nil, error
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
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

func (c *Client) GetSecurityGroupList(project_id string, location string) (map[string]interface{}, error) {

	urlSecurityGroups := c.Api_endpoint + "security_group/"
	req, err := http.NewRequest("GET", urlSecurityGroups, nil)
	if err != nil {
		return nil, err
	}
	projectIDInt, err := strconv.Atoi(project_id)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get security groups list")
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUPS LIST READ | response code %d", response.StatusCode)
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
		log.Printf("[ERROR] CLIENT GET SECURITY GROUP | error when unmarshalling | %s", err)
		return nil, err
	}
	return jsonRes, nil
}

func CheckResponseCreatedStatus(response *http.Response) error {
	if response.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 201 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 201 status code: %v - %s", response.StatusCode, respBody.String())
	}
	return nil
}

func (c *Client) CheckNodeLCMState(nodeId string, project_id string, location string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/check-lcm-state/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[error]  CLIENT CheckNodeLCMState |  error = %v", err)
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	log.Printf("[info] CLIENT | CheckNodeLCMState | res = %+v", res)
	if err != nil {
		log.Printf("[ERROR] CLIENT  | CheckNodeLCMState | ERROR WHILE UNMARSHALLING")
		return nil, err
	}
	return res, nil
}
