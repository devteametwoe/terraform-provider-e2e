package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

//production url  -> "https://api.e2enetworks.com/myaccount/api/v1/nodes/"

// groot url -> "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/"
func (c *Client) NewBlockStorage(item *models.BlockStorageCreate, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}

	UrlBlockStorage := c.Api_endpoint + "block_storage/"

	log.Printf("[INFO] CLIENT NEWBLOCKSTORAGE | BEFORE REQUEST")
	req, err := http.NewRequest("POST", UrlBlockStorage, &buf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

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

func (c *Client) GetBlockStorage(blockStorageID string, project_id int, location string) (map[string]interface{}, error) {
	urlBlockStorage := c.Api_endpoint + "block_storage/" + blockStorageID + "/"
	req, err := http.NewRequest("GET", urlBlockStorage, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | BLOCK STORAGE READ")
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT BLOCK STORAGE READ | after response %d", response.StatusCode)
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
		log.Printf("[ERROR] CLIENT GET BLOCK STORAGE | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) DeleteBlockStorage(blockStorageID string, project_id int, location string) error {

	urlNode := c.Api_endpoint + "block_storage/" + blockStorageID + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateBlockStorage(item *models.BlockStorageUpgrade, blockStorageID string, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	urlNode := c.Api_endpoint + "block_storage/" + blockStorageID + "/vm/upgrade/"
	req, err := http.NewRequest("PUT", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
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
func (c *Client) AttachOrDetachBlockStorage(item *models.BlockStorageAttach, Action string, blockStorageID string, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	urlNode := c.Api_endpoint + "block_storage/" + blockStorageID + "/vm/" + Action + "/"
	req, err := http.NewRequest("PUT", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	log.Printf("[INFO] CLIENT | %s BLOCK STORAGE | before request %+v", Action, req)
	response, err := c.HttpClient.Do(req)
	log.Printf("[INFO] CLIENT | %s BLOCK STORAGE | after response %+v", Action, response)
	if err == nil {
		err = CheckResponseStatusForBlock(response)
	}
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
func (c *Client) GetBlockStoragePlans(project_id int, location string) (map[string]interface{}, error) {
	urlBlockStorage := c.Api_endpoint + "block_storage/plans/"
	req, err := http.NewRequest("GET", urlBlockStorage, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	client := &http.Client{}
	log.Printf("[INFO] CLIENT | GET BLOCK STORAGE PLANS, BEFORE REQUEST %+v", req)
	response, err := client.Do(req)
	log.Printf("[INFO] CLIENT | GET BLOCK STORAGE PLANS, AFTER RESPONSE response %+v", response)
	if err != nil {
		return nil, err
	}
	err = CheckResponseStatusForBlock(response)
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
	log.Printf("[INFO] CLIENT | GET BLOCK STORAGE PLANS, AFTER RESPONSE jsonRes %+v", jsonRes)
	return jsonRes, nil
}

func addParamsAndHeaders(req *http.Request, Api_key string, Auth_token string, project_id int, location string) *http.Request {
	params := req.URL.Query()
	params.Add("apikey", Api_key)
	params.Add("project_id", strconv.Itoa(project_id))
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	return req
}

func CheckResponseStatusForBlock(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, errorResponse.Errors)
	}
	return nil
}
