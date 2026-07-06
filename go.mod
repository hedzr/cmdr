module github.com/hedzr/cmdr/v2

go 1.25.0

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

require (
	github.com/hedzr/evendeep v1.4.0
	github.com/hedzr/is v0.9.3
	github.com/hedzr/logg v0.9.3
	github.com/hedzr/store v1.4.0
	github.com/hedzr/store/codecs/json v1.4.0
	github.com/hedzr/store/providers/file v1.4.0
	golang.org/x/exp v0.0.0-20260212183809-81e46e3db34a
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)
