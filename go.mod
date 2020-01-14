module github.com/hedzr/cmdr

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/hedzr/logex v1.1.5
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/hedzr/errors.v2 v2.0.11
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
)

// exclude github.com/sirupsen/logrus v1.4.1
// exclude github.com/sirupsen/logrus v1.4.2

// replace github.com/hedzr/logex v0.0.0 => ../logex

// replace github.com/hedzr/pools v0.0.0 => ../pools

// replace github.com/hedzr/errors => ../errors

// replace gopkg.in/hedzr/errors.v2 => ../errors
