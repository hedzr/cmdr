/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func resetOsArgs() {
	os.Args = []string{}
	for _, s := range tool.SavedOsArgs {
		os.Args = append(os.Args, s)
	}
}

func prepareStreams() (outX, errX *bytes.Buffer) {
	outX = bytes.NewBufferString("")
	errX = bytes.NewBufferString("")
	var outBuf = bufio.NewWriterSize(outX, 16384)
	var errBuf = bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	return
}

func prepareConfD(t *testing.T) func() {
	cmdr.SetPredefinedLocationsForTesting("./.tmp.yaml")

	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}

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
	t.Logf("%v, %v", cmdr.GetOptions(), cmdr.GetUsingConfigFiles())
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

func tLog(a ...interface{}) {}

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

	tLog(cmdr.Get("app.debug"))
	tLog(cmdr.GetR("debug"))
	tLog(cmdr.GetBool("app.debug"))
	tLog(cmdr.GetBoolR("debug"))
	tLog(cmdr.GetBoolRP("", "debug"))
	tLog(cmdr.GetBoolP("app", "debug"))
	tLog(cmdr.GetBool("app.debug", false))
	tLog(cmdr.GetBoolR("debug", false))
	tLog(cmdr.GetBoolRP("", "debug", false))
	tLog(cmdr.GetBoolP("app", "debug", false))

	tLog(cmdr.GetInt("app.retry"))
	tLog(cmdr.GetIntR("retry"))
	tLog(cmdr.GetIntRP("", "retry"))
	tLog(cmdr.GetIntP("app", "retry"))
	tLog(cmdr.GetInt64("app.retry"))
	tLog(cmdr.GetInt64R("retry"))
	tLog(cmdr.GetInt64RP("", "retry"))
	tLog(cmdr.GetInt64P("app", "retry"))
	tLog(cmdr.GetInt("app.retry", 1))
	tLog(cmdr.GetIntR("retry", 1))
	tLog(cmdr.GetIntRP("", "retry", 1))
	tLog(cmdr.GetIntP("app", "retry", 1))
	tLog(cmdr.GetInt64("app.retry", 1))
	tLog(cmdr.GetInt64R("retry", 1))
	tLog(cmdr.GetInt64RP("", "retry", 1))
	tLog(cmdr.GetInt64P("app", "retry", 1))
	tLog(cmdr.GetUint("app.retry"))
	tLog(cmdr.GetUintP("app", "retry"))
	tLog(cmdr.GetUintR("retry"))
	tLog(cmdr.GetUintRP("", "retry"))
	tLog(cmdr.GetUint64("app.retry"))
	tLog(cmdr.GetUint64R("retry"))
	tLog(cmdr.GetUint64RP("", "retry"))
	tLog(cmdr.GetUint64P("app", "retry"))
	tLog(cmdr.GetUint("app.retry", 1))
	tLog(cmdr.GetUintP("app", "retry", 1))
	tLog(cmdr.GetUintR("retry", 1))
	tLog(cmdr.GetUintRP("", "retry", 1))
	tLog(cmdr.GetUint64("app.retry", 1))
	tLog(cmdr.GetUint64R("retry", 1))
	tLog(cmdr.GetUint64RP("", "retry", 1))
	tLog(cmdr.GetUint64P("app", "retry", 1))

	tLog(cmdr.GetKibibytes("app.retry", 1))
	tLog(cmdr.GetKibibytesR("retry", 1))
	tLog(cmdr.GetKibibytesRP("", "retry", 1))
	tLog(cmdr.GetKibibytesP("app", "retry", 1))
	tLog(cmdr.GetKilobytes("app.retry", 1))
	tLog(cmdr.GetKilobytesR("retry", 1))
	tLog(cmdr.GetKilobytesRP("", "retry", 1))
	tLog(cmdr.GetKilobytesP("app", "retry", 1))

	tLog(cmdr.GetComplex64("app.retry"))
	tLog(cmdr.GetComplex64P("app", "retry"))
	tLog(cmdr.GetComplex64R("retry"))
	tLog(cmdr.GetComplex64RP("", "retry"))
	tLog(cmdr.GetComplex64P("app", "retry"))
	tLog(cmdr.GetComplex128("app.retry"))
	tLog(cmdr.GetComplex128R("retry"))
	tLog(cmdr.GetComplex128RP("", "retry"))
	tLog(cmdr.GetComplex128P("app", "retry"))

	tLog(cmdr.GetFloat32("app.retry"))
	tLog(cmdr.GetFloat32P("app", "retry"))
	tLog(cmdr.GetFloat32R("retry"))
	tLog(cmdr.GetFloat32RP("", "retry"))
	tLog(cmdr.GetFloat32P("app", "retry"))
	tLog(cmdr.GetFloat64("app.retry"))
	tLog(cmdr.GetFloat64R("retry"))
	tLog(cmdr.GetFloat64RP("", "retry"))
	tLog(cmdr.GetFloat64P("app", "retry"))
	tLog(cmdr.GetFloat32("app.retry", 1))
	tLog(cmdr.GetFloat32P("app", "retry", 1))
	tLog(cmdr.GetFloat32R("retry", 1))
	tLog(cmdr.GetFloat32RP("", "retry", 1))
	tLog(cmdr.GetFloat32P("app", "retry", 1))
	tLog(cmdr.GetFloat64("app.retry", 1))
	tLog(cmdr.GetFloat64R("retry", 1))
	tLog(cmdr.GetFloat64RP("", "retry", 1))
	tLog(cmdr.GetFloat64P("app", "retry", 1))

	tLog(cmdr.GetString("app.version"))
	tLog(cmdr.GetStringR("version"))
	tLog(cmdr.GetStringRP("", "version"))
	tLog(cmdr.GetStringP("", "app.version"))
	tLog(cmdr.GetString("app.version", ""))
	tLog(cmdr.GetStringR("version", ""))
	tLog(cmdr.GetStringRP("", "version", ""))
	tLog(cmdr.GetStringP("", "app.version", ""))

	tLog(cmdr.GetStringNoExpand("app.version", "1"))
	tLog(cmdr.GetStringNoExpandR("version", "2"))
	tLog(cmdr.GetStringNoExpandRP("", "version", "3"))
	tLog(cmdr.GetStringNoExpandP("", "app.version", "4"))

	if cmdr.WrapWithRxxtPrefix("ms") != "app.ms" {
		t.Fatal("WrapWithRxxtPrefix failed")
	}

	tLog(cmdr.GetMap("app.ms.tags"))
	tLog(cmdr.GetMapR("app.ms.tags"))
	tLog(cmdr.GetStringSlice("app.ms.tags.modify.set"))
	tLog(cmdr.GetStringSliceP("app", "ms.tags.modify.set"))
	tLog(cmdr.GetStringSliceR("ms.tags.modify.set"))
	tLog(cmdr.GetStringSliceRP("ms.tags", "modify.set"))
	tLog(cmdr.GetIntSlice("app.ms.tags.modify.xed"))
	tLog(cmdr.GetIntSliceP("app", "ms.tags.modify.xed"))
	tLog(cmdr.GetIntSliceR("ms.tags.modify.xed"))
	tLog(cmdr.GetIntSliceRP("ms.tags", "modify.xed"))

	tLog(cmdr.GetDuration("app.ms.tags.modify.v"))
	tLog(cmdr.GetDurationP("app", "ms.tags.modify.v"))
	tLog(cmdr.GetDurationR("ms.tags.modify.v"))
	tLog(cmdr.GetDurationRP("ms.tags", "modify.v"))
	tLog(cmdr.GetDuration("app.ms.tags.modify.v", time.Second))
	tLog(cmdr.GetDurationP("app", "ms.tags.modify.v", time.Second))
	tLog(cmdr.GetDurationR("ms.tags.modify.v", time.Second))
	tLog(cmdr.GetDurationRP("ms.tags", "modify.v", time.Second))

	tLog(cmdr.GetInt64Slice("app.ms.tags.modify.xed"))
	tLog(cmdr.GetInt64SliceP("app", "ms.tags.modify.xed"))
	tLog(cmdr.GetInt64SliceR("ms.tags.modify.xed"))
	tLog(cmdr.GetInt64SliceRP("ms.tags", "modify.xed"))
	tLog(cmdr.GetUint64Slice("app.ms.tags.modify.xed"))
	tLog(cmdr.GetUint64SliceP("app", "ms.tags.modify.xed"))
	tLog(cmdr.GetUint64SliceR("ms.tags.modify.xed"))
	tLog(cmdr.GetUint64SliceRP("ms.tags", "modify.xed"))

	// comma separator string -> int slice
	tLog(cmdr.GetIntSlice("app.ms.tags.modify.ued"))
	// string slice -> int slice
	tLog(cmdr.GetIntSlice("app.ms.tags.modify.wed"))

	tLog(cmdr.GetInt64P("app", "retry"))
	tLog(cmdr.GetUintP("app", "retry"))
	tLog(cmdr.GetUint64P("app", "retry"))

	cmdr.Set("ms.tags.modify.v", "")
	tLog(cmdr.GetDuration("app.ms.tags.modify.v"))
	cmdr.Set("ms.tags.modify.v", "3s")
	tLog(cmdr.GetDuration("app.ms.tags.modify.v"))

}

func postWorks(t *testing.T) {
	if cx := cmdr.FindSubCommand("ms", &rootCmdForTesting.Command); cx == nil {
		t.Fatal("cannot find `ms`")
	} else if cy := cmdr.FindSubCommand("list", cx); cy == nil {
		t.Fatal("cannot find `list`")
	} else if cz := cmdr.FindSubCommand("yy", cy); cz != nil {
		t.Fatal("should not find `yy` for 'ms list'")
	}
	if cx := cmdr.FindSubCommandRecursive("modify", &rootCmdForTesting.Command); cx == nil {
		t.Fatal("cannot find `tags`")
	} else {
		if cmdr.FindFlag("spasswd", cx) != nil {
			t.Fatal("should not find `spasswd` for 'ms tags modify'")
		}
	}
	if cmdr.FindFlag("spasswd", &rootCmdForTesting.Command) == nil {
		t.Fatal("cannot find `spasswd`")
	}
	if cmdr.FindFlagRecursive("add", &rootCmdForTesting.Command) == nil {
		t.Fatal("cannot find `add`")
	}
}
