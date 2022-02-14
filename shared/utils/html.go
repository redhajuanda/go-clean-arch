package utils

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// ParseHTMLTemplate parses html template
func ParseHTMLTemplate(templateFileName string, data interface{}) ([]byte, error) {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, errors.Wrap(err, "error parse file")
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "error execute template")
	}
	return buf.Bytes(), nil
}
