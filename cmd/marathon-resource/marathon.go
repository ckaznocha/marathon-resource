package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	gomarathon "github.com/gambol99/go-marathon"
)

const (
	pathApp          = "/v2/apps/%s"
	pathAppVersions  = "/v2/apps/%s/versions"
	pathAppAtVersion = "/v2/apps/%s/versions/%s"
	pathDeployments  = "/v2/deployments"
	pathDeployment   = "/v2/deployments/%s"

	jsonContentType = "application/json"
)

type (
	doer interface {
		Do(req *http.Request) (*http.Response, error)
	}
	//Marathoner is an interface to interact with marathon
	marathoner interface {
		LatestVersions(appID string, version string) ([]string, error)
		GetApp(appID, version string) (gomarathon.Application, error)
		UpdateApp(gomarathon.Application) (gomarathon.DeploymentID, error)
		CheckDeployment(deploymentID string) (bool, error)
		DeleteDeployment(deploymentID string) error
	}
	marathon struct {
		client doer
		url    *url.URL
	}
)

func newMarathoner(client doer, uri *url.URL) marathoner {
	return &marathon{client: client, url: uri}
}

func (m *marathon) handleReq(
	method string,
	path string,
	payload io.Reader,
	wantCode int,
	resObj interface{},
) error {
	u := *m.url
	u.Path = path
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", jsonContentType)
	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != wantCode {
		return fmt.Errorf("Expected %d response code but got %d", wantCode, res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(resObj); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (m *marathon) LatestVersions(appID, version string) ([]string, error) {
	var v gomarathon.ApplicationVersions
	if err := m.handleReq(
		http.MethodGet,
		fmt.Sprintf(pathAppVersions, appID),
		nil,
		http.StatusOK,
		&v,
	); err != nil {
		return nil, err
	}
	return newerTimestamps(v.Versions, version)
}

func (m *marathon) GetApp(appID, version string) (gomarathon.Application, error) {
	var app gomarathon.Application
	err := m.handleReq(
		http.MethodGet,
		fmt.Sprintf(pathAppAtVersion, appID, version),
		nil,
		http.StatusOK,
		&app,
	)
	return app, err
}

func (m *marathon) UpdateApp(inApp gomarathon.Application) (gomarathon.DeploymentID, error) {
	var (
		payload, _ = json.Marshal(inApp)
		deployment gomarathon.DeploymentID
	)
	err := m.handleReq(
		http.MethodPut,
		fmt.Sprintf(pathApp, inApp.ID),
		bytes.NewReader(payload),
		http.StatusOK,
		&deployment,
	)
	return deployment, err
}

func (m *marathon) CheckDeployment(deploymentID string) (bool, error) {
	var (
		deployments []gomarathon.Deployment
	)
	err := m.handleReq(
		http.MethodGet,
		pathDeployments,
		nil,
		http.StatusOK,
		&deployments,
	)
	fmt.Println(err)

	for _, v := range deployments {
		fmt.Println(v.ID)
		if v.ID == deploymentID {
			return true, nil
		}
	}
	return false, err
}

func (m *marathon) DeleteDeployment(deploymentID string) error {
	return m.handleReq(
		http.MethodDelete,
		fmt.Sprintf(pathDeployment, deploymentID),
		nil,
		http.StatusOK,
		nil,
	)
}
