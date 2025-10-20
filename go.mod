module github.com/hedzr/cmdr/v2

go 1.24.0

toolchain go1.24.5

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/store/providers/file => ../libs.store/providers/file

require (
	github.com/hedzr/evendeep v1.3.61
	github.com/hedzr/is v0.8.61
	github.com/hedzr/logg v0.8.61
	github.com/hedzr/store v1.3.61
	github.com/hedzr/store/codecs/json v1.3.61
	github.com/hedzr/store/providers/file v1.3.61
	golang.org/x/exp v0.0.0-20251017212417-90e834f514db
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
)
