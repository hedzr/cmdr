module github.com/hedzr/cmdr

go 1.16

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../15.logex

//replace gopkg.in/hedzr/errors.v2 => ../05.errors

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/log v1.5.0
	github.com/hedzr/logex v1.5.1
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
	gopkg.in/hedzr/errors.v2 v2.1.5
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
