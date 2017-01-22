package main

import (
	"encoding/json"
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

var logger = logrus.New()

func main() {
	var (
		input   behaviors.InputJSON
		decoder = json.NewDecoder(os.Stdin)
		encoder = json.NewEncoder(os.Stdout)
	)

	logger.Out = os.Stderr
	logger.Level = logrus.InfoLevel

	if len(os.Args) < 2 {
		logger.Fatal("You must one or more arguments")
	}

	if err := decoder.Decode(&input); err != nil {
		logger.WithError(err).Fatal("Failed to decode stdin")
	}

	uri, err := url.Parse(input.Source.URI)
	if err != nil {
		logger.WithError(err).Fatalf("Malformed URI %s", input.Source.URI)
	}

	m := marathon.NewMarathoner(&http.Client{}, uri, input.Source.BasicAuth, input.Source.APIToken, logger)

	switch os.Args[1] {
	case check:
		output, err := behaviors.Check(input, m)
		if err != nil {
			logger.WithError(err).Fatalf("Unable to get APP info from marathon: %s", err)
		}
		if err = encoder.Encode(output); err != nil {
			logger.WithError(err).Fatalf("Failed to write output: %s", err)
		}
		return

	case in:
		output, err := behaviors.In(input, m)
		if err != nil {
			logger.WithError(err).Fatalf("Unable to get APP info from marathon: %s", err)
		}
		if err = encoder.Encode(output); err != nil {
			logger.WithError(err).Fatalf("Failed to write output: %s", err)
		}
		return

	case out:
		output, err := behaviors.Out(input, os.Args[2], m)
		if err != nil {
			logger.WithError(err).Fatalf("Unable to deply APP to marathon: %s", err)
		}
		if err = encoder.Encode(output); err != nil {
			logger.WithError(err).Fatalf("Failed to write output: %s", err)
		}
		return
	}
}
