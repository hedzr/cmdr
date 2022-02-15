package cmdr

import (
	"bufio"
	cmdrBase "github.com/hedzr/cmdr-base"
	"runtime"
	"sync"
)

// ExecWorker is a core logic worker and holder
type ExecWorker struct {
	switchCharset string

	// beforeXrefBuildingX, afterXrefBuiltX HookFunc
	beforeXrefBuilding      []HookFunc
	beforeConfigFileLoading []HookFunc
	afterConfigFileLoading  []HookFunc
	afterXrefBuilt          []HookFunc
	afterAutomaticEnv       []HookOptsFunc
	beforeHelpScreen        []HookHelpScreenFunc
	afterHelpScreen         []HookHelpScreenFunc
	preActions              []Handler
	postActions             []Invoker

	envPrefixes         []string
	rxxtPrefixes        []string
	predefinedLocations []string // predefined config file locations, the main config store
	secondaryLocations  []string // secondary locations, these configs will be merged into main config store
	alterLocations      []string // alter config file locations, so we can write back the changes
	pluginsLocations    []string
	extensionsLocations []string

	shouldIgnoreWrongEnumValue bool

	enableVersionCommands     bool
	enableHelpCommands        bool
	enableVerboseCommands     bool
	enableCmdrCommands        bool
	enableGenerateCommands    bool
	treatUnknownCommandAsArgs bool

	watchMainConfigFileToo   bool
	doNotLoadingConfigFiles  bool
	doNotWatchingConfigFiles bool
	confDFolderName          string
	watchChildConfigFiles    bool

	globalShowVersion   func()
	globalShowBuildInfo func()

	currentHelpPainter Painter

	bufferedStdio bool
	defaultStdout *bufio.Writer
	defaultStderr *bufio.Writer

	// rootCommand the root of all commands
	rootCommand *RootCommand

	rxxtOptions           *Options
	onOptionMergingSet    OnOptionSetCB
	onOptionSet           OnOptionSetCB
	savedOptions          []*Options
	writeBackAlterConfigs bool

	similarThreshold            float64 //
	noDefaultHelpScreen         bool    // disable printing help screen while '--help' hit
	noColor                     bool    // 'no-color'
	noEnvOverrides              bool    // 'no-env-overriders'
	strictMode                  bool    // 'strict-mode'
	noUnknownCmdTip             bool    // don't invoke unknownOptionHandler while unknown command found
	noCommandAction             bool    // disable invoking command action even if it's valid
	noPluggableAddons           bool    //
	noPluggableExtensions       bool    //
	noUseOnSwitchCharHitHandler bool    // don't invoke onSwitchCharHitHandler handler while a '-' found
	inCompleting                bool    // allow partial matching at last cmdline args

	logexInitialFunctor Handler
	logexPrefix         string
	logexSkipFrames     int

	afterArgsParsed Handler

	envVarToValueMap map[string]func() string

	helpTailLine string

	onSwitchCharHitHandler   OnSwitchCharHitCB
	onPassThruCharHitHandler OnPassThruCharHitCB

	addons []cmdrBase.PluginEntry

	lastPkg     *ptpkg
	hitCommands []*Command
	hitFlags    []*Flag
}

// GetHitCommands returns all matched sub-commands from commandline
func GetHitCommands() []*Command { return internalGetWorker().hitCommands }

// GetHitFlags returns all matched flags from commandline
func GetHitFlags() []*Flag { return internalGetWorker().hitFlags }

// ExecOption is the functional option for Exec()
type ExecOption func(w *ExecWorker)

func internalGetWorker() (w *ExecWorker) {
	uniqueWorkerLock.RLock()
	w = uniqueWorker
	uniqueWorkerLock.RUnlock()
	return
}

