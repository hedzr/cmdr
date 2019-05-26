/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"text/template"
)

func tplApply(tmpl string, data interface{}) string {
	var w = new(bytes.Buffer)
	var tpl = template.Must(template.New("y").Parse(tmpl))
	if err := tpl.Execute(w, data); err != nil {
		logrus.Errorf("tpl execute error: %v", err)
	}
	return w.String()
}

func isRoot() bool {
	return os.Getuid() == 0
}

func shellRunAuto(name string, arg ...string) error {
	err, output := shellRun(name, arg...)
	if err != nil {
		logrus.Fatalf("shellRunAuto err: %v\n\noutput:\n%v", err, output.String())
	}
	return err
}

func shellRun(name string, arg ...string) (err error, output bytes.Buffer) {
	cmd := exec.Command(name, arg...)
	// cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &output
	err = cmd.Run()
	return
}
