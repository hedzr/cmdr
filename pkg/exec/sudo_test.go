// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestRun(t *testing.T) {
	_ = Run("ls", "-a")
	_, _, _ = Sudo("ls", "-a")
}

func TestRunWithOutput(t *testing.T) {
	_, out, err := RunWithOutput("ls", "-la", "/not-exits")
	t.Logf("stdout: %v", out)
	t.Logf("stderr: %v", err)
}

func TestRunCommand(t *testing.T) {
	cmd := exec.Command("ls", "-lah", "/not-exits")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Logf("cmd.Run() failed with %s\n", err)
	}
}

func TestCombineOutputs(t *testing.T) {
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd.Run() failed with %s\n%s", err, out)
	}
	t.Logf("combined out:\n%s\n", string(out))
}

func TestAlone(t *testing.T) {
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var err error
	cmd := exec.Command("ls", "-lah", "/not-exits")
	stderr, err = cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	// var errBuffer bytes.Buffer
	// cmd.Stderr = &errBuffer
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	// fmt.Printf("%s\n", slurp)

	err = cmd.Wait()
	if err != nil {
		// var sb bytes.Buffer
		// in := bufio.NewScanner(stderr)
		// for in.Scan() {
		//	sb.Write(in.Bytes())
		// }
		t.Logf("cmd.Run() failed with %s\n%s", err, slurp)
	}
	out, _ := ioutil.ReadAll(stdout)
	t.Logf(" output:\n%s\n", string(out))
}
