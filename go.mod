module github.com/hedzr/cmdr

go 1.13

// exclude github.com/sirupsen/logrus v1.4.1
// exclude github.com/sirupsen/logrus v1.4.2

// replace github.com/hedzr/log => ../log

// replace github.com/hedzr/logex => ../logex

// replace github.com/hedzr/pools => ../pools

// replace github.com/hedzr/errors => ../errors

// replace gopkg.in/hedzr/errors.v2 => ../errors

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/hedzr/log v0.2.0
	github.com/hedzr/logex v1.2.12
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	gopkg.in/hedzr/errors.v2 v2.1.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
