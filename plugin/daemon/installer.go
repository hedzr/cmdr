/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v2"
	"io/ioutil"
	"log"
	"os"
)

func runInstaller(cmd *cmdr.Command, args []string) (err error) {
	if !isRoot() {
		log.Fatal("This program must be run as root! (sudo)")
	}

	var (
		fileName string
		contents string
		data     *tplData
	)

	if cmdr.FileExists(systemdDir) {
		data = &tplData{
			*cmd.GetRoot(),
			cmdr.GetExcutablePath(),
		}

		fileName = fmt.Sprintf("%s/%s@.service", systemdDir, data.AppName)
		contents = tplApply(tplService, data)
		err = ioutil.WriteFile(fileName, []byte(contents), 0644)
		if err != nil {
			return
		}

		fileName = fmt.Sprintf("%v/%v", defaultsDir, data.AppName)
		contents = tplApply(tplDefault, data)
		err = ioutil.WriteFile(fileName, []byte(contents), 0755)
		if err != nil {
			return
		}

		_ = shellRunAuto("systemctl", "daemon-reload")

		name := fmt.Sprintf("%s@%s.service", data.AppName, os.Getenv("USER"))
		_ = shellRunAuto("systemctl", "enable", name)

		println(fmt.Sprintf("Usage:\n$ sudo systemctl daemon-reload\n$ sudo systemctl start %v", name))
	} else {
		println(fmt.Sprintf("Systemd not found, cannot install service."))
	}

	return
}

func runUninstaller(cmd *cmdr.Command, args []string) (err error) {
	if !isRoot() {
		log.Fatal("This program must be run as root! (sudo)")
	}

	var (
		fileName string
		data     *tplData
	)

	if cmdr.FileExists(systemdDir) {
		data = &tplData{
			*cmd.GetRoot(),
			cmdr.GetExcutablePath(),
		}

		fileName = fmt.Sprintf("%s/%s@.service", systemdDir, data.AppName)
		if cmdr.FileExists(fileName) {
			err = os.Remove(fileName)
		}
		_ = shellRunAuto("systemctl", "daemon-reload")
	}
	return
}

// ErrNoRoot error object: `MUST have administrator privileges`
var ErrNoRoot = errors.New("MUST have administrator privileges")

type tplData struct {
	cmdr.RootCommand
	BinPath string
}

const systemdDir = "/etc/systemd/system"

const defaultsDir = "/etc/default"

const tplDefault = `
### {{.AppName}} configurations
### executable: {{.BinPath}}

# PORT=3211

# OPTIONS="--port 3211"
GLOBAL_OPTIONS=""
OPTIONS=""

`

const tplService = `
#
[Unit]
Description={{.AppName}} Service for %i
# Documentation=man:sshd(8) man:sshd_config(5) man:{{.AppName}}(1)
After=network.target
# Wants=syslog.service
ConditionPathExists={{.BinPath}}

[Install]
WantedBy=multi-user.target

[Service]
#Type=idle
Type=forking
User=%i
#Group=%i
LimitNOFILE=65535

KillMode=process
Restart=on-failure
RestartSec=23s
# RestartLimitIntervalSec=60

EnvironmentFile=/etc/default/{{.AppName}}
WorkingDirectory=%h
#          start: --addr, --port,
#           todo: --pid
# global options: --config
ExecStart={{.BinPath}} $GLOBAL_OPTIONS server start $OPTIONS
#           stop: -1/--hup, -9/--kill,
### TODO ExecStop={{.BinPath}} $GLOBAL_OPTIONS server stop -1
### TODO ExecReload=/bin/kill -HUP $MAINPID
ExecStop={{.BinPath}} $GLOBAL_OPTIONS server stop -3
ExecReload={{.BinPath}} $GLOBAL_OPTIONS server restart

# # make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=-/bin/mkdir /var/run/{{.AppName}}
ExecStartPre=-/bin/mkdir /var/lib/{{.AppName}}
ExecStartPre=-/bin/mkdir /var/log/{{.AppName}}
ExecStartPre=-/bin/chown -R %i: /var/run/{{.AppName}} /var/lib/{{.AppName}}
# ExecStartPre=-/bin/chown -R syslog:adm /var/log/{{.AppName}}
ExecStartPre=-/bin/chown -R %i: /var/log/{{.AppName}}

# # enable coredump
# ExecStartPre=ulimit -c unlimited

StandardOutput=syslog
StandardError=syslog
SyslogIdentifier={{.AppName}}


`
