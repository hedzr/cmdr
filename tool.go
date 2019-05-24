/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

func manBr(s string) string {
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		lines = append(lines, l+"\n.br")
	}
	return strings.Join(lines, "\n")
}

func manExamples(s string, data interface{}) string {
	var (
		sources  = strings.Split(s, "\n")
		lines    []string
		lastLine string
	)
	for _, l := range sources {
		if strings.HasPrefix(l, "$ {{.AppName}}") {
			lines = append(lines, `.TP \w'{{.AppName}}\ 'u
.BI {{.AppName}} \ `+manWs(l[14:]))
		} else {
			if len(lastLine) == 0 {
				lastLine = strings.TrimSpace(l)
				// ignore multiple empty lines, compat them as one line.
				if len(lastLine) != 0 {
					lines = append(lines, lastLine+"\n.br")
				}
			} else {
				lastLine = strings.TrimSpace(l)
				lines = append(lines, lastLine+"\n.br")
			}
		}
	}
	return tplApply(strings.Join(lines, "\n"), data)
}

func tplApply(tmpl string, data interface{}) string {
	var w = new(bytes.Buffer)
	var tpl = template.Must(template.New("x").Parse(tmpl))
	if err := tpl.Execute(w, data); err != nil {
		logrus.Errorf("tpl execute error: %v", err)
	}
	return w.String()
}

func manWs(fmtStr string, args ...interface{}) string {
	str := fmt.Sprintf(fmtStr, args...)
	str = strings.ReplaceAll(strings.TrimSpace(str), " ", `\ `)
	return str
}
