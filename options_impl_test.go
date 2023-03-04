// Copyright © 2023 Hedzr Yeh.

package cmdr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hedzr/evendeep"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
)

func TestSaveCheckpoint(t *testing.T) {

	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()
	cmdr.RemoveAllOnConfigLoadedListeners()
	cmdr.StopExitingChannelForFsWatcherAlways()

	t.Run(`save options as checkpoint`, func(t *testing.T) {
		defer prepareConfD(t)()

		var err error
		os.Args = []string{"consul-tags", "--config", "./conf.d"}
		conf.AppName = `consul-tags`
		// cmdr.SetInternalOutputStreams(nil, nil)
		if err = cmdr.Exec(rootCmdForTesting(), cmdr.WithNoWatchConfigFiles(true), cmdr.WithInternalOutputStreams(nil, nil)); err != nil {
			t.Fatal(err)
		}

		fmt.Println(cmdr.DumpAsString())

		if err = cmdr.SaveCheckpoint(); err != nil {
			t.Fatal(err)
		}

		t.Log("[XXX] cmdr.SaveCheckpoint() ENDED ---------")

		if cnt := cmdr.CheckpointSize(); cnt != 1 {
			t.Fatalf(`expecting CheckpointSize() returns 1 but got %d`, cnt)
		}

		t.Log("[XXX] cmdr.CheckpointSize() is ok ---------\n\n\n\n")

		aMap := map[string][]int{"v1": []int{3, 6}, "v2": []int{9}}
		if err = cmdr.MergeWith(map[string]interface{}{
			"app": map[string]interface{}{
				conf.AppName: map[string]interface{}{
					"a-map": aMap,
				},
			},
		}); err != nil {
			t.Fatal(err)
		}

		t.Log("[XXX] cmdr.MergeWith() ENDED ---------\n\n\n\n")

		m := cmdr.GetMapR(conf.AppName)
		if diff, equal := evendeep.DeepDiff(m["a-map"], aMap); !equal {
			t.Fatalf(`expecting fetch the new appended option ok, but they are different. Diff: %+v`, diff)
		}

		t.Log("[XXX] cmdr.GetMapR() ENDED ---------\n\n\n\n")

		cmdr.ResetOptions()

		t.Log("[XXX] cmdr.ResetOptions() ENDED ---------\n\n\n\n")

		if err = cmdr.RestoreCheckpoint(); err != nil {
			t.Fatal(err)
		}

		t.Log("[XXX] cmdr.RestoreCheckpoint() ENDED ---------\n\n\n\n")

		if cnt := cmdr.CheckpointSize(); cnt != 0 {
			t.Fatalf(`expecting CheckpointSize() returns 0 but got %d`, cnt)
		}

		m = cmdr.GetMapR(conf.AppName)
		if m != nil {
			t.Fatalf(`expecting m is nil, but got: %v`, m)
		}

		m = cmdr.GetMapR("ms.tags")
		if m == nil {
			t.Fatalf(`expecting the options had been restored but it's empty noe. cmdr.options:\n%v`, cmdr.DumpAsString())
		}
		if v := m["float"]; v != float32(357) {
			t.Fatalf(`expecting the options subkey app.ms.tags.clear = false, but it's %+v.`, v)
		}

		resetOsArgs()
		cmdr.ResetOptions()
	})

}
