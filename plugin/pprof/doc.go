/*
 * Copyright © 2021 Hedzr Yeh.
 */

// Package pprof provides the profiling command-line options and
// adapts to go tool pprof.
//
// Examples:
//
//    if err := cmdr.Exec(buildRootCmd(),
//        pprof.GetCmdrProfilingOptions(...),
//    ); err != nil {
//        log.Fatalf("error occurs in app running: %+v\n", err)
//    }
//
// Examples:
//
//    func buildRootCmd() (rootCmd *cmdr.RootCommand) {
//        root := cmdr.Root(appName, cmdr.Version)
//        rootCmd = root.RootCommand()
//        pprof.AttachToCmdr(root)
//        return
//    }
//
// Examples:
//
//    func buildRootCmd() (rootCmd *cmdr.RootCommand) {
//        root := cmdr.Root(appName, cmdr.Version).
//           Copyright(copyright, "hedzr")
//        rootCmd = root.RootCommand()
//        pprof.AttachToCmdr(root.RootCmdOpt(), ...)
//        return
//    }
//
//
package pprof
