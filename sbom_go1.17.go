// Copyright © 2022 Hedzr Yeh.

//go:build !go1.18
// +build !go1.18

package cmdr

// sbomAttacher adds `sbom` subcommand as cmdr RootCommand if it has not been defined by user.
//
// But nothing to do if golang version is too small.
func sbomAttacher(w *ExecWorker, root *RootCommand) {}
