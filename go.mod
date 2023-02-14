module github.com/hedzr/cmdr

go 1.18

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace gopkg.in/hedzr/errors.v3 => ../05.errors

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../libs/logex

//replace github.com/hedzr/deepcopy => ../30.deepcopy

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/hedzr/cmdr-base v1.0.0
	github.com/hedzr/evendeep v0.3.1
	github.com/hedzr/log v1.6.1
	github.com/hedzr/logex v1.6.1
	golang.org/x/crypto v0.6.0
	golang.org/x/net v0.6.0
	gopkg.in/hedzr/errors.v3 v3.1.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
