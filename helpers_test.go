/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func resetOsArgs() {
	os.Args = []string{}
	for _, s := range cmdr.SavedOsArgs {
		os.Args = append(os.Args, s)
	}
}

func prepareConfD(t *testing.T) func() {
	cmdr.SetPredefinedLocations([]string{"./.tmp.yaml"})

	var clcl = &cfgLoaded{}
	cfg(t, clcl)

	return func() {
		_ = os.Remove("conf.d/tmp.yaml")
		_ = os.Remove("conf.d/tmp.json")
		_ = os.Remove("conf.d/tmp.toml")
		_ = os.Remove("conf.d")
		_ = os.Remove(".tmp.json")
		_ = os.Remove(".tmp.toml")
		cmdr.SetOnConfigLoadedListener(clcl, false)
		cmdr.RemoveOnConfigLoadedListener(clcl)
	}
}

type cfgLoaded struct {
}

func (s *cfgLoaded) OnConfigReloaded() {
	//
}

func cfg(t *testing.T, clcl cmdr.ConfigReloaded) {
	cmdr.AddOnConfigLoadedListener(clcl)

	_ = ioutil.WriteFile(".tmp.yaml", []byte(`
app:
  debug: false
  ms:
    tags:
      modify:
        wed: [3, 4]
`), 0644)
	_ = cmdr.EnsureDir("conf.d")

	_ = ioutil.WriteFile("conf.d/tmp.yaml", []byte(`
app:
  debug: false
  ms:
    tags:
      modify:
        wed: [3, 4]
`), 0644)
	// _ = cmdr.LoadConfigFile(".tmp.json")
	// _ = cmdr.LoadConfigFile(".tmp.toml")
	if err := cmdr.LoadConfigFile(".tmp.yaml"); err != nil {
		t.Fatal(err)
	}

	t.Logf("%v, %v", cmdr.GetUsedConfigFile(), cmdr.GetUsedConfigSubDir())
	_ = ioutil.WriteFile("conf.d/tmp.yaml", []byte(`
app:
  debug: true
  ms:
    tags:
      modify:
        wed: [3, 4]
`), 0644)
	_ = ioutil.WriteFile("conf.d/tmp.json", []byte(`{"app":{"debug":false}}`), 0644)
	_ = ioutil.WriteFile("conf.d/tmp.toml", []byte(``), 0644)

}

type testStruct struct {
	Debug bool
}

