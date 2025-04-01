module github.com/hedzr/cmdr/v2

go 1.23.0

toolchain go1.23.3

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

require (
	github.com/hedzr/evendeep v1.3.11
	github.com/hedzr/is v0.7.11
	github.com/hedzr/logg v0.8.11
	github.com/hedzr/store v1.3.11
	github.com/hedzr/store/codecs/json v1.3.11
	github.com/hedzr/store/providers/file v1.3.11
	golang.org/x/crypto v0.36.0
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
)
