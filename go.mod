module github.com/hedzr/cmdr

go 1.13

//replace github.com/hedzr/cmdr-base => ../00.cmdr-base

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../15.logex

//replace gopkg.in/hedzr/errors.v2 => ../05.errors

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.5.1
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/log v1.3.22
	github.com/hedzr/logex v1.3.22
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
	gopkg.in/hedzr/errors.v2 v2.1.5
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
