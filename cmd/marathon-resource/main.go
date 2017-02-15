package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/behaviors"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/marathon"
)

const (
	check = "check"
	in    = "in"
	out   = "out"
)

func main() {
	var (
		err     error
		input   behaviors.InputJSON
		output  interface{}
		decoder = json.NewDecoder(os.Stdin)
		encoder = json.NewEncoder(os.Stdout)

		logger   = logrus.New()
		logFatal = func(err error, msg string) {
			// This is just to make logrus play nice with the tests.
			// `Fatal()` calls os.Exit(int) which is annoying to test and
			// `WithError(error)` is hard to mock because it returns a concrete
			// type.
			if len(os.Getenv("GO_TESTING")) == 0 {
				logger.WithError(err).Fatal(msg)
			}
			panic(fmt.Sprintf("%s: %v", msg, err))
		}
	)

	logger.Out = os.Stderr

	if len(os.Args) < 2 {
		logFatal(
			fmt.Errorf("You must supply more than %d arguments", len(os.Args)),
			"Incorrect number of arguments",
		)
	}

	if err = decoder.Decode(&input); err != nil {
		logFatal(err, "Failed to decode stdin")
	}

	uri, err := url.Parse(input.Source.URI)
	if err != nil {
		logFatal(err, fmt.Sprintf("Malformed URI %s", input.Source.URI))
	}

	m := marathon.NewMarathoner(&http.Client{}, uri, input.Source.BasicAuth, input.Source.APIToken, logger)

	switch os.Args[1] {
	case check:
		if output, err = behaviors.Check(input, m); err != nil {
			logFatal(err, "Unable to get APP versions from marathon")
		}
	case in:
		if output, err = behaviors.In(input, m); err != nil {
			logFatal(err, "Unable to get APP info from marathon")
		}
	case out:
		if output, err = behaviors.Out(input, os.Args[2], m); err != nil {
			logFatal(err, "Unable to deploy APP to marathon")
		}
	default:
		logFatal(
			fmt.Errorf("%q is not a valid behavior", os.Args[1]),
			"Unknown argument",
		)
	}

	if err = encoder.Encode(output); err != nil {
		logFatal(err, "Failed to write output")
	}
}
