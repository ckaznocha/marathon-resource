package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
)

const (
	check = "check"
	in    = "in"
	out   = "out"
)

type (
	params struct {
		AppJSON      string     `json:"app_json"`
		Replacements []metadata `json:"replacements"`
	}
	authCreds struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	source struct {
		AppID     string     `json:"app_id"`
		URI       string     `json:"uri"`
		BasicAuth *authCreds `json:"basic_auth"`
	}
	version struct {
		Ref string `json:"ref"`
	}
	inputJSON struct {
		Params  params  `json:"params"`
		Source  source  `json:"source"`
		Version version `json:"version"`
	}
	checkOut []version
	metadata struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	ioOut struct {
		Version  version    `json:"version"`
		Metadata []metadata `json:"metadata"`
	}
)

var logger = logrus.New()

func main() {
	var (
		input   inputJSON
		decoder = json.NewDecoder(os.Stdin)
		/*encoder*/ _ = json.NewEncoder(os.Stdout)
	)

	if len(os.Args) < 2 {
		logger.Fatal("You must one or more arguments")
	}

	if err := decoder.Decode(&input); err != nil {
		logger.WithError(err).Fatal("Failed to decode stdin")
	}

	/*uri*/ _, err := url.Parse(input.Source.URI)
	if err != nil {
		logger.WithError(err).Fatalf("Malformed URI %s", input.Source.URI)
	}

	// m := newMarathoner(&http.Client{}, uri, source.BasicAuth)

	switch os.Args[1] {
	case check:
		//TODO: do check
	case in:
		//TODO: do in
	case out:
		//TODO: do out
	}
}
