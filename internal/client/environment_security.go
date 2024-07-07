package client

import (
	"fmt"
	"io"
	"net/http"
)

const environmentSecurityEndpoint = "%s/securityroles/scopes/distributedtask.environmentreferencerole/roleassignments/resources/%s_%s?api-version=%s"

func (c *Client) GetEnvironmentSecurityMembers(environmentId string) ([]EnvironmentSecurityAccess, error) {
	url := fmt.Sprintf(environmentSecurityEndpoint, c.organizationBaseUrl(), c.projectId, environmentId, APIVersion7_1Preview1)
	req, err := c.newJsonRequest("GET", nil, url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var valueList EnvironmentSecurityAccessList
	defer resp.Body.Close()
	unpackJSON(resp.Body, &valueList)
	return valueList.Value, nil
}

func (c *Client) AddMemberToEnvironmentSecurity(environmentId, memberId string, roleName RoleName) (*EnvironmentSecurityAccess, error) {
	body := []CreateEnvironmentSecurityAccess{{
		RoleName: roleName,
		UserId:   memberId,
	}}
	url := fmt.Sprintf(environmentSecurityEndpoint, c.organizationBaseUrl(), c.projectId, environmentId, APIVersion7_1Preview1)
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
	var valueList EnvironmentSecurityAccessList
	unpackJSON(resp.Body, &valueList)

	if len(valueList.Value) == 0 {
		return nil, fmt.Errorf("failed to create environment security: member Id '%s' is possibly wrong", memberId)
	}
	return &valueList.Value[0], nil
}

func (c *Client) DeleteMemberInEnvironmentSecurity(environmentId, memberId string) error {
	body := []string{memberId}
	url := fmt.Sprintf(environmentSecurityEndpoint, c.organizationBaseUrl(), c.projectId, environmentId, APIVersion7_1Preview1)
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
