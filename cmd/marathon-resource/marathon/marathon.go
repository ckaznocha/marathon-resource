package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/dates"
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
	Marathoner interface {
		LatestVersions(appID string, version string) ([]string, error)
		GetApp(appID, version string) (Application, error)
		UpdateApp(Application) (DeploymentID, error)
		CheckDeployment(deploymentID string) (bool, error)
		DeleteDeployment(deploymentID string) error
	}
	marathon struct {
		client doer
		url    *url.URL
		auth   *AuthCreds
	}

	//AuthCreds will be used for HTTP basic auth
	AuthCreds struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	//Application is a marathon application
	Application gomarathon.Application

	//DeploymentID is a marathon deploymentID
	DeploymentID gomarathon.DeploymentID
)

//NewMarathoner returns a new marathoner
func NewMarathoner(client doer, uri *url.URL, auth *AuthCreds) Marathoner {
	return &marathon{client: client, url: uri}
}

func (m *marathon) handleReq(
	method string,
	path string,
	payload io.Reader,
	wantCodes []int,
	resObj interface{},
) error {
	u := *m.url
	u.Path = path
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", jsonContentType)
	if m.auth != nil {
		req.SetBasicAuth(m.auth.UserName, m.auth.Password)
	}
	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	gotWantCode := false
	for _, wantCode := range wantCodes {
		if res.StatusCode == wantCode {
			gotWantCode = true
			break
		}
	}

	if !gotWantCode {
		return fmt.Errorf(
			"Expected one of %v responses code but got %d",
			wantCodes,
			res.StatusCode,
		)
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
		[]int{http.StatusOK},
		&v,
	); err != nil {
		return nil, err
	}
	return dates.NewerTimestamps(v.Versions, version)
}

func (m *marathon) GetApp(appID, version string) (Application, error) {
	var app Application
	err := m.handleReq(
		http.MethodGet,
		fmt.Sprintf(pathAppAtVersion, appID, version),
		nil,
		[]int{http.StatusOK},
		&app,
	)
	return app, err
}

func (m *marathon) UpdateApp(inApp Application) (DeploymentID, error) {
	var (
		payload, _ = json.Marshal(inApp)
		deployment DeploymentID
	)
	err := m.handleReq(
		http.MethodPut,
		fmt.Sprintf(pathApp, inApp.ID),
		bytes.NewReader(payload),
		[]int{http.StatusOK, http.StatusCreated},
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
		[]int{http.StatusOK},
		&deployments,
	)

	for _, v := range deployments {
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
		[]int{http.StatusOK},
		nil,
	)
}
