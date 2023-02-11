module github.com/hedzr/cmdr

go 1.18

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace gopkg.in/hedzr/errors.v3 => ../05.errors

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../15.logex

//replace github.com/hedzr/deepcopy => ../30.deepcopy

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/evendeep v0.3.1
	github.com/hedzr/log v1.5.57
	github.com/hedzr/logex v1.5.57
	golang.org/x/crypto v0.5.0
	golang.org/x/net v0.5.0
	gopkg.in/hedzr/errors.v3 v3.0.23
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/term v0.4.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
