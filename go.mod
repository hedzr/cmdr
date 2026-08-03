module github.com/hedzr/cmdr/v2

go 1.25.0

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

require (
	github.com/hedzr/evendeep v1.4.3
	github.com/hedzr/is v0.9.5
	github.com/hedzr/logg v0.9.3
	github.com/hedzr/store v1.4.3
	github.com/hedzr/store/codecs/json v1.4.3
	github.com/hedzr/store/providers/file v1.4.3
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)
