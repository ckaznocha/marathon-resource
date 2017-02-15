package behaviors

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
)

func parsePayload(p Params, path string) (io.Reader, error) {
	var (
		replacements = map[string]string{}
		buf          = bytes.NewBuffer([]byte{})
	)
	replacements = replaceStrings(p.Replacements, replacements)
	replacements, err := replaceFiles(p.ReplacementFiles, replacements, path)
	if err != nil {
		return nil, err
	}

	tmpl, err := raymond.ParseFile(filepath.Join(path, p.AppJSON))
	if err != nil {
		return nil, err
	}
	app, err := tmpl.Exec(replacements)
	if err != nil {
		return nil, err
	}

	if _, err = buf.WriteString(app); err != nil {
		return nil, err
	}
	return buf, nil
}

func replaceStrings(
	metadata []Metadata,
	replacements map[string]string,
) map[string]string {
	for _, v := range metadata {
		replacements[v.Name] = v.Value
	}

	return replacements
}

func replaceFiles(
	metadata []Metadata,
	replacements map[string]string,
	path string,
) (map[string]string, error) {
	for _, v := range metadata {
		fileValue, err := ioutil.ReadFile(filepath.Join(path, v.Value))
		if err != nil {
			return replacements, fmt.Errorf(
				"Error replacing %s from replacement_files: %v",
				v.Name,
				err,
			)
		}
		replacements[v.Name] = strings.TrimSpace(string(fileValue))
	}

	return replacements, nil
}