type testServerStruct struct {
	Retry int
	Enum  string
	Tail  int
	Head  int
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

func resetFlagsAndLog(t *testing.T) {
	// reset all option values
	cmdr.Set("kv.port", 8500)
	cmdr.Set("ms.tags.port", 8500)
	cmdr.SetNx("app.help", false)
	cmdr.SetNx("app.help-zsh", false)
	cmdr.SetNx("app.help-bash", false)
	cmdr.SetNx("app.debug", false) // = cmdr.Set("debug", false)
	cmdr.SetNx("debug", false)
	cmdr.SetNx("app.verbose", false)
	cmdr.SetNx("help", false)
	cmdr.Set("generate.shell.zsh", false)
	cmdr.Set("generate.shell.bash", false)

	// SetNx(key, nil) shouldn't clear an node owned children
	cmdr.Set("generate.shell", nil)
	if cmdr.GetMapR("generate.shell") == nil {
		t.Fatal("SetNx(key, nil) shouldn't clear an node owned children!!")
	}

	// cmdr.Set("app.generate.shell.auto", false)

	_ = os.Setenv("APP_DEBUG", "1")

	t.Log(cmdr.Get("app.debug"))
	t.Log(cmdr.GetR("debug"))
	t.Log(cmdr.GetBool("app.debug"))
	t.Log(cmdr.GetBoolR("debug"))
	t.Log(cmdr.GetBoolRP("", "debug"))
	t.Log(cmdr.GetInt("app.retry"))
	t.Log(cmdr.GetIntR("retry"))
	t.Log(cmdr.GetIntRP("", "retry"))
	t.Log(cmdr.GetInt64("app.retry"))
	t.Log(cmdr.GetInt64R("retry"))
	t.Log(cmdr.GetInt64RP("", "retry"))
	t.Log(cmdr.GetUint("app.retry"))
	t.Log(cmdr.GetUintP("app", "retry"))
	t.Log(cmdr.GetUintR("retry"))
	t.Log(cmdr.GetUintRP("", "retry"))
	t.Log(cmdr.GetUint64("app.retry"))
	t.Log(cmdr.GetUint64R("retry"))
	t.Log(cmdr.GetUint64RP("", "retry"))
	t.Log(cmdr.GetFloat32("app.retry"))
	t.Log(cmdr.GetFloat32P("app", "retry"))
	t.Log(cmdr.GetFloat32R("retry"))
	t.Log(cmdr.GetFloat32RP("", "retry"))
	t.Log(cmdr.GetFloat32P("app", "retry"))
	t.Log(cmdr.GetFloat64("app.retry"))
	t.Log(cmdr.GetFloat64R("retry"))
	t.Log(cmdr.GetFloat64RP("", "retry"))
	t.Log(cmdr.GetFloat64P("app", "retry"))
	t.Log(cmdr.GetString("app.version"))
	t.Log(cmdr.GetStringR("version"))
	t.Log(cmdr.GetStringRP("", "version"))
	t.Log(cmdr.GetStringP("", "app.version"))

	if cmdr.WrapWithRxxtPrefix("ms") != "app.ms" {
		t.Fatal("WrapWithRxxtPrefix failed")
	}

	t.Log(cmdr.GetMap("app.ms.tags"))
	t.Log(cmdr.GetMapR("app.ms.tags"))
	t.Log(cmdr.GetStringSlice("app.ms.tags.modify.set"))
	t.Log(cmdr.GetStringSliceP("app", "ms.tags.modify.set"))
	t.Log(cmdr.GetStringSliceR("ms.tags.modify.set"))
	t.Log(cmdr.GetStringSliceRP("ms.tags", "modify.set"))
	t.Log(cmdr.GetIntSlice("app.ms.tags.modify.xed"))
	t.Log(cmdr.GetIntSliceP("app", "ms.tags.modify.xed"))
	t.Log(cmdr.GetIntSliceR("ms.tags.modify.xed"))
	t.Log(cmdr.GetIntSliceRP("ms.tags", "modify.xed"))
	t.Log(cmdr.GetDuration("app.ms.tags.modify.v"))
	t.Log(cmdr.GetDurationP("app", "ms.tags.modify.v"))
	t.Log(cmdr.GetDurationR("ms.tags.modify.v"))
	t.Log(cmdr.GetDurationRP("ms.tags", "modify.v"))

	// comma separator string -> int slice
	t.Log(cmdr.GetIntSlice("app.ms.tags.modify.ued"))
	// string slice -> int slice
	t.Log(cmdr.GetIntSlice("app.ms.tags.modify.wed"))

	t.Log(cmdr.GetInt64P("app", "retry"))
	t.Log(cmdr.GetUintP("app", "retry"))
	t.Log(cmdr.GetUint64P("app", "retry"))
}

func postWorks(t *testing.T) {
	if cx := cmdr.FindSubCommand("ms", &rootCmd.Command); cx == nil {
		t.Fatal("cannot find `ms`")
	} else if cy := cmdr.FindSubCommand("list", cx); cy == nil {
		t.Fatal("cannot find `list`")
	} else if cz := cmdr.FindSubCommand("yy", cy); cz != nil {
		t.Fatal("should not find `yy` for 'ms list'")
	}
	if cx := cmdr.FindSubCommandRecursive("modify", &rootCmd.Command); cx == nil {
		t.Fatal("cannot find `tags`")
	} else {
		if cmdr.FindFlag("spasswd", cx) != nil {
			t.Fatal("should not find `spasswd` for 'ms tags modify'")
		}
	}
	if cmdr.FindFlag("spasswd", &rootCmd.Command) == nil {
		t.Fatal("cannot find `spasswd`")
	}
	if cmdr.FindFlagRecursive("add", &rootCmd.Command) == nil {
		t.Fatal("cannot find `add`")
	}
}
