package sumologic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	AccessID    string
	AccessKey   string
	Environment string
	BaseURL     *url.URL
	httpClient  HttpClient
}

var endpoints = map[string]string{
	"us1":  "https://api.sumologic.com/api/v1/",
	"us2":  "https://api.us2.sumologic.com/api/v1/",
	"eu":   "https://api.eu.sumologic.com/api/v1/",
	"au":   "https://api.au.sumologic.com/api/v1/",
	"de":   "https://api.de.sumologic.com/api/v1/",
	"jp":   "https://api.jp.sumologic.com/api/v1/",
	"ca":   "https://api.ca.sumologic.com/api/v1/",
	"nite": "https://nite-api.sumologic.net/api/",
}

var rateLimiter = time.Tick(time.Minute / 240)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (s *Client) PostWithCookies(urlPath string, payload interface{}) ([]byte, []*http.Cookie, error) {
	relativeURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, nil, err
	}

	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, sumoURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(s.AccessID, s.AccessKey)

	<-rateLimiter
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respCookie := resp.Cookies()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode >= 400 {
		errMsgFormat := fmt.Sprintf("Status code %d - %%s", resp.StatusCode)
		log.Printf("[DEBUG] client PostWithCookies returned status %d\n%s\n", resp.StatusCode, d)
		requestBody, err := json.MarshalIndent(payload, "", "  ")
		log.Printf("[TRACE] request body \n%s\n", requestBody)
		var errorResponse ErrorResponse
		err = json.Unmarshal(d, &errorResponse)
		if err != nil {
			log.Printf("[DEBUG] client PostWithCookies response unmarshalling failed %s\n%s\n", err, d)
			return nil, nil, errors.New(fmt.Sprintf(errMsgFormat, err))
		}
		return nil, nil, errors.New(fmt.Sprintf(errMsgFormat, errorResponse.Message))
	}

	return d, respCookie, nil
}

func (s *Client) GetWithCookies(urlPath string, cookies []*http.Cookie) ([]byte, string, error) {
	relativeURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, "", err
	}

	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	req, err := http.NewRequest(http.MethodGet, sumoURL.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(s.AccessID, s.AccessKey)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	<-rateLimiter
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode == 404 {
		return nil, "", nil
	} else if resp.StatusCode >= 400 {
		errMsgFormat := fmt.Sprintf("Status code %d - %%s", resp.StatusCode)
		log.Printf("[DEBUG] client GetWithCookies returned status %d\n%s\n", resp.StatusCode, d)
		var errorResponse ErrorResponse
		err = json.Unmarshal(d, &errorResponse)
		if err != nil {
			log.Printf("[DEBUG] client GetWithCookies response unmarshalling failed %s\n%s\n", err, d)
			return nil, "", errors.New(fmt.Sprintf(errMsgFormat, err))
		}
		return nil, "", errors.New(fmt.Sprintf(errMsgFormat, errorResponse.Message))
	}

	return d, resp.Header.Get("ETag"), nil
}

func (s *Client) Post(urlPath string, payload interface{}) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)
	log.Printf("[DEBUG] THOMASKAOPAYLOAD: %v\n", payload)
	body, _ := json.Marshal(payload)
	log.Printf("[DEBUG] THOMASKAOBody: %s\n", body)
	req, _ := http.NewRequest(http.MethodPost, sumoURL.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(s.AccessID, s.AccessKey)

	<-rateLimiter
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	d, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		errMsgFormat := fmt.Sprintf("Status code %d - %%s", resp.StatusCode)
		log.Printf("[DEBUG] client Post returned status %d\n%s\n", resp.StatusCode, d)
		requestBody, err := json.MarshalIndent(payload, "", "  ")
		log.Printf("[TRACE] request body \n%s\n", requestBody)
		var errorResponse ErrorResponse
		err = json.Unmarshal(d, &errorResponse)
		if err != nil {
			log.Printf("[DEBUG] client Post response unmarshalling failed %s\n%s\n", err, d)
			return nil, errors.New(fmt.Sprintf(errMsgFormat, err))
		}
		return nil, errors.New(fmt.Sprintf(errMsgFormat, errorResponse.Message))
	}

	return d, nil
}

func (s *Client) Put(urlPath string, payload interface{}) ([]byte, error) {
	log.Printf("[DEBUG] SSAIN: CLIENT URLPATH: " + urlPath)
	SumoMutexKV.Lock(urlPath)
	defer SumoMutexKV.Unlock(urlPath)

	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)
	log.Printf("[DEBUG] client Put path %s", sumoURL.String())
	_, etag, _ := s.Get(sumoURL.String())

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, sumoURL.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Match", etag)

	req.SetBasicAuth(s.AccessID, s.AccessKey)

	<-rateLimiter
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	d, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		errMsgFormat := fmt.Sprintf("Status code %d - %%s", resp.StatusCode)
		log.Printf("[DEBUG] client Put returned status %d\n%s\n", resp.StatusCode, d)
		requestBody, err := json.MarshalIndent(payload, "", "  ")
		log.Printf("[TRACE] request body \n%s\n", requestBody)
		var errorResponse ErrorResponse
		err = json.Unmarshal(d, &errorResponse)
		if err != nil {
			log.Printf("[DEBUG] client Put response unmarshalling failed %s\n%s\n", err, d)
			return nil, errors.New(fmt.Sprintf(errMsgFormat, err))
		}
		return nil, errors.New(fmt.Sprintf(errMsgFormat, errorResponse.Message))
	}

	return d, nil
}

func (s *Client) Get(urlPath string) ([]byte, string, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)
	log.Printf("[DEBUG] client Get path %s", sumoURL.String())
	req, _ := http.NewRequest(http.MethodGet, sumoURL.String(), nil)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(s.AccessID, s.AccessKey)

	<-rateLimiter
	resp, _ := s.httpClient.Do(req)

	d, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 404 {
		return nil, "", nil
	} else if resp.StatusCode >= 400 {
		errMsgFormat := fmt.Sprintf("Status code %d - %%s", resp.StatusCode)
		log.Printf("[DEBUG] client Get returned status %d\n%s\n", resp.StatusCode, d)
		var errorResponse ErrorResponse
		err := json.Unmarshal(d, &errorResponse)
		if err != nil {
			log.Printf("[DEBUG] client Get response unmarshalling failed %s\n%s\n", err, d)
			return nil, "", errors.New(fmt.Sprintf(errMsgFormat, err))
		}
		return nil, "", errors.New(fmt.Sprintf(errMsgFormat, errorResponse.Message))
	}

	return d, resp.Header.Get("ETag"), nil
}

func (s *Client) Delete(urlPath string) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	req, _ := http.NewRequest(http.MethodDelete, sumoURL.String(), nil)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(s.AccessID, s.AccessKey)

	<-rateLimiter
	resp, _ := s.httpClient.Do(req)

	d, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(d))
	}

	return d, nil
}

func NewClient(accessID, accessKey, environment string) (*Client, error) {
	client := Client{
		AccessID:    accessID,
		AccessKey:   accessKey,
		Environment: environment,
		httpClient:  http.DefaultClient,
	}

	client.BaseURL, _ = url.Parse(endpoints[client.Environment])
	return &client, nil
}
