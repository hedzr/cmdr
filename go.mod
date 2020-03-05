module github.com/hedzr/cmdr

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/hedzr/logex v1.1.5
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/hedzr/errors.v2 v2.0.11
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71
)

// exclude github.com/sirupsen/logrus v1.4.1
// exclude github.com/sirupsen/logrus v1.4.2

// replace github.com/hedzr/logex v0.0.0 => ../logex

// replace github.com/hedzr/pools v0.0.0 => ../pools

// replace github.com/hedzr/errors => ../errors

// replace gopkg.in/hedzr/errors.v2 => ../errors
