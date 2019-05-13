/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
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

	b, err := yaml.Marshal(demoOpts)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("demo-opts.yaml", b, 0644)
	if err != nil {
		t.Fatal(err)
	}

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
