package actions

//Params holds the values supported in by the concourse `params` array
type Params struct {
	AppJSON      string     `json:"app_json"`
	Replacements []Metadata `json:"replacements"`
}

//AuthCreds will be used for HTTP basic auth
type AuthCreds struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

//Source holds the values supported in by the concourse `source` array
type Source struct {
	AppID     string     `json:"app_id"`
	URI       string     `json:"uri"`
	BasicAuth *AuthCreds `json:"basic_auth"`
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
