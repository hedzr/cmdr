module github.com/hedzr/cmdr/v2/loaders

go 1.21

// replace gopkg.in/hedzr/errors.v3 => ../../../24/libs.errors

// replace github.com/hedzr/go-errors/v2 => ../../libs.errors

// replace github.com/hedzr/env => ../../libs.env

// replace github.com/hedzr/is => ../../libs.is

// replace github.com/hedzr/logg => ../../libs.logg

// replace github.com/hedzr/go-utils/v2 => ../../libs.utils

// replace github.com/hedzr/evendeep => ../../libs.diff

// replace github.com/hedzr/go-common/v2 => ../../libs.common

// replace github.com/hedzr/go-log/v2 => ../../libs.log

// replace github.com/hedzr/store => ../../libs.store

// replace github.com/hedzr/store/codecs/hcl => ../../libs.store/codecs/hcl

// replace github.com/hedzr/store/codecs/hjson => ../../libs.store/codecs/hjson

// replace github.com/hedzr/store/codecs/json => ../../libs.store/codecs/json

// replace github.com/hedzr/store/codecs/nestext => ../../libs.store/codecs/nestext

// replace github.com/hedzr/store/codecs/toml => ../../libs.store/codecs/toml

// replace github.com/hedzr/store/codecs/yaml => ../../libs.store/codecs/yaml

// replace github.com/hedzr/store/providers/env => ../../libs.store/providers/env

// replace github.com/hedzr/store/providers/file => ../../libs.store/providers/file

replace github.com/hedzr/cmdr/v2 => ../

require (
	github.com/hedzr/cmdr/v2 v2.0.2
	github.com/hedzr/logg v0.5.20
	github.com/hedzr/store v1.0.7
	github.com/hedzr/store/codecs/hcl v1.0.7
	github.com/hedzr/store/codecs/hjson v1.0.7
	github.com/hedzr/store/codecs/json v1.0.7
	github.com/hedzr/store/codecs/nestext v1.0.7
	github.com/hedzr/store/codecs/toml v1.0.7
	github.com/hedzr/store/codecs/yaml v1.0.7
	github.com/hedzr/store/providers/env v1.0.7
	github.com/hedzr/store/providers/file v1.0.7
	gopkg.in/hedzr/errors.v3 v3.3.2
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hedzr/evendeep v1.1.10 // indirect
	github.com/hedzr/is v0.5.19 // indirect
	github.com/hjson/hjson-go/v4 v4.4.0 // indirect
	github.com/npillmayer/nestext v0.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
