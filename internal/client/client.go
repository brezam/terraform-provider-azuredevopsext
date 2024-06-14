package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"terraform-provider-azuredevopsext/util"
)

const (
	environmentSecurityEndpoint = "securityroles/scopes/distributedtask.environmentreferencerole/roleassignments/resources"
)

type Client struct {
	auth          string
	orgServiceUrl string
	projectId     string
}

func NewClient(pat, orgServiceUrl, projectId string) *Client {
	return &Client{
		auth:          "Basic " + base64.StdEncoding.EncodeToString([]byte("pat"+":"+pat)),
		orgServiceUrl: orgServiceUrl,
		projectId:     projectId,
	}
}

func (c *Client) GetEnvironmentSecurityMembers(environmentId string) ([]EnvironmentSecurityAccess, error) {
	url := c.organizationUrl(environmentSecurityEndpoint, fmt.Sprintf("%s_%s?api-version=7.1-preview.1", c.projectId, environmentId))
	req, err := c.newJsonRequest("GET", nil, url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var valueList ValueList[EnvironmentSecurityAccess]
	defer resp.Body.Close()
	util.UnpackJSON(resp.Body, &valueList)
	return valueList.Value, nil
}

func (c *Client) AddMemberToEnvironmentSecurity(environmentId, memberId string, roleName RoleName) (*EnvironmentSecurityAccess, error) {
	body := []CreateEnvironmentSecurityAccess{{
		RoleName: roleName,
		UserId:   memberId,
	}}
	url := c.organizationUrl(environmentSecurityEndpoint, fmt.Sprintf("%s_%s?api-version=7.1-preview.1", c.projectId, environmentId))
	req, err := c.newJsonRequest("PUT", body, url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		stringBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create environment security: response: %s, body: %s", resp.Status, stringBody)
	}
	var valueList ValueList[EnvironmentSecurityAccess]
	util.UnpackJSON(resp.Body, &valueList)

	if len(valueList.Value) == 0 {
		return nil, fmt.Errorf("failed to create environment security: member Id '%s' is possibly wrong", memberId)
	}
	return &valueList.Value[0], nil
}

func (c *Client) DeleteMemberInEnvironmentSecurity(environmentId, memberId string) error {
	body := []string{memberId}
	url := c.organizationUrl(environmentSecurityEndpoint, fmt.Sprintf("%s_%s?api-version=7.1-preview.1", c.projectId, environmentId))
	req, err := c.newJsonRequest("PATCH", body, url)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	}
	return fmt.Errorf("failed to delete environment security: response: %s", resp.Status)
}

// private
func (c *Client) organizationUrl(endpoints ...string) string {
	url := fmt.Sprintf("%s/_apis", strings.TrimLeft(c.orgServiceUrl, "/"))
	for _, endpoint := range endpoints {
		url += "/" + strings.TrimLeft(endpoint, "/")
	}
	return url
}

func (c *Client) newJsonRequest(method string, body any, url string) (*http.Request, error) {
	var request *http.Request
	var err error
	if body != nil {
		buffer, err := util.PackJSON(body)
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
