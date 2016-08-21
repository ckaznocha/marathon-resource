package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/mocks"
	gomarathon "github.com/gambol99/go-marathon"
	"github.com/golang/mock/gomock"
)

func Test_marathon_LatestVersions(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"versions":["2015-02-11T09:31:50.021Z","2014-03-01T23:42:20.938Z"]}`)),
		},
		nil,
	)
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusConflict,
			Body:       ioutil.NopCloser(strings.NewReader(`{"versions":["2015-02-11T09:31:50.021Z","2014-03-01T23:42:20.938Z"]}`)),
		},
		nil,
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		appID   string
		version string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{"Works", fields{mockClient, u}, args{"foo", "2015-02-11T09:31:50.021Z"}, []string{"2015-02-11T09:31:50.021Z"}, false},
		{"Errors", fields{mockClient, u}, args{"foo", "2015-02-11T09:31:50.021Z"}, nil, true},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		got, err := m.LatestVersions(tt.args.appID, tt.args.version)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.LatestVersions(%v) error = %v, wantErr %v", tt.name, tt.args.appID, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. marathon.LatestVersion(%v) = %v, want %v", tt.name, tt.args.appID, got, tt.want)
		}
	}
}

func Test_marathon_handleReq(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
		},
		nil,
	)
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
		},
		nil,
	)
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{]`)),
		},
		nil,
	)
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		nil,
		errors.New("Something went wrong"),
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		method   string
		path     string
		payload  io.Reader
		wantCode int
		resObj   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"No body",
			fields{mockClient, u},
			args{
				http.MethodGet,
				"/",
				nil,
				http.StatusOK,
				&map[string]string{},
			},
			false,
		},
		{
			"Bad Status code",
			fields{mockClient, u},
			args{
				http.MethodGet,
				"/",
				nil,
				http.StatusOK,
				&map[string]string{},
			},
			true,
		},
		{
			"Error",
			fields{mockClient, u},
			args{
				http.MethodGet,
				"/",
				nil,
				http.StatusOK,
				&map[string]string{},
			},
			true,
		},
		{
			"Error",
			fields{mockClient, u},
			args{
				"😂",
				"/",
				nil,
				http.StatusOK,
				&map[string]string{},
			},
			true,
		},
		{
			"Error",
			fields{mockClient, u},
			args{
				http.MethodGet,
				"/",
				nil,
				http.StatusOK,
				&[]string{},
			},
			true,
		},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		if err := m.handleReq(tt.args.method, tt.args.path, tt.args.payload, tt.args.wantCode, tt.args.resObj); (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.handleReq(%v, %v, %v, %v, %v) error = %v, wantErr %v", tt.name, tt.args.method, tt.args.path, tt.args.payload, tt.args.wantCode, tt.args.resObj, err, tt.wantErr)
		}
	}
}

func Test_marathon_GetApp(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	in, _ := json.Marshal(gomarathon.Application{ID: "hello-app"})
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(in)),
		},
		nil,
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		appID   string
		version string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    gomarathon.Application
		wantErr bool
	}{
		{"Works", fields{mockClient, u}, args{"hello-app", "2015-02-11T09:31:50.021Z"}, gomarathon.Application{ID: "hello-app"}, false},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		got, err := m.GetApp(tt.args.appID, tt.args.version)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.GetApp(%v, %v) error = %v, wantErr %v", tt.name, tt.args.appID, tt.args.version, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. marathon.GetApp(%v, %v) = %v, want %v", tt.name, tt.args.appID, tt.args.version, got, tt.want)
		}
	}
}

func Test_marathon_UpdateApp(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	out, _ := json.Marshal(gomarathon.DeploymentID{DeploymentID: "foo"})
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(out)),
		},
		nil,
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		inApp gomarathon.Application
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    gomarathon.DeploymentID
		wantErr bool
	}{
		{"Works", fields{mockClient, u}, args{gomarathon.Application{ID: "foo-app"}}, gomarathon.DeploymentID{DeploymentID: "foo"}, false},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		got, err := m.UpdateApp(tt.args.inApp)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.UpdateApp(%v) error = %v, wantErr %v", tt.name, tt.args.inApp, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. marathon.UpdateApp(%v) = %v, want %v", tt.name, tt.args.inApp, got, tt.want)
		}
	}
}

func Test_marathon_CheckDeployment(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	out, _ := json.Marshal([]gomarathon.Deployment{{ID: "foo", XXStepsRaw: json.RawMessage("{}")}})
	mockClient.EXPECT().Do(gomock.Any()).Times(2).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(out)),
		},
		nil,
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		deploymentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{"Works", fields{mockClient, u}, args{"foo"}, true, false},
		{"Works", fields{mockClient, u}, args{"bar"}, false, false},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		got, err := m.CheckDeployment(tt.args.deploymentID)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.CheckDeployment(%v) error = %v, wantErr %v", tt.name, tt.args.deploymentID, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. marathon.CheckDeployment(%v) = %v, want %v", tt.name, tt.args.deploymentID, got, tt.want)
		}
	}
}

func Test_marathon_DeleteDeployment(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		mockClient = mocks.NewMockdoer(ctrl)
		u, _       = url.Parse("http://foo.bar/")
	)
	defer ctrl.Finish()
	mockClient.EXPECT().Do(gomock.Any()).Times(1).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		},
		nil,
	)
	type fields struct {
		client doer
		url    *url.URL
	}
	type args struct {
		deploymentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Works", fields{mockClient, u}, args{"foo"}, false},
	}
	for _, tt := range tests {
		m := &marathon{
			client: tt.fields.client,
			url:    tt.fields.url,
		}
		if err := m.DeleteDeployment(tt.args.deploymentID); (err != nil) != tt.wantErr {
			t.Errorf("%q. marathon.DeleteDeployment(%v) error = %v, wantErr %v", tt.name, tt.args.deploymentID, err, tt.wantErr)
		}
	}
}

func Test_newMarathoner(t *testing.T) {
	type args struct {
		client doer
		uri    *url.URL
	}
	tests := []struct {
		name string
		args args
		want marathoner
	}{
		{"Works", args{http.DefaultClient, &url.URL{}}, &marathon{http.DefaultClient, &url.URL{}}},
	}
	for _, tt := range tests {
		if got := newMarathoner(tt.args.client, tt.args.uri); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. newMarathoner() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
