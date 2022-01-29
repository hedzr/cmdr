package internal

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log"
	"github.com/hedzr/log/basics"
	"github.com/hedzr/log/closers"
	"gopkg.in/hedzr/errors.v2"
	"runtime"
	"sync"
)

// App is a global singleton GlobalApp instance
func App() *GlobalApp { return app_ }

// ---------------------------------------------

// DBX returns DB layer (wrapped on GORM)
//func (s *GlobalApp) DBX() dbl.DB { return s.dbx }

// GormDB returns the underlying GORM DB object in DB layer (for fast, simple coding)
//func (s *GlobalApp) GormDB() *gorm.DB { return s.dbx.DBE() }

func (s *GlobalApp) RootCommand() *cmdr.RootCommand { return s.cmd.GetRoot() }
func (s *GlobalApp) AppName() string                { return conf.AppName }
func (s *GlobalApp) AppTag() string                 { return conf.AppName } // appTag: appName or serviceID
func (s *GlobalApp) AppTitle() string               { return cmdr.GetStringR("app-title") }
func (s *GlobalApp) AppModuleName() string          { return cmdr.GetStringR("app-module-name") }

// Cache returns Cache/Redis Service
//func (s *GlobalApp) Cache() *cache.Hub { return s.cache }

// ---------------------------------------------

func (s *GlobalApp) Close() {
	log.Debug("* *App shutting down ...")
	s.Basic.Close()
}

// ---------------------------------------------

type GlobalApp struct {
	basics.Basic

	muInit sync.RWMutex
	//dbx    dbl.DB
	cmd *cmdr.Command

	//cache *cache.Hub

	//cron cron.Jobs
}

func createApp() {
	app_ = &GlobalApp{}
}

var once_ sync.Once
var app_ *GlobalApp

func init() {
	once_.Do(func() {
		createApp()
	})
}

// ---------------------------------------------

func (s *GlobalApp) Init(cmd *cmdr.Command, args []string) (err error) {
	// initialize all infrastructures here, such as: DB, Cache, MQ, ...

	log.Debugf("* *App initializing...OS: %v, ARCH: %v", runtime.GOOS, runtime.GOARCH)
	log.Debugf("  cmdr: InDebugging/IsDebuggerAttached: %v, DebugMode/TraceMode: %v/%v, LogLevel: %v", cmdr.InDebugging(), cmdr.GetDebugMode(), cmdr.GetTraceMode(), cmdr.GetLoggerLevel())

	s.cmd = cmd

	ce := errors.NewContainer("")
	ce.Attach(s.initDB())
	ce.Attach(s.initCache())
	ce.Attach(s.initCron())

	// TODO add your basic components initializations here

	err = ce.Error()

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

	return
}
