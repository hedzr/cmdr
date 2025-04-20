package worker

import (
	"context"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
)

type genDocS struct{}

func (w *genDocS) onAction(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive,unused
	outDir := cmd.Store().MustString("dir")
	fmt.Printf("# generating docpages (output-dir: %s) ...\n", outDir)
	return
}

//
//
// /////////////////////////////////////////
//
//

/*
func genDoc(command cli.Cmd, args []string) (err error) {
	prefix := strings.Join(append(internalGetWorker().rxxtPrefixes, "generate.doc"), ".")
	// logrus.Infof("OK gen doc: hit=%v", cmd.strHit)
	var painter Painter
	switch command.strHit {
	case "mkd", "m", "markdown":
		painter = newMarkdownPainter()
	case "pdf":
		painter = newMarkdownPainter()
	// case "man", "manual", "manpage", "man-page":
	// 	painter = newManPainter()
	// case "docx":
	// 	painter = newMarkdownPainter()
	// case "tex":
	// 	painter = newMarkdownPainter()
	default: // , "doc", "d"
		if GetBoolP(prefix, "markdown") { //nolint:gocritic //like it
			painter = newMarkdownPainter()
		} else if GetBoolP(prefix, "pdf") { //nolint:gocritic //like it
			painter = newMarkdownPainter()
			// } else if GetBoolP(prefix, "tex") {
			// 	painter = newMarkdownPainter()
		} else { //nolint:gocritic //like it
			painter = newMarkdownPainter()
		}
	}

	// fmt.Printf("  .  . args = [%v]\n", args)
	w := internalGetWorker()
	err = WalkAllCommands(func(cmd *Command, index, level int) (err error) {
		painter.Reset()
		// fmt.Printf("  .  .  cmd = %v\n", cmd.GetTitleNames())

		dirname := GetStringP(prefix, "dir")
		if err = dir.EnsureDir(dirname); err != nil {
			return
		}

		fn := cmd.root.AppName
		if !cmd.IsRoot() {
			cmds := replaceAll(backtraceCmdNames(cmd, false), ".", "-")
			if len(cmds) > 0 {
				fn += "-" + cmds
			}
		}
		fn = fmt.Sprintf("%s/%v.md", dirname, fn)

		w.paintFromCommand(painter, cmd, false)
		if err = dir.WriteFile(fn, painter.Results(), 0o600); err == nil {
			log.Printf("'%v' generated...", fn)
		}
		return
	})

	return
}
*/
