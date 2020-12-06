/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

type (
	tomlConfig struct {
		Title   string
		Runmode string
		App     app
		Owner   ownerInfo
		DB      database `toml:"database"`
		Servers map[string]server
		Clients clients
	}

	app struct {
		Debug bool
	}

	ownerInfo struct {
		Name string
		Org  string `toml:"organization"`
		Bio  string
		DOB  time.Time
	}

	database struct {
		Server  string
		Ports   []int
		ConnMax int `toml:"connection_max"`
		Enabled bool
	}

	server struct {
		IP string
		DC string
	}

	clients struct {
		Data  [][]interface{}
		Hosts []string
	}
)

func TestHasParent(t *testing.T) {
	s := cmdr.BaseOpt{
		Name:  "A",
		Short: "A",
		Full:  "Abcuse",
	}
	if s.HasParent() {
		t.Failed()
	}
	if s.GetTitleNames() != "A, Abcuse" {
		t.Failed()
	}
}

func TestSetGetStringSlice(t *testing.T) {
	cmdr.Set("A", []int{3, 7})

	oo := cmdr.GetStringSlice("app.A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSlice on int slice")
	}
	oo = cmdr.GetStringSliceR("A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSliceR on int slice")
	}

	cmdr.Set("A", "3,7")
	oo = cmdr.GetStringSlice("app.A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSlice on int slice")
	}
	oo = cmdr.GetStringSliceR("A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSliceR on int slice")
	}
}

func TestSetGetStringSlice2(t *testing.T) {
	cmdr.Set("A", []float32{3, 7})
	oo := cmdr.GetStringSlice("app.A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSlice on int slice")
	}
	oo = cmdr.GetStringSliceR("A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSliceR on int slice")
	}

	cmdr.Set("A", []byte("3,7"))
	oo = cmdr.GetStringSlice("app.A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSlice on int slice")
	}
	oo = cmdr.GetStringSliceR("A")
	if "3" != oo[0] || "7" != oo[1] {
		t.Fatal("wrong GetStringSliceR on int slice")
	}

	cmdr.Set("A", 99)
	oo = cmdr.GetStringSlice("app.A")
	if "99" != oo[0] {
		t.Fatal("wrong GetStringSlice on int slice")
	}
	oo = cmdr.GetStringSliceR("A")
	if "99" != oo[0] {
		t.Fatal("wrong GetStringSliceR on int slice")
	}
}

func TestSetGetIntSlice(t *testing.T) {
	// int slice

	cmdr.Set("A", []string{"3", "7"})
	oi := cmdr.GetIntSlice("app.A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSlice on int slice 1 ")
	}
	oi = cmdr.GetIntSliceR("A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSliceR on int slice 1 ")
	}

	cmdr.Set("A", []int{3, 7})
	oi = cmdr.GetIntSlice("app.A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSlice on int slice 1 ")
	}
	oi = cmdr.GetIntSliceR("A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSliceR on int slice 1 ")
	}

	cmdr.Set("A", "3,7")
	oi = cmdr.GetIntSlice("app.A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSlice on int slice 2")
	}
	oi = cmdr.GetIntSliceR("A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSliceR on int slice 2")
	}
}

func TestSetGetIntSlice2(t *testing.T) {
	// int slice

	cmdr.Set("A", []float32{3, 7})
	oi := cmdr.GetIntSlice("app.A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSlice on int slice 3")
	}
	oi = cmdr.GetIntSliceR("A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSliceR on int slice 3")
	}

	cmdr.Set("A", []byte("3,7"))
	oi = cmdr.GetIntSlice("app.A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSlice on int slice 4")
	}
	oi = cmdr.GetIntSliceR("A")
	if 3 != oi[0] || 7 != oi[1] {
		t.Fatal("wrong GetIntSliceR on int slice 4")
	}

	cmdr.Set("A", "99")
	oi = cmdr.GetIntSlice("app.A")
	if 99 != oi[0] {
		t.Fatal("wrong GetIntSlice on int slice 5")
	}
	oi = cmdr.GetIntSliceR("A")
	if 99 != oi[0] {
		t.Fatal("wrong GetIntSliceR on int slice 5")
	}

	cmdr.Set("A", 99)
	oi = cmdr.GetIntSlice("app.A")
	if 99 != oi[0] {
		t.Fatal("wrong GetIntSlice on int slice 5")
	}
	oi = cmdr.GetIntSliceR("A")
	if 99 != oi[0] {
		t.Fatal("wrong GetIntSliceR on int slice 5")
	}
}

var (
	tomlSample = []byte(`

runmode="devel"

title = "TOML Example"

[app]
debug=true

[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] 
# just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]

`)
)

