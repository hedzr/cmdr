/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/hedzr/cmdr"
	"reflect"
	"strings"
	"testing"
)

var (
// demoOptx = cmdr.NewOptionsWith(map[string]interface{}{
// 	"runmode":         "devel",
// 	"env-prefix":      "DOT",
// 	"app.version-sim": "",
// 	"app.version":     "",
// 	"app.logger.file": "",
// })
//
// demoOpts = &cmdr.OptOne{
// 	Children: map[string]*cmdr.OptOne{
// 		"runmode":    {Value: "devel"},
// 		"env-prefix": {Value: "DOT"},
// 		"app": {
// 			Children: map[string]*cmdr.OptOne{
// 				"generate": {
// 					Children: map[string]*cmdr.OptOne{
// 						"shell": {
// 							Children: map[string]*cmdr.OptOne{
// 								"bash": {Value: false},
// 								"zsh":  {Value: false},
// 								"auto": {Value: false},
// 							},
// 						},
// 						"manual": {
// 							Children: map[string]*cmdr.OptOne{
// 								"pdf": {Value: false},
// 								"tex": {Value: false},
// 							},
// 						},
// 					},
// 				},
// 				"ms": {
// 					Children: map[string]*cmdr.OptOne{
// 						"name": {Value: ""},
// 						"id":   {Value: ""},
// 						"tags": {
// 							// Value: nil,
// 							Children: map[string]*cmdr.OptOne{
// 								"ls":  {Value: true,},
// 								"add": {Value: true,},
// 								"rm": {
// 									Value: true,
// 									Children: map[string]*cmdr.OptOne{
// 										"list": {Value: []string{},},
// 									},
// 								},
// 								"toggle": {
// 									Children: map[string]*cmdr.OptOne{
// 										"set":   {Value: []string{},},
// 										"reset": {Value: []string{},},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }
)

func TestLoadConfigFile(t *testing.T) {

	err := cmdr.LoadConfigFile("../../ci/etc/devops/devops.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + cmdr.DumpAsString())

}

func TestDemoOptsWriting(t *testing.T) {

	// b, err := yaml.Marshal(demoOpts)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// err = ioutil.WriteFile("demo-opts.yaml", b, 0644)
	// if err != nil {
	// 	t.Fatal(err)
	// }

}

func doubleSlice(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() != reflect.Slice {
		fmt.Println("The interface is not a slice.")
		return nil
	}

	v := reflect.ValueOf(s)
	newLen := v.Len()
	newCap := (v.Cap() + 1) * 2
	typ := reflect.TypeOf(s).Elem()

	t := reflect.MakeSlice(reflect.SliceOf(typ), newLen, newCap)
	reflect.Copy(t, v)
	return t.Interface()
}

func TestReflectOfSlice(t *testing.T) {
	xs := doubleSlice([]string{"foo", "bar"}).([]string)
	fmt.Println("data =", xs, "len =", len(xs), "cap =", cap(xs))

	ys := doubleSlice([]int{3, 1, 4}).([]int)
	fmt.Println("data =", ys, "len =", len(ys), "cap =", cap(ys))
}

func TestExec(t *testing.T) {
	if rootCmd.SubCommands[1].SubCommands[0].Flags[0] == rootCmd.SubCommands[2].Flags[0] {
		t.Log(rootCmd.SubCommands[1].SubCommands[0].Flags)
		t.Log(rootCmd.SubCommands[2].Flags)
		t.Fatal("should not equal.")
	}

	flags := *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag)
	t.Log(flags)

	var err error
	var outX = bytes.NewBufferString("")
	var errX = bytes.NewBufferString("")
	var outBuf = bufio.NewWriterSize(outX, 16384)
	var errBuf = bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)

	defer func() {
		t.Log("--------- stdout")
		t.Log(outX.String())
	}()

	for sss, verifier := range execTestings {
		cmdr.Set("app.kv.port", 8500)
		cmdr.Set("app.ms.tags.port", 8500)

		if err = cmdr.InternalExecFor(rootCmd, strings.Split(sss, " ")); err != nil {
			t.Fatal(err)
		} else {
			if err = verifier(t); err != nil {
				t.Fatal(err)
			}
		}
	}

	if errX.Len() > 0 {
		t.Log("--------- stderr")
		t.Fatalf("Error!! %v", errX.String())
	}
}