// internalResetWorkerNoLock makes a new instance of pointer
// to ExecWorker and updates uniqueWorker global variable.
//
func internalResetWorkerNoLock() (w *ExecWorker) {
	w = &ExecWorker{
		envPrefixes:  []string{"CMDR"},
		rxxtPrefixes: []string{"app"},

		predefinedLocations: []string{
			"./ci/etc/$APPNAME/$APPNAME.yml",       // for developer
			"/etc/$APPNAME/$APPNAME.yml",           // regular location
			"/usr/local/etc/$APPNAME/$APPNAME.yml", // regular macOS HomeBrew location
			"/opt/etc/$APPNAME/$APPNAME.yml",       // regular location
			"/var/lib/etc/$APPNAME/$APPNAME.yml",   // regular location
			"$HOME/.config/$APPNAME/$APPNAME.yml",  // per user
			"$THIS/$APPNAME.yml",                   // executable's directory
			"$APPNAME.yml",                         // current directory
			// "$XDG_CONFIG_HOME/$APPNAME/$APPNAME.yml", // ?? seldom defined | generally it's $HOME/.config
			// "./ci/etc/%s/%s.yml",
			// "/etc/%s/%s.yml",
			// "/usr/local/etc/%s/%s.yml",
			// "$HOME/.%s/%s.yml",
			// "$HOME/.config/%s/%s.yml",
		},

		secondaryLocations: []string{
			"/ci/etc/$APPNAME/conf/$APPNAME.yml",
			"/etc/$APPNAME/conf/$APPNAME.yml",
			"/usr/local/etc/$APPNAME/conf/$APPNAME.yml",
			"$HOME/.$APPNAME/$APPNAME.yml", // ext location per user
		},

		alterLocations: []string{
			"/ci/etc/$APPNAME/alter/$APPNAME.yml",
			"/etc/$APPNAME/alter/$APPNAME.yml",
			"/usr/local/etc/$APPNAME/alter/$APPNAME.yml",
			"./bin/$APPNAME.yml",              // for developer, current bin directory
			"/var/lib/$APPNAME/.$APPNAME.yml", //
			"$THIS/.$APPNAME.yml",             // executable's directory
		},

		pluginsLocations: []string{
			"./ci/local/share/$APPNAME/addons",
			"$HOME/.local/share/$APPNAME/addons",
			"$HOME/.$APPNAME/addons",
			"/usr/local/share/$APPNAME/addons",
			"/usr/share/$APPNAME/addons",
		},
		extensionsLocations: []string{
			"./ci/local/share/$APPNAME/ext",
			"$HOME/.local/share/$APPNAME/ext",
			"$HOME/.$APPNAME/ext",
			"/usr/local/share/$APPNAME/ext",
			"/usr/share/$APPNAME/ext",
		},

		shouldIgnoreWrongEnumValue: true,

		enableVersionCommands:     true,
		enableHelpCommands:        true,
		enableVerboseCommands:     true,
		enableCmdrCommands:        true,
		enableGenerateCommands:    true,
		treatUnknownCommandAsArgs: true,

		doNotLoadingConfigFiles: false,

		currentHelpPainter: new(helpPainter),

		defaultStdout: nil, //bufio.NewWriterSize(os.Stdout, 16384),
		defaultStderr: nil, //bufio.NewWriterSize(os.Stderr, 16384),

		rxxtOptions: newOptions(),

		similarThreshold:    similarThreshold,
		noDefaultHelpScreen: false,

		helpTailLine: defaultTailLine,

		confDFolderName: confDFolderNameConst,
	}

	WithEnvVarMap(nil)(w)

	w._setSwChars(runtime.GOOS)
	//if runtime.GOOS == "windows" {
	//	w.switchCharset = "-/~"
	//}

	uniqueWorker = w
	return
}

func (w *ExecWorker) _setSwChars(os string) {
	if os == "windows" {
		w.switchCharset = "-/~"
	} else {
		w.switchCharset = "-~/"
	}
	//if sw, ok := switchCharMap[runtime.GOOS]; ok {
	//	w.switchCharset = sw
	//} else {
	//	w.switchCharset = "-~/"
	//}
}

func init() {
	onceWorkerInitial.Do(func() {

		noResetWorker = true
		//switchCharMap = map[string]string{
		//	"windows": "-/~",
		//}

		// create the uniqueWorker first time
		_ = internalResetWorkerNoLock()
	})
}

var onceWorkerInitial sync.Once   // once initializer for some global variables
var uniqueWorkerLock sync.RWMutex //
var uniqueWorker *ExecWorker      // NOTE that pointer to uniqueWorker can be updated, it's not an initial-once pointer
var noResetWorker bool            //
//var switchCharMap map[string]string //

const confDFolderNameConst = "conf.d"
