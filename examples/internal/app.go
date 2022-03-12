package internal

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log"
	"github.com/hedzr/log/basics"
	"github.com/hedzr/log/closers"
	"github.com/hedzr/log/dir"
	"gopkg.in/hedzr/errors.v3"
	"runtime"
	"sync"
)

// App is a global singleton GlobalApp instance
func App() *GlobalApp { return appUniqueInstance }

// ---------------------------------------------

// RootCommand returns cmdr.RootCommand
func (s *GlobalApp) RootCommand() *cmdr.RootCommand { return s.cmd.GetRoot() }

// AppName returns app name
func (s *GlobalApp) AppName() string {
	return sel(conf.AppName, s.cmd.GetRoot().Name, s.cmd.GetRoot().AppName)
}

// AppVersion returns app version
func (s *GlobalApp) AppVersion() string { return sel(conf.Version, s.cmd.GetRoot().Version) }

// AppTag returns app tag name (app name or service id)
func (s *GlobalApp) AppTag() string { return sel(conf.ServerTag, conf.ServerID, conf.AppName) } // appTag: appName or serviceID

// CmdrVersion returns app tag name (app name or service id)
func (s *GlobalApp) CmdrVersion() string { return sel(cmdr.GetString("cmdr.Version"), cmdr.Version) } // appTag: appName or serviceID

//// AppTitle returns app title line
//func (s *GlobalApp) AppTitle() string { return cmdr.GetStringR("app-title") }

//// AppModuleName returns app module name
//func (s *GlobalApp) AppModuleName() string { return cmdr.GetStringR("app-module-name") }

// DBX returns DB layer (wrapped on GORM)
//func (s *GlobalApp) DBX() dbl.DB { return s.dbx }

// GormDB returns the underlying GORM DB object in DB layer (for fast, simple coding)
//func (s *GlobalApp) GormDB() *gorm.DB { return s.dbx.DBE() }

// Cache returns Cache/Redis Service
//func (s *GlobalApp) Cache() *cache.Hub { return s.cache }

func sel(ss ...string) (ret string) {
	for _, s := range ss {
		if len(s) > 0 {
			ret = s
			break
		}
	}
	return
}

// ---------------------------------------------

// Close cleanups internal resources and free any basic infrastructure if necessary
func (s *GlobalApp) Close() {
	log.Debug("* *App shutting down ...")
	s.Basic.Close()
}

// ---------------------------------------------

// GlobalApp is a general global object
type GlobalApp struct {
	basics.Basic

	muInit sync.RWMutex
	cmd    *cmdr.Command

	//dbx    dbl.DB
	//cache  *cache.Hub
	//cron   cron.Jobs
}

var onceForApp sync.Once
var appUniqueInstance *GlobalApp

func init() {
	onceForApp.Do(func() {
		appUniqueInstance = &GlobalApp{}
	})
}

// ---------------------------------------------

// NewAppOption returns a cmdr.ExecOption so that you can attach it
// into your application.
//
func NewAppOption() cmdr.ExecOption {
	return func(w *cmdr.ExecWorker) {
		cmdr.WithGlobalPreActions(appUniqueInstance.Init)(w) // appUniqueInstance will be closed automatically

		//no need to do:
		//cmdr.WithGlobalPostActions(func(cmd *cmdr.Command, args []string) { appUniqueInstance.Close() })(w)

		//cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		//	cmdr.NewBool(false).
		//		Titles("trace", "tr").
		//		Description("enable trace mode for tcp/mqtt send/recv data dump", "").
		//		//Action(func(cmd *cmdr.Command, args []string) (err error) {
		//		//	println("trace mode on")
		//		//	cmdr.SetTraceMode(true)
		//		//	return
		//		//}).
		//		Group(cmdr.SysMgmtGroup).
		//		AttachToRoot(root)
		//}, nil)
	}
}

// Init do initial stuffs
func (s *GlobalApp) Init(cmd *cmdr.Command, args []string) (err error) {
	// initialize all infrastructures here, such as: DB, Cache, MQ, ...

	log.Debugf("* *App initializing...OS: %v, ARCH: %v", runtime.GOOS, runtime.GOARCH)
	log.Debugf("  cmdr: InDebugging/IsDebuggerAttached: %v, DebugMode/TraceMode: %v/%v, LogLevel: %v", cmdr.InDebugging(), cmdr.GetDebugMode(), cmdr.GetTraceMode(), cmdr.GetLoggerLevel())
	log.Debugf("  pwd: %v, exe: %v", dir.GetCurrentDir(), dir.GetExecutablePath())

	s.cmd = cmd

	ce := errors.New("")
	defer ce.Defer(&err)
	ce.Attach(s.initDB())
	ce.Attach(s.initCache())
	ce.Attach(s.initCron())

	// TODO add your basic components initializations here

	closers.RegisterPeripheral(s)
	return
}

func (s *GlobalApp) initCron() (err error) {
	s.muInit.Lock()
	defer s.muInit.Unlock()

	//if s.cron == nil {
	//	s.cron = cron.New().AddToPeripheral(&s.Basic)
	//}

	return
}

func (s *GlobalApp) initCache() (err error) {
	s.muInit.Lock()
	defer s.muInit.Unlock()

	//if s.cache == nil {
	//	s.cache = cache.New().AddToPeripheral(&s.Basic)
	//}

	return
}

func (s *GlobalApp) initDB() (err error) {
	s.muInit.Lock()
	defer s.muInit.Unlock()

	// do sth.

	return
}
