module github.com/hedzr/cmdr

go 1.13

//replace github.com/hedzr/cmdr-base => ../cmdr-base

//replace github.com/hedzr/log => ../10.log

//replace github.com/hedzr/logex => ../logex

//replace gopkg.in/hedzr/errors.v2 => ../errors

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/hedzr/cmdr-base v0.1.3
	github.com/hedzr/log v0.3.19
	github.com/hedzr/logex v1.3.19
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc
	gopkg.in/hedzr/errors.v2 v2.1.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
