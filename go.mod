module github.com/hedzr/cmdr/v2

go 1.22.7
toolchain go1.23.7

// github.com/hedzr/cmdr/v2 => ../cmdr.v3

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

// replace github.com/hedzr/cmdr/v2/loaders => ./loaders

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
	github.com/hedzr/evendeep v1.2.15
	github.com/hedzr/is v0.6.10
	github.com/hedzr/logg v0.7.22
	github.com/hedzr/store v1.2.15
	github.com/hedzr/store/codecs/json v1.2.15
	github.com/hedzr/store/providers/file v1.2.15
	golang.org/x/crypto v0.35.0
	golang.org/x/exp v0.0.0-20250210185358-939b2ce775ac
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
)
