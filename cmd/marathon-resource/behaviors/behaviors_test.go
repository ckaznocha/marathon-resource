package behaviors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/marathon"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/mocks"
	gomarathon "github.com/gambol99/go-marathon"
	"github.com/golang/mock/gomock"
)

func TestOut(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		mockMarathoner = mocks.NewMockMarathoner(ctrl)
	)
	defer ctrl.Finish()

	gomock.InOrder(
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(6).Return(gomarathon.DeploymentID{DeploymentID: "foo", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{}, errors.New("Something went wrong")),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "baz", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "quux", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "zork", Version: "bar"}, nil),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().CheckDeployment("foo").Times(3).Return(false, nil),
		mockMarathoner.EXPECT().CheckDeployment("bing").Times(1).Return(false, nil),
		mockMarathoner.EXPECT().CheckDeployment("foo").Times(3).Return(false, nil),
		mockMarathoner.EXPECT().CheckDeployment("bing").Times(1).Return(false, errors.New("something bad happened")),
		mockMarathoner.EXPECT().CheckDeployment("baz").Times(2).Return(true, nil),
		mockMarathoner.EXPECT().CheckDeployment("quux").Times(1).Return(false, errors.New("something bad happened")),
		mockMarathoner.EXPECT().CheckDeployment("zork").Times(2).Return(true, nil),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().DeleteDeployment("baz").Times(1).Return(nil),
		mockMarathoner.EXPECT().DeleteDeployment("zork").Times(1).Return(errors.New("no way")),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().RestartApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "bing", Version: "bar"}, nil),
		mockMarathoner.EXPECT().RestartApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{}, errors.New("no way")),
		mockMarathoner.EXPECT().RestartApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "bing", Version: "bar"}, nil),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().LatestVersions(gomock.Any(), "").Times(1).Return([]string{"bar"}, nil),
		mockMarathoner.EXPECT().LatestVersions(gomock.Any(), "").Times(2).Return([]string{"baz"}, nil),
		mockMarathoner.EXPECT().LatestVersions(gomock.Any(), "").Times(1).Return([]string{}, errors.New("no way")),
		mockMarathoner.EXPECT().LatestVersions(gomock.Any(), "").Times(1).Return([]string{"baz"}, nil),
		mockMarathoner.EXPECT().LatestVersions(gomock.Any(), "").Times(1).Return([]string{"baz"}, nil),
	)

	type args struct {
		input       InputJSON
		appJSONPath string
		apiclient   marathon.Marathoner
	}
	tests := []struct {
		name    string
		args    args
		want    IOOutput
		wantErr bool
	}{
		{
			"Works",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{Version: Version{Ref: "bar"}},
			false,
		},
		{
			"No update",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{Version: Version{Ref: "baz"}},
			false,
		},
		{
			"No update, restart",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2, RestartIfNoUpdate: true},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{Version: Version{Ref: "bar"}},
			false,
		},
		{
			"Errors fetching latest versions",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2, RestartIfNoUpdate: true},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Errors restarting app",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2, RestartIfNoUpdate: true},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Errors on second deployment check",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2, RestartIfNoUpdate: true},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Bad app json file",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "ajson", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Bad app json file",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app_bad.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Error from UpdateApp",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Deployment times out",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Check deployment errors",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Delete deployment errors",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
	}

	for _, tt := range tests {
		got, err := Out(tt.args.input, tt.args.appJSONPath, tt.args.apiclient)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Out() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Out() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIn(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		mockMarathoner = mocks.NewMockMarathoner(ctrl)
	)
	defer ctrl.Finish()

	gomock.InOrder(
		mockMarathoner.EXPECT().GetApp("bar", "foo").Times(1).Return(gomarathon.Application{Version: "foo"}, nil),
		mockMarathoner.EXPECT().GetApp("baz", "quux").Times(1).Return(gomarathon.Application{}, errors.New("Bad stuff")),
	)

	type args struct {
		input     InputJSON
		apiclient marathon.Marathoner
	}
	tests := []struct {
		name    string
		args    args
		want    IOOutput
		wantErr bool
	}{
		{
			"Works",
			args{
				input: InputJSON{
					Source:  Source{AppID: "bar"},
					Version: Version{Ref: "foo"},
				},
				apiclient: mockMarathoner,
			},
			IOOutput{Version: Version{Ref: "foo"}},
			false,
		},
		{
			"Errors",
			args{
				input: InputJSON{
					Source:  Source{AppID: "baz"},
					Version: Version{Ref: "quux"},
				},
				apiclient: mockMarathoner,
			},
			IOOutput{},
			true,
		},
	}
	for _, tt := range tests {
		got, err := In(tt.args.input, tt.args.apiclient)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. In() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. In() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCheck(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		mockMarathoner = mocks.NewMockMarathoner(ctrl)
	)
	defer ctrl.Finish()

	gomock.InOrder(
		mockMarathoner.EXPECT().LatestVersions("bar", "").Times(1).Return([]string{"a", "b", "c"}, nil),
		mockMarathoner.EXPECT().LatestVersions("bar", "").Times(1).Return([]string{}, errors.New("totally whack")),
	)

	type args struct {
		input     InputJSON
		apiclient marathon.Marathoner
	}
	tests := []struct {
		name    string
		args    args
		want    CheckOutput
		wantErr bool
	}{
		{
			"Works",
			args{
				input: InputJSON{
					Source:  Source{AppID: "bar"},
					Version: Version{Ref: ""},
				},
				apiclient: mockMarathoner,
			},
			CheckOutput{
				Version{Ref: "a"},
				Version{Ref: "b"},
				Version{Ref: "c"},
			},
			false,
		},
		{
			"Errors",
			args{
				input: InputJSON{
					Source:  Source{AppID: "bar"},
					Version: Version{Ref: ""},
				},
				apiclient: mockMarathoner,
			},
			CheckOutput{},
			true,
		},
	}
	for _, tt := range tests {
		got, err := Check(tt.args.input, tt.args.apiclient)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Check() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Check() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
