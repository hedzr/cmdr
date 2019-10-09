/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"os"
)

// QuitSignals return a channel for quit signal raising up.
func QuitSignals() chan os.Signal {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	return quit
}

// // StopSelf will terminate the app gracefully
// func StopSelf() {
// 	//
// }
