module github.com/hedzr/cmdr

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/hedzr/errors v1.1.7
	github.com/hedzr/logex v1.1.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

// exclude github.com/sirupsen/logrus v1.4.1
// exclude github.com/sirupsen/logrus v1.4.2

// replace github.com/hedzr/logex v0.0.0 => ../logex

// replace github.com/hedzr/pools v0.0.0 => ../pools

// replace github.com/hedzr/errors v0.0.0 => ../errors
