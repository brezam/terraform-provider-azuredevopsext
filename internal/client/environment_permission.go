package client

import (
	"fmt"
	"net/http"
	"strconv"
)

const pipelinePermissionsEndpoint = "%s/pipelines/pipelinePermissions/environment/%s?api-version=%s"

func (c *Client) ListEnvironmentPipelinePermissions(environmentId string) ([]PipelinePermission, error) {
	url := fmt.Sprintf(pipelinePermissionsEndpoint, c.projectBaseUrl(), environmentId, APIVersion7_1Preview1)
	req, err := c.newJsonRequest("GET", nil, url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var pipelinePermissionResponse PipelinePermissionList
	defer resp.Body.Close()
	unpackJSON(resp.Body, &pipelinePermissionResponse)
	return pipelinePermissionResponse.Pipelines, nil
}

func (c *Client) AddEnvironmentPipelinePermission(environmentId string, pipelineId string) error {
	return c.setEnvironmentPipelinePermission(environmentId, pipelineId, true)
}

func (c *Client) RemoveEnvironmentPipelinePermission(environmentId string, pipelineId string) error {
	return c.setEnvironmentPipelinePermission(environmentId, pipelineId, false)
}

func (c *Client) setEnvironmentPipelinePermission(environmentId string, pipelineId string, authorized bool) error {
	url := fmt.Sprintf(pipelinePermissionsEndpoint, c.projectBaseUrl(), environmentId, APIVersion7_1Preview1)
	pipeIdInt, err := strconv.Atoi(pipelineId)
	if err != nil {
		return err
	}
	body := &PipelinePermissionList{
		Pipelines: []PipelinePermission{{Id: pipeIdInt, Authorized: authorized}},
	}
	req, err := c.newJsonRequest("PATCH", body, url)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to add pipeline permission, response status: %d", resp.StatusCode)
	}
	return nil
}