func TestTomlLoad(t *testing.T) {
	var (
		err    error
		b      []byte
		mm     map[string]map[string]interface{}
		config tomlConfig
		meta   toml.MetaData
	)

	if err = ioutil.WriteFile(".tmp.toml", tomlSample, 0644); err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Remove(".tmp.toml")
		_ = os.Remove(".tmp.2.toml")
	}()

	mm = make(map[string]map[string]interface{})
	if err = toml.Unmarshal(b, &mm); err != nil {
		return
	}

	t.Log(mm)

	if meta, err = toml.DecodeFile(".tmp.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	t.Log(config)
	t.Log(meta)

	if err = cmdr.SaveObjAsToml(config, ".tmp.2.toml"); err != nil {
		t.Fatal(err)
	}

	// if err = cmdr.LoadConfigFile(".tmp.toml"); err != nil {
	// 	t.Fatal(err)
	// }

}

func TestConfigFiles(t *testing.T) {
	var err error

	cmdr.Set("no-watch-conf-dir", true)

	defer func() {
		_ = os.Remove(".tmp.json")
		_ = os.Remove(".tmp.yaml")
		_ = os.Remove(".tmp.toml")
	}()

	// try loading cfg again for gocov
	if _, _, err = cmdr.LoadConfigFile(".tmp.yaml"); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(".tmp.yaml")

	// try loading cfg again for gocov
	if _, _, err = cmdr.LoadConfigFile(".tmp.yaml"); err != nil {
		t.Fatal(err)
	}

	_ = ioutil.WriteFile(".tmp.yaml", []byte(`
app'x':"
`), 0644)

	// try loading cfg again for gocov
	if _, _, err = cmdr.LoadConfigFile(".tmp.yaml"); err == nil {
		t.Fatal("loading cfg file should be failed (err != nil), but it returns nil as err.")
	}
	_ = os.Remove(".tmp.yaml")

	_ = ioutil.WriteFile(".tmp.json", []byte(`{"app":{"debug":errrrr}}`), 0644)
	if _, _, err = cmdr.LoadConfigFile(".tmp.json"); err == nil {
		t.Fatal(err)
	}

	_ = ioutil.WriteFile(".tmp.json", []byte(`{"app":{"debug":false}}`), 0644)
	// try loading cfg again for gocov
	if _, _, err = cmdr.LoadConfigFile(".tmp.json"); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(".tmp.json")

	_ = ioutil.WriteFile(".tmp.toml", []byte(`
runmode=devel
`), 0644)
	if _, _, err = cmdr.LoadConfigFile(".tmp.toml"); err == nil {
		t.Fatal(err)
	}

	_, _, _ = cmdr.LoadConfigFile(".tmp.x.toml")

	_ = ioutil.WriteFile(".tmp.toml", []byte(`
runmode="devel"
[app]
debug=true
`), 0644)
	// try loading cfg again for gocov
	if _, _, err = cmdr.LoadConfigFile(".tmp.toml"); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(".tmp.toml")
}

func TestLaunchEditor2(t *testing.T) {
	if b, err := tool.LaunchEditorWith("cat", "/etc/passwd"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(b))
	}

	if _, err := tool.LaunchEditorWith("cat", "/etc/not-exists"); err != nil {
		// t.Fatal("should have an error return for non-exist file")
		t.Fatalf(`cmdr.LaunchEditorWith("cat", "/etc/not-exists") failed: %v`, err)
	}
}

func TestLaunch(t *testing.T) {
	_ = tool.Launch("ls")
	_ = os.Setenv("EDITOR", "ls")
	_, _ = tool.LaunchEditor("EDITOR")
}

func TestNormalizeDir(t *testing.T) {
	if cmdr.NormalizeDir("./a") != path.Join(cmdr.GetCurrentDir(), "./a") {
		t.Failed()
	}
	if cmdr.NormalizeDir("../a") != path.Join(cmdr.GetCurrentDir(), "../a") {
		t.Failed()
	}
	if cmdr.NormalizeDir("~/a") != path.Join(os.Getenv("HOME"), "a") {
		t.Failed()
	}
	if cmdr.NormalizeDir("v/a") != "v/a" {
		t.Failed()
	}
	_ = os.Setenv("EDITOR", "ls")
	_, _ = tool.LaunchEditor("EDITOR")

	_ = tool.Launch("ls", "/not-exists")

	// _ = cmdr.LaunchSudo("ls", "/not-exists")
}

func TestNoColorMode(t *testing.T) {
	cmdr.GetStrictMode()
	cmdr.GetDebugMode()
	cmdr.GetVerboseMode()
	cmdr.GetQuietMode()
	cmdr.GetNoColorMode()
}
