/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "os"

// AddOnBeforeXrefBuilding add hook func
// Deprecated from v1.5.0
func AddOnBeforeXrefBuilding(cb HookXrefFunc) {
	uniqueWorker.AddOnBeforeXrefBuilding(cb)
}

// AddOnAfterXrefBuilt add hook func
// Deprecated from v1.5.0
func AddOnAfterXrefBuilt(cb HookXrefFunc) {
	uniqueWorker.AddOnAfterXrefBuilt(cb)
}

// ExecWith is main entry of `cmdr`.
// Deprecated from v1.5.0
func ExecWith(rootCmd *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookXrefFunc) (err error) {
	w := uniqueWorker
	if beforeXrefBuildingX != nil {
		w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
	}
	if afterXrefBuiltX != nil {
		w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
	}
	err = w.InternalExecFor(rootCmd, os.Args)
	return
}
