package worker

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func TestWorkerS_Pre(t *testing.T) {
	ctx := context.Background()

	app, ww := cleanApp(t, ctx, true)

	// app := buildDemoApp()
	// ww := postBuild(app)
	// ww.InitGlobally()
	// assert.EqualTrue(t, ww.Ready())
	//
	// ww.ForceDefaultAction = true
	// ww.wrHelpScreen = &discardP{}
	// ww.wrDebugScreen = os.Stdout
	// ww.wrHelpScreen = os.Stdout

	ww.setArgs([]string{"--debug"})
	ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app

	err := ww.Run(ctx,
		withTasksBeforeParse(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			runner.Root().SelfAssert()
			t.Logf("root.SelfAssert() passed.")
			return
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkerS_Pre_v1(t *testing.T) {
	ctx := context.Background()
	app, ww := cleanApp(t, ctx, true)

	ww.setArgs([]string{"--debug", "-v"})
	ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app

	err := ww.Run(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkerS_Pre_v3(t *testing.T) {
	ctx := context.Background()
	app, ww := cleanApp(t, ctx, true)

	ww.setArgs([]string{"--debug", "-vv", "--verbose"})
	ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app

	err := ww.Run(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

type cmdrRunTest struct {
	args     string
	verifier taskAfterParse
	opts     []cli.Opt
}

type cmdrRunTests struct {
	// t     *testing.T
	items []cmdrRunTest
}

func TestWorkerS_Parse(t *testing.T) { //nolint:revive
	ctx := context.Background()
	for i, c := range testWorkerParseCases.items {
		if c.args == "" && c.verifier == nil {
			continue
		}

		t.Log()
		t.Log()
		t.Log()
		t.Logf("--------------- test #%d: Parsing %q\n", i, c.args)

		app, ww := cleanApp(t, ctx, false, c.opts...)
		ww.Config.Store = store.New()
		// ww.Config.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}

		ww.setArgs(append([]string{app.Name()}, strings.Split(c.args, " ")...))
		ww.tasksAfterParse = []taskAfterParse{c.verifier}
		ww.Config.TasksBeforeRun = []cli.Task{aTaskBeforeRun}
		err := ww.Run(ctx) // withTasksBeforeRun(taskBeforeRun),withTasksAfterParse(c.verifier))
		// err := app.Run()
		if err != nil {
			_ = app
			t.Fatal(err)
		}
	}
}

var (
	aTaskBeforeRun = func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { return } //nolint:revive

	testWorkerParseCases = cmdrRunTests{[]cmdrRunTest{
		{},
		{},

		// ~~tree
		{args: "~~tree", opts: []cli.Opt{
			withEnv(map[string]string{"FORCE_RUN": "1"}),
			withHelpScreenWriter(os.Stdout),
		}, verifier: func(w *workerS, ctx *parseCtx, errParsed error) (err error) { //nolint:revive
			if errorsv3.Is(errParsed, cli.ErrUnmatchedFlag) {
				logz.Print("[INFO] ErrUnmatchedFlag FOUND, that's expecting.", "err", errParsed)
				return nil
			}
			return errParsed
		}},

		// ~~tree 1
		{args: "ms t t --tree", verifier: func(w *workerS, ctx *parseCtx, errParsed error) (err error) { //nolint:revive
			if errorsv3.Is(errParsed, cli.ErrUnmatchedFlag) {
				logz.Print("[INFO] ErrUnmatchedFlag FOUND, that's expecting.", "err", errParsed)
				return nil
			}
			return errParsed
		}},

		// ~~tree 2
		{args: "ms t t ~~tree -v", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			if errorsv3.Is(errParsed, cli.ErrUnmatchedFlag) {
				logz.Fatal("ErrUnmatchedFlag FOUND, that's NOT expecting.")
			}
			ctx := context.TODO()
			if !pc.matchedFlags[pc.flag(ctx, "tree")].DblTilde {
				logz.Fatal("expecting DblTilde is true but fault.")
			}
			return errParsed
		}},

		// hit times
		{args: "m snd -n -wn cool fog --pp box", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			if !regexp.MustCompile(`UNKNOWN (Cmd|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				logz.Print("expect 'UNKNOWN Flag FOUND' error, but not matched.") // "--pp"
			}
			hitTest(pc, "dry-run", 2)
			hitTest(pc, "wet-run", 1)
			argsAre(pc, "cool", "fog")
			return nil /* errParsed */
		}},

		// general commands and flags
		{args: "jump to --full -f --dry-run", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "full", 2)
			hitTest(pc, "dry-run", 1)
			return errParsed
		}},

		// compact flags
		{args: "-qvqDq gen --debug sh --zsh -b -Dwmann --dry-run", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "quiet", 3)
			hitTest(pc, "debug", 3)
			hitTest(pc, "verbose", 1)
			hitTest(pc, "wet-run", 1)
			hitTest(pc, "dry-run", 2)
			return errParsed
		}},

		// general, unknown cmd/flg errors
		{args: "m snd --help"},
		{args: "m unk snd cool", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			if !regexp.MustCompile(`UNKNOWN (Cmd|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				logz.Print("expect 'UNKNOWN Cmd FOUND' error, but not matched.") // "unk"
			}
			return /* errParsed */
		}, opts: []cli.Opt{cli.WithUnmatchedAsError(true)}},
		{args: "m snd -n -wn cool fog --pp box", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			if !regexp.MustCompile(`UNKNOWN (Cmd|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				logz.Print("expect 'UNKNOWN Flag FOUND' error, but not matched.") // "--pp"
			}
			hitTest(pc, "dry-run", 2)
			hitTest(pc, "wet-run", 1)
			argsAre(pc, "cool", "fog")
			return nil /* errParsed */
		}},

		// headLike
		{args: "server start -f -129", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "foreground", 1)
			hitTest(pc, "head", 1)
			hitTest(pc, "tail", 0)
			valTest(pc, "head", 129) //nolint:revive
			return errParsed
		}},

		// toggle group
		{args: "generate shell --bash --zsh -p", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			ctx := context.Background()
			if f := pc.flag(ctx, "shell"); f != nil {
				assertEqual(f.MatchedTG().MatchedTitle, "powershell")
			}
			return errParsed
		}},

		// valid args
		{args: "server start -e apple -e zig", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			valTest(pc, "enum", "zig")
			return errParsed
		}},

		// parsing slice (-cr foo,bar,noz), compact flag with value (-mmt3)
		{args: "ms t modify -2 -cr foo,bar,noz -nfool -mmi3", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "money", 1)
			valTest(pc, "both", true)
			valTest(pc, "clear", true)
			valTest(pc, "name", "fool")
			valTest(pc, "id", "3")
			valTest(pc, "remove", []string{"foo", "bar", "noz"})
			return errParsed
		}},

		// parsing slice (-cr foo,bar,noz), compact flag with value (-mmt3)
		// merge/append to slice
		{args: "ms t modify -2 -cr foo,bar,noz -n fool -mmr 1", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "money", 1)
			valTest(pc, "both", true)
			valTest(pc, "clear", true)
			valTest(pc, "name", "fool")
			valTest(pc, "remove", []string{"foo", "bar", "noz", "1"})
			return errParsed
		}},

		{args: "ms t t -K -2 -cun foo,bar,noz", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "insecure", 1)
			valTest(pc, "insecure", true)
			valTest(pc, "both", true)
			valTest(pc, "clear", true)
			valTest(pc, "unset", []string{"foo", "bar", "noz"})
			return errParsed
		}},

		// parse duration
		{args: "m -dur 9s", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "dry-run", 0)
			hitTest(pc, "wet-run", 0)
			// argsAre(pc, "cool", "fog")

			valTest(pc, "duration", 9*time.Second) //nolint:revive
			return nil                             /* errParsed */
		}},

		// parse integer, float, complex
		{args: "m -i -9 -i64 3 -u 8 -u64 13 -f 3.14 -f64 2.718 -c64 1.23+4.567i -c128 2.313+9.87i", verifier: func(w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			hitTest(pc, "dry-run", 0)
			hitTest(pc, "wet-run", 0)
			// argsAre(pc, "cool", "fog")

			valTest(pc, "int", -9)                           //nolint:revive
			valTest(pc, "int64", int64(3))                   //nolint:revive
			valTest(pc, "uint", uint(8))                     //nolint:revive
			valTest(pc, "uint64", uint64(13))                //nolint:revive
			valTest(pc, "float32", float32(3.14))            //nolint:revive
			valTest(pc, "float64", 2.718)                    //nolint:revive
			valTest(pc, "complex64", complex64(1.23+4.567i)) //nolint:revive
			valTest(pc, "complex128", 2.313+9.87i)           //nolint:revive
			return nil                                       /* errParsed */
		}},

		{},
		{},
		{},
		{},
		{},
		{},
	}}
)

func argsAre(s *parseCtx, list ...string) {
	if !reflect.DeepEqual(s.positionalArgs, list) {
		panic(fmt.Sprintf("expect positional args are %v but got %v (for cmd %v)", list, s.positionalArgs, s.LastCmd()))
	}
}

func hitTest(s *parseCtx, longTitle string, times int) {
	cc := s.LastCmd()
	ctx := context.Background()
	if f := cc.FindFlagBackwards(ctx, longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if f.GetTriggeredTimes() != times {
		panic(fmt.Sprintf("expect hit times is %d but got %d (for flag %v)", times, f.GetTriggeredTimes(), f))
	}
}

func valTest(s *parseCtx, longTitle string, val any) {
	cc := s.LastCmd()
	ctx := context.Background()
	if f := cc.FindFlagBackwards(ctx, longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if !reflect.DeepEqual(f.DefaultValue(), val) {
		panic(fmt.Sprintf("expect flag's value is '%v' but got '%v' (for flag %v)", val, f.DefaultValue(), f))
	}
}

func assertEqual(expect, actual any, msgs ...any) {
	if expect != actual {
		logz.Fatal(fmt.Sprintf("expecting %v but got %v", actual, expect))
	}
	_ = msgs
}
