module github.com/hedzr/cmdr

go 1.17

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace gopkg.in/hedzr/errors.v3 => ../05.errors

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../15.logex

//replace github.com/hedzr/deepcopy => ../30.deepcopy

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/fsnotify/fsnotify v1.5.4
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/log v1.5.53
	github.com/hedzr/logex v1.5.53
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	gopkg.in/hedzr/errors.v3 v3.0.21
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
