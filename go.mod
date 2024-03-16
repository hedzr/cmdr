module github.com/hedzr/cmdr/v2

go 1.21

//  github.com/hedzr/cmdr/v2 => ../cmdr.v3

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/go-errors/v2 => ../libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/evendeep => ../libs.diff

// replace github.com/hedzr/store => ../libs.store

// replace github.com/hedzr/go-utils/v2 => ../libs.utils

// replace github.com/hedzr/go-common/v2 => ../libs.common

// replace github.com/hedzr/go-log/v2 => ../libs.log

// replace github.com/hedzr/go-cabin/v2 => ../libs.cabin

replace github.com/hedzr/cmdr/v2/loaders => ./loaders

// replace github.com/hedzr/store/codecs/hcl => ../libs.store/codecs/hcl

// replace github.com/hedzr/store/codecs/hjson => ../libs.store/codecs/hjson

// replace github.com/hedzr/store/codecs/json => ../libs.store/codecs/json

// replace github.com/hedzr/store/codecs/nestext => ../libs.store/codecs/nestext

// replace github.com/hedzr/store/codecs/toml => ../libs.store/codecs/toml

// replace github.com/hedzr/store/codecs/yaml => ../libs.store/codecs/yaml

// replace github.com/hedzr/store/providers/consul => ../libs.store/providers/consul

// replace github.com/hedzr/store/providers/env => ../libs.store/providers/env

// replace github.com/hedzr/store/providers/etcd => ../libs.store/providers/etcd

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

// replace github.com/hedzr/store/providers/fs => ../libs.store/providers/fs

// replace github.com/hedzr/store/providers/maps => ../libs.store/providers/maps

require (
	github.com/hedzr/cmdr/v2/loaders v0.0.0-00010101000000-000000000000
	github.com/hedzr/evendeep v1.1.6
	github.com/hedzr/is v0.5.17
	github.com/hedzr/logg v0.5.13
	github.com/hedzr/store v0.5.1
	golang.org/x/crypto v0.21.0
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
	gopkg.in/hedzr/errors.v3 v3.3.1
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hedzr/store/codecs/hcl v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/codecs/hjson v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/codecs/json v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/codecs/nestext v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/codecs/toml v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/codecs/yaml v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/providers/env v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hedzr/store/providers/file v0.0.0-20240309055652-9bfe2570640d // indirect
	github.com/hjson/hjson-go/v4 v4.4.0 // indirect
	github.com/npillmayer/nestext v0.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
