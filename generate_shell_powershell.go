package cmdr

import (
	"fmt"
	"io"
)

func (w *ExecWorker) genShellPowershell(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	fmt.Println(`# todo powershell`)
	return
}
