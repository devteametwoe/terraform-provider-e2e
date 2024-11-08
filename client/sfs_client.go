package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client)NewSfs(item *models.SfsCreate, project_id string, location string)(map[string]interface{}, error){
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlSfs := c.Api_endpoint + "efs/"+ "create/"
	log.Printf("[INFO] Client NEWNODE | BEFORE REQUEST")
	req, err := http.NewRequest("POST", UrlSfs, &buf)
	if err != nil {
		return nil, err
	}
	AddParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

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
func (c *Client) GetSfs(SfsId string , project_id string, location string) (map[string]interface{}, error) {

	UrlSfs := c.Api_endpoint + "efs/" + SfsId + "/"
	req, err := http.NewRequest("GET", UrlSfs, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Client | NODE READ")
	AddParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)


	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Client NODE READ | after response %d", response.StatusCode)
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
		log.Printf("[ERROR] Client GET NDE | error when unmarshalling")
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteSFs(SfsId string, project_id string , location string) error {

	UrlSfs := c.Api_endpoint + "efs/" + "delete/"+ SfsId + "/"
	req, err := http.NewRequest("DELETE", UrlSfs, nil)
	if err != nil {
		return err
	}

	AddParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

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


func (c *Client) GetSfss(location string, project_id string) (*models.ResponseSfss, error) {

	urlGetSfsss := c.Api_endpoint + "efs/"
	req, err := http.NewRequest("GET", urlGetSfsss, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT GET NODES sfs_list")
	AddParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

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
	
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseSfss{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func AddParamsAndHeaders(req *http.Request, Api_key string, Auth_token string, project_id string, location string) *http.Request {
	params := req.URL.Query()
	params.Add("apikey", Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	return req
}