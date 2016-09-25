package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/actions"
)

const (
	check = "check"
	in    = "in"
	out   = "out"
)

var logger = logrus.New()

func main() {
	var (
		input   actions.InputJSON
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

	//m := marathon.NewMarathoner(&http.Client{}, uri, source.BasicAuth)

	switch os.Args[1] {
	case check:
		//TODO: do check
	case in:
		//TODO: do in
	case out:
		//TODO: do out
	}
}
