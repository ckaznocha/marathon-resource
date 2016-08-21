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
		AppJSON      string            `json:"app_json"`
		Replacements map[string]string `json:"replacements"`
	}
	source struct {
		AppID string `json:"app_id"`
		URI   string `json:"uri"`
	}
	version struct {
		Ref string `json:"ref"`
	}
	inputJSON struct {
		Params  params  `json:"params"`
		Source  source  `json:"source"`
		Version version `json:"version"`
	}
)

func main() {
	var (
		input   inputJSON
		logger  = logrus.New()
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

	// m := newMarathoner(http.DefaultClient, uri)

	switch os.Args[1] {
	case check:
		//TODO: do check
	case in:
		//TODO: do in
	case out:
		//TODO: do out
	}
}
