/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

// Worker returns unexported worker for testing
func Worker() *ExecWorker {
	return uniqueWorker
}

