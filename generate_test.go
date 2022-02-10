package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"os"
	"strings"
	"testing"
)

func TestGenShell1(t *testing.T) {
	cmdr.InternalResetWorkerForTest()
	cmdr.ResetOptions()

	// copyRootCmd = rootCmdForTesting
	var rootCmdX = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "consul-tags",
			},
		},
	}

	_ = dir.EnsureDir("man1")

	defer func() {
		_ = os.Remove(".tmp.bash")
		_ = dir.RemoveDirRecursive("man1")
		_ = dir.RemoveDirRecursive("man3")
		_ = dir.RemoveDirRecursive("conf.d")
	}()

	var commands = []string{
		"consul-tags gen sh --powershell",
		"consul-tags gen sh --fish",
		"consul-tags gen sh --zsh",
		"consul-tags gen sh --elvish",
		"consul-tags gen sh --fig",
		"consul-tags gen sh --bash",
		"consul-tags gen sh --bash -o .tmp.bash --dir man1/",
		"consul-tags gen sh --bash -o .tmp.bash",
		"consul-tags gen man",
		"consul-tags gen man -d man1",
		"consul-tags gen doc",
		"consul-tags gen doc -d man1",
		"consul-tags gen mkd -d man1",
		"consul-tags gen pdf -d man1",
		"consul-tags gen docx -d man1",
		"consul-tags gen tex -d man1",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)

		if err := cmdr.Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	cmdr.Set("generate.manual.dir", "man1")
	_, _ = cmdr.GenManualForCommandForTest(&rootCmdX.Command)
	t.Log("done")
}

func TestForGenerateCommands(t *testing.T) {
	copyRootCmd = rootCmdForTesting

	cmdr.InternalResetWorkerForTest()
	cmdr.ResetOptions()
	cmdr.Set("no-watch-conf-dir", true)

	defer func() {
		_ = os.Remove(".tmp.1.json")
		_ = os.Remove(".tmp.1.yaml")
		_ = os.Remove(".tmp.1.toml")
	}()

	var commands = []string{
		"consul-tags gen doc --markdown",
		"consul-tags gen shell --auto",
		"consul-tags gen shell --auto --force-bash",
		"consul-tags gen doc",
		"consul-tags gen pdf",
		"consul-tags gen docx",
		"consul-tags gen tex",
		"consul-tags gen markdown",
		"consul-tags gen d",
		"consul-tags gen doc --pdf",
		"consul-tags gen doc --tex",
		"consul-tags gen doc --doc",
		"consul-tags gen doc --docx",
		"consul-tags gen shell --bash",
		"consul-tags gen shell --zsh",
		"consul-tags gen shell",
	}
	for _, cc := range commands {
		cmdr.Set("generate.shell.zsh", false)
		cmdr.Set("generate.shell.bash", false)
		cmdr.Set("generate.shell.auto", false)
		cmdr.Set("generate.shell.force-bash", false)
		cmdr.Set("generate.doc.pdf", false)
		cmdr.Set("generate.doc.markdown", false)
		cmdr.Set("generate.doc.tex", false)
		cmdr.Set("generate.doc.doc", false)
		cmdr.Set("generate.doc.docx", false)

		os.Args = strings.Split(cc, " ")
		fmt.Printf("  . args = [%v], go ...\n", os.Args)
		// cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmdForTesting, cmdr.WithInternalOutputStreams(nil, nil)); err != nil {
			t.Fatal(err)
		}
		// time.Sleep(time.Second)
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestForGenerateDoc(t *testing.T) {
	copyRootCmd = rootCmdForTesting

	cmdr.InternalResetWorkerForTest()
	cmdr.ResetOptions()
	cmdr.Set("no-watch-conf-dir", true)

	defer func() {
		_ = dir.RemoveDirRecursive("docs")
	}()

	var commands = []string{
		"consul-tags gen doc",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		// cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmdForTesting, cmdr.WithInternalOutputStreams(nil, nil)); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestForGenerateMan(t *testing.T) {
	copyRootCmd = rootCmdForTesting

	cmdr.InternalResetWorkerForTest()
	cmdr.ResetOptions()
	cmdr.Set("no-watch-conf-dir", true)

	defer func() {
		_ = os.Remove("man1")
		_ = os.Remove("man3")
	}()

	var commands = []string{
		"consul-tags gen man",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		// cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmdForTesting, cmdr.WithInternalOutputStreams(nil, nil)); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}
