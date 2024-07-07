package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	auth          string
	orgServiceUrl string
	projectId     string
}

func New(pat, orgServiceUrl, projectId string) *Client {
	return &Client{
		auth:          "Basic " + base64.StdEncoding.EncodeToString([]byte("pat"+":"+pat)),
		orgServiceUrl: orgServiceUrl,
		projectId:     projectId,
	}
}

func (c *Client) organizationBaseUrl() string {
	return fmt.Sprintf("%s/_apis", strings.TrimRight(c.orgServiceUrl, "/"))
}

func (c *Client) projectBaseUrl() string {
	return fmt.Sprintf("%s/%s/_apis", strings.TrimRight(c.orgServiceUrl, "/"), c.projectId)
}

func (c *Client) newJsonRequest(method string, body any, url string) (*http.Request, error) {
	var request *http.Request
	var err error
	if body != nil {
		buffer, err := packJSON(body)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(method, url, buffer)
		request.Header.Add("Content-Type", "application/json")
		if err != nil {
			return nil, err
		}
	} else {
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	request.Header.Add("Authorization", c.auth)
	return request, nil
}

func unpackJSON(obj io.Reader, container any) error {
	bodyBytes, err := io.ReadAll(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, container)
	if err != nil {
		return err
	}
	return nil
}

func packJSON(obj any) (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}
