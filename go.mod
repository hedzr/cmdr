module github.com/hedzr/cmdr

go 1.17

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../15.logex

//replace gopkg.in/hedzr/errors.v2 => ../05.errors

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/log v1.5.11
	github.com/hedzr/logex v1.5.11
	golang.org/x/crypto v0.0.0-20220210151621-f4118a5b28e2
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	gopkg.in/hedzr/errors.v2 v2.1.9
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
