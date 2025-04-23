module github.com/hedzr/cmdr/v2

go 1.23.0

toolchain go1.23.3

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

require (
	github.com/hedzr/evendeep v1.3.15
	github.com/hedzr/is v0.7.15
	github.com/hedzr/logg v0.8.15
	github.com/hedzr/store v1.3.15
	github.com/hedzr/store/codecs/json v1.3.15
	github.com/hedzr/store/providers/file v1.3.15
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0
	golang.org/x/term v0.31.0
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
