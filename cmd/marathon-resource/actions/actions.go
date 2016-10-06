package actions

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/marathon"
	gomarathon "github.com/gambol99/go-marathon"
)

//Params holds the values supported in by the concourse `params` array
type Params struct {
	AppJSON      string     `json:"app_json"`
	TimeOut      int        `json:"time_out"`
	Replacements []Metadata `json:"replacements"`
}

//Source holds the values supported in by the concourse `source` array
type Source struct {
	AppID     string              `json:"app_id"`
	URI       string              `json:"uri"`
	BasicAuth *marathon.AuthCreds `json:"basic_auth"`
}

//Version maps to a concousre version
type Version struct {
	Ref string `json:"ref"`
}

//InputJSON is what all concourse actions will pass to us
type InputJSON struct {
	Params  Params  `json:"params"`
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

//CheckOutput is what concourse expects as the result of a `check`
type CheckOutput []Version

//Metadata holds a concourse metadata entry
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//IOOutput is the return concourse expects from an `in` or and `out`
type IOOutput struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// Out shall delploy an APP to marathon based on marathon.json file.
func Out(input InputJSON, appJSONPath string, apiclient marathon.Marathoner) (IOOutput, error) {

	// jsondata, err := ioutil.ReadFile(filepath.Join(os.Args[2], appjsonpath))

	jsondata, err := parsePayload(input.Params, appJSONPath)
	if err != nil {
		return IOOutput{}, err
	}

	var marathonAPP gomarathon.Application
	if err := json.NewDecoder(jsondata).Decode(&marathonAPP); err != nil {
		return IOOutput{}, err
	}

	did, err := apiclient.UpdateApp(marathonAPP)

	if err != nil {
		return IOOutput{}, err
	}

	timer := time.NewTimer(time.Duration(input.Params.TimeOut) * time.Second)
	deploying := true

	// Check if APP was deployed.
deployloop:
	for {

		select {
		case <-timer.C:
			break deployloop
		default:
			var err error
			deploying, err = apiclient.CheckDeployment(did.DeploymentID)
			if err != nil {
				return IOOutput{}, err
			}
			if !deploying {
				break deployloop
			}
		}
		time.Sleep(1 * time.Second)
	}
	if deploying {
		err := apiclient.DeleteDeployment(did.DeploymentID)
		if err != nil {
			return IOOutput{}, err
		}
		return IOOutput{}, errors.New("Could not deply")
	}

	return IOOutput{Version: Version{Ref: did.Version}}, nil

}

// In shall fetch info on current version
func In(input InputJSON, apiclient marathon.Marathoner) (IOOutput, error) {

	app, err := apiclient.GetApp(input.Source.AppID, input.Version.Ref)
	if err != nil {
		return IOOutput{}, err
	}

	return IOOutput{Version: Version{Ref: app.Version}}, nil

}

// Check shall get the latest versions
func Check(input InputJSON, apiclient marathon.Marathoner) (CheckOutput, error) {

	versions, err := apiclient.LatestVersions(input.Source.AppID, input.Version.Ref)
	if err != nil {
		return CheckOutput{}, err
	}

	var out = CheckOutput{}
	for _, v := range versions {
		out = append(out, Version{Ref: v})
	}

	return out, nil

}
