package behaviors

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/aymerick/raymond"
	"io/ioutil"
	"strings"
)

func parsePayload(p Params, path string) (io.Reader, error) {
	var (
		replacments = map[string]string{}
		buf         = bytes.NewBuffer([]byte{})
	)
	for _, v := range p.Replacements {
		fileValue, err := ioutil.ReadFile(filepath.Join(path, v.Value))
		if err != nil {
			replacments[v.Name] = v.Value
		} else {
			replacments[v.Name] = strings.TrimSpace(string(fileValue))
		}
	}
	tmpl, err := raymond.ParseFile(filepath.Join(path, p.AppJSON))
	if err != nil {
		return nil, err
	}
	app, err := tmpl.Exec(replacments)
	if err != nil {
		return nil, err
	}

	if _, err = buf.WriteString(app); err != nil {
		return nil, err
	}
	return buf, nil
}