var (
	// testing args
	execTestings = map[string]func(t *testing.T) error{
		"consul-tags ms tags modify -h ~~debug --port8509 --prefix/": func(t *testing.T) error {
			if cmdr.GetInt("app.ms.tags.port") != 8509 || cmdr.GetString("app.ms.tags.prefix") != "/" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") {
				return fmt.Errorf("something wrong 1. |%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"))
			}
			return nil
		},
		"consul-tags -? -vD --vv kv backup --prefix'' -h ~~debug": func(t *testing.T) error {
			if cmdr.GetInt("app.kv.port") != 8500 || cmdr.GetString("app.kv.prefix") != "" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetBool("app.verbose") || !cmdr.GetBool("app.debug") {
				return errors.New("something wrong 2.")
			}
			return nil
		},
		"consul-tags -vD --vv ms tags modify --prefix'' --help ~~debug --prefix\"\" --prefix'cmdr' --prefix\"app\" --prefix=/ --prefix/ --prefix /": func(t *testing.T) error {
			if cmdr.GetInt("app.ms.tags.port") != 8500 || cmdr.GetString("app.ms.tags.prefix") != "/" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetBool("app.verbose") || !cmdr.GetBool("app.debug") {
				return fmt.Errorf("something wrong 3. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"),
					cmdr.GetBool("app.verbose"), !cmdr.GetBool("app.debug"))
			}
			return nil
		},
		"consul-tags -vD ms tags modify --prefix'' -? ~~debug --port8509 -p8507 -p=8506 -p 8503": func(t *testing.T) error {
			if cmdr.GetInt("app.ms.tags.port") != 8503 || cmdr.GetString("app.ms.tags.prefix") != "" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetBool("app.verbose") || !cmdr.GetBool("app.debug") {
				return fmt.Errorf("something wrong 4. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"),
					cmdr.GetBool("app.verbose"), !cmdr.GetBool("app.debug"))
			}
			return nil
		},
	}

	// testing rootCmd

	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name:  "consul-tags",
				Flags: []*cmdr.Flag{
					// global options here.
				},
			},
			// PreAction:  Pre,
			// PostAction: Post,
			SubCommands: []*cmdr.Command{
				// dnsCommands,
				// playCommand,
				// generatorCommands,
				serverCommands,
				msCommands,
				kvCommands,
			},
		},

		AppName: "consul-tags",
		Version: "0.0.1",
		// Version:    consul_tags.Version,
		// VersionInt: consul_tags.VersionInt,
		Copyright: "consul-tags is an effective devops tool",
		Author:    "Hedzr Yeh <hedzrz@gmail.com>",
	}

	serverCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr"},
			Description: "server ops: for linux service/daemon.",
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup"},
					Description: "startup this system service/daemon.",
					// Action:impl.ServerStart,
				},
				// PreAction: impl.ServerStartPre,
				// PostAction: impl.ServerStartPost,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					Description: "stop this system service/daemon.",
					// Action:impl.ServerStop,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restart",
					Aliases:     []string{"reload"},
					Description: "restart this system service/daemon.",
					// Action:impl.ServerRestart,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Full:        "status",
					Aliases:     []string{"st"},
					Description: "display its running status as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "i",
					Full:        "install",
					Aliases:     []string{"setup"},
					Description: "install as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "u",
					Full:        "uninstall",
					Aliases:     []string{"remove"},
					Description: "remove from a system service/daemon.",
				},
			},
		},
	}

	kvCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "kvstore",
			Full:        "kv",
			Aliases:     []string{"kvstore"},
			Description: "consul kv store operations...",
			Flags:       *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "b",
					Full:        "backup",
					Aliases:     []string{"bk", "bf", "bkp"},
					Description: "Dump Consul's KV database to a JSON/YAML file",
					// Action:      kvBackup,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "o",
								Full:                    "output",
								Description:             "Write output to a file (*.json / *.yml)",
								DefaultValuePlaceholder: "FILE",
							},
							DefaultValue: "consul-backup.json",
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restore",
					Description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					// Action:      kvRestore,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "i",
								Full:                    "input",
								Description:             "read the input file (*.json / *.yml)",
								DefaultValuePlaceholder: "FILE",
							},
							DefaultValue: "consul-backup.json",
						},
					},
				},
			},
		},
	}

	msCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "microservices",
			Full:        "ms",
			Aliases:     []string{"microservice", "micro-service"},
			Description: "micro-service operations...",
			Flags: []*cmdr.Flag{
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "n",
						Full:                    "name",
						Description:             "name of the service",
						DefaultValuePlaceholder: "NAME",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "i",
						Full:                    "id",
						Description:             "unique id of the service",
						DefaultValuePlaceholder: "ID",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "a",
						Full:        "all",
						Description: "all services",
					},
					DefaultValue: false,
				},
			},
		},
		SubCommands: []*cmdr.Command{
			tagsCommands,
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "l",
					Full:        "list",
					Aliases:     []string{"ls", "lst", "dir"},
					Description: "list services.",
					// Action:      msList,
				},
			},
		},
	}

	tagsCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Short:       "t",
			Full:        "tags",
			Aliases:     []string{},
			Description: "tags op.",
			Flags:       *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "ls",
					Full:        "list",
					Aliases:     []string{"l", "lst", "dir"},
					Description: "list tags.",
					// Action:      msTagsList,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "a",
					Full:        "add",
					Aliases:     []string{"create", "new"},
					Description: "add tags.",
					// Action:      msTagsAdd,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "ls",
								Full:                    "list",
								Aliases:                 []string{"l", "lst", "dir"},
								Description:             "a comma list to be added",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "rm",
					Aliases:     []string{"remove", "erase", "delete", "del"},
					Description: "remove tags.",
					// Action:      msTagsRemove,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "ls",
								Full:                    "list",
								Aliases:                 []string{"l", "lst", "dir"},
								Description:             "a comma list to be added.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "m",
					Full:        "modify",
					Aliases:     []string{"mod", "update", "change"},
					Description: "modify tags.",
					// Action:      msTagsModify,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "a",
								Full:                    "add",
								Description:             "a comma list to be added.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "r",
								Full:                    "rm",
								Aliases:                 []string{"remove", "erase", "del"},
								Description:             "a comma list to be removed.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "toggle",
					Aliases:     []string{"tog", "switch"},
					Description: "toggle tags for ms.",
					// Action:      msTagsToggle,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "x",
								Full:                    "address",
								Description:             "the address of the service (by id or name)",
								DefaultValuePlaceholder: "HOST:PORT",
							},
							DefaultValue: "",
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "s",
								Full:                    "set",
								Description:             "set to `tag` which service specified by --address",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "u",
								Full:                    "unset",
								Aliases:                 []string{"reset"},
								Description:             "and reset the others service nodes to `tag`",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
		},
	}

	consulConnectFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "a",
				Full:                    "addr",
				Description:             "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')",
				DefaultValuePlaceholder: "HOST[:PORT]",
			},
			DefaultValue: "consul.ops.local",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "p",
				Full:                    "port",
				Description:             "Consul port",
				DefaultValuePlaceholder: "PORT",
			},
			DefaultValue: 8500,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "K",
				Full:                    "insecure",
				Description:             "Skip TLS host verification",
				DefaultValuePlaceholder: "PORT",
			},
			DefaultValue: true,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "prefix",
				Description:             "Root key prefix",
				DefaultValuePlaceholder: "ROOT",
			},
			DefaultValue: "/",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "cacert",
				Description:             "Client CA cert",
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "cert",
				Description:             "Client cert",
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "scheme",
				Description:             "Consul connection scheme (HTTP or HTTPS)",
				DefaultValuePlaceholder: "SCHEME",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "u",
				Full:                    "username",
				Description:             "HTTP Basic auth user",
				DefaultValuePlaceholder: "USERNAME",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "pw",
				Full:                    "password",
				Aliases:                 []string{"passwd", "pwd"},
				Description:             "HTTP Basic auth password",
				DefaultValuePlaceholder: "PASSWORD",
			},
			DefaultValue: "",
		},
	}
)
