package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) AddSshKey(item models.AddSshKey, project_id string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlSshKey := c.Api_endpoint + "ssh_keys/"
	log.Printf("[INFO] %s", UrlSshKey)
	req, err := http.NewRequest("POST", UrlSshKey, &buf)
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
	log.Printf("inside add ssh key req = %+v, res = %+v", req, response)
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

func (c *Client) GetSshKey(label string, project_id string) (map[string]interface{}, error) {
	UrlSshKey := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", UrlSshKey, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("label", label)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// err = CheckResponseStatus(response)
	// if err != nil {
	// 	return nil, err
	// }

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", resBody)
	stringresponse := string(resBody)
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteSshKey(pk string, project_id string, location string) error {
	UrlSshKey := c.Api_endpoint + "delete_ssh_key/" + pk + "/"
	req, err := http.NewRequest("DELETE", UrlSshKey, nil)
	if err != nil {
		return err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	// err = CheckResponseStatus(response)
	// if err != nil {
	// 	return err
	// }
	defer response.Body.Close()
	return nil
}

func (c *Client) GetSshKeys(location string, project_id string) (*models.SshKeyResponse, error) {

	urlSshKeys := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", urlSshKeys, nil)
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
