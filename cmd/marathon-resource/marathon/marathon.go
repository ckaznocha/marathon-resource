package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/dates"
	gomarathon "github.com/gambol99/go-marathon"
)

const (
	pathApp          = "/v2/apps/%s"
	pathAppRestart   = "/v2/apps/%s/restart"
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
		GetApp(appID, version string) (gomarathon.Application, error)
		UpdateApp(gomarathon.Application) (gomarathon.DeploymentID, error)
		RestartApp(appID string) (gomarathon.DeploymentID, error)
		CheckDeployment(deploymentID string) (bool, error)
		DeleteDeployment(deploymentID string) error
	}
	marathon struct {
		client   doer
		url      *url.URL
		auth     *AuthCreds
		apiToken string
		logger   logrus.FieldLogger
	}

	//AuthCreds will be used for HTTP basic auth
	AuthCreds struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
)

//NewMarathoner returns a new marathoner
func NewMarathoner(
	client doer,
	uri *url.URL,
	auth *AuthCreds,
	apiToken string,
	logger logrus.FieldLogger) Marathoner {
	return &marathon{
		client:   client,
		url:      uri,
		auth:     auth,
		apiToken: apiToken,
		logger:   logger,
	}
}

func (m *marathon) handleReq(
	method string,
	resourcePath string,
	payload io.Reader,
	wantCodes []int,
	resObj interface{},
) error {
	u := *m.url
	u.Path = path.Join(u.Path, resourcePath)
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", jsonContentType)
	if m.auth != nil {
		req.SetBasicAuth(m.auth.UserName, m.auth.Password)
	}
	if m.apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token=%s", m.apiToken))
	}

	m.logger.WithFields(
		logrus.Fields{
			"Method": req.Method,
			"URL":    req.URL.String(),
		},
	).Info("Sending HTTP API request to Marathon")
	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

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

	if res.Body == nil || resObj == nil {
		return nil
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

func (m *marathon) GetApp(appID, version string) (gomarathon.Application, error) {
	var app gomarathon.Application
	err := m.handleReq(
		http.MethodGet,
		fmt.Sprintf(pathAppAtVersion, appID, version),
		nil,
		[]int{http.StatusOK},
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
		[]int{http.StatusOK, http.StatusCreated},
		&deployment,
	)
	return deployment, err
}

func (m *marathon) RestartApp(appID string) (gomarathon.DeploymentID, error) {
	var (
		deployment gomarathon.DeploymentID
	)
	err := m.handleReq(
		http.MethodPost,
		fmt.Sprintf(pathAppRestart, appID),
		nil,
		[]int{http.StatusOK},
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
