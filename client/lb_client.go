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

func (c *Client) AddParamsAndHeader(req *http.Request, location string, project_id string) (*http.Request, error) {
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	return req, nil
}

func (c *Client) NewLoadBalancer(item *models.LoadBalancerCreate, project_id string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlEndPoint := c.Api_endpoint + "appliances/load-balancers/"
	log.Printf("[INFO] CLIENT NEWLOADBALANCER| BEFORE REQUEST")
	buf, err = RemoveExtraKeysLoadBalancer(&buf)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", UrlEndPoint, &buf)
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

func (c *Client) GetLoadBalancerInfo(lbId string, location string, project_id string) (map[string]interface{}, error) {
	urlLbInfo := c.Api_endpoint + "appliances/" + lbId + "/"
	req, err := http.NewRequest("GET", urlLbInfo, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] CLIENT | LOAD BALANCER READ")

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		return nil, err
	}

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT LOAD BALANCER READ | response code %d", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("======================INSIDE ERROR FROM READFROM======================")
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		log.Printf("=======================ERROR FROM API NOT GETTING 200 ====================")
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	log.Printf("======================NOW DEFER CLOSE RESPONSE BODY ===========================")
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	log.Printf("================STRING RESPONSE=========================%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET LOAD BALANCER | error when unmarshalling | %s", err)
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) DeleteLoadBalancer(lbId string, location string, project_id string) error {
	urlLbInfo := c.Api_endpoint + "appliances/" + lbId + "/"
	req, err := http.NewRequest("DELETE", urlLbInfo, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] CLIENT | LOAD BALANCER DELETE")

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		return err
	}

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

func (c *Client) UpdateLoadBalancerAction(data map[string]interface{}, lbId string, location string, project_id string) error {
	UrlEndPoint := c.Api_endpoint + "appliances/load-balancers/" + lbId + "/actions/"

	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ERROR] UpdateLoadBalancerAction | PAYLOAD_ERROR | %s", err)
		return err
	}

	req, err := http.NewRequest("PUT", UrlEndPoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("[ERROR] UpdateLoadBalancerAction | HTTP_NEW_REQUEST_ERROR | %s", err)
		return err
	}

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		log.Printf("[ERROR] UpdateLoadBalancerAction | ADDING_PARAMS_HEADER_ERROR | %s", err)
		return err
	}

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

func (c *Client) IPV6LoadBalancerAction(data map[string]interface{}, lbId string, location string, project_id string) error {
	UrlEndPoint := c.Api_endpoint + "appliances/load-balancers/" + lbId + "/ipv6/"

	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ERROR] IPV6LoadBalancerAction | PAYLOAD_ERROR | %s", err)
		return err
	}

	req, err := http.NewRequest("PUT", UrlEndPoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("[ERROR] IPV6LoadBalancerAction | HTTP_NEW_REQUEST_ERROR | %s", err)
		return err
	}

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		log.Printf("[ERROR] IPV6LoadBalancerAction | ADDING_PARAMS_HEADER_ERROR | %s", err)
		return err
	}

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

func (c *Client) LoadBalancerBackendUpdate(item *models.LoadBalancerCreate, lbId string, location string, project_id string) (map[string]interface{}, error) {
	UrlEndPoint := c.Api_endpoint + "appliances/load-balancers/" + lbId + "/"

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | JSON_NEW_ENCODER_ERROR | %s", err)
		return nil, err
	}

	buf, err = RemoveExtraKeysLoadBalancer(&buf)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | RemoveExtraKeysLoadBalancer | %s", err)
		return nil, err
	}

	log.Printf("================LOAD BALANCER UPDATE API INFO==================, %s, %s", UrlEndPoint, &buf)
	req, err := http.NewRequest("PUT", UrlEndPoint, &buf)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | NEW_REQUEST_ERROR | %s", err)
		return nil, err
	}

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | ADDING_PARAMS_HEADER_ERROR | %s", err)
		return nil, err
	}

	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | ERROR_WHILE_EXECUTING_REQUEST | %s", err)
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | CHECK_RESPONSE_STATUS | %s", err)
		return nil, err
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] LoadBalancerBackendUpdate | UNMARSHAL_RESPONSE | %s", err)
		return nil, err
	}
	log.Printf("[INFO] LoadBalancerBackendUpdate | LOAD BALANCER API UPDATE CALL SUCCESS | RESPONSE | %s", jsonRes)
	return jsonRes, nil
}
