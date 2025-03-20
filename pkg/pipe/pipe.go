package pipe

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// PipeToReader do Shell pipe operation to a target program.
//
// For example:
//
//	err = PipeToReader(func(stdout io.Writer){
//	    fmt.FPrintf(stdout, "hello, world")
//	}, "less", "-R")
func PipeToReader(cb func(stdout io.Writer), program string, args ...string) (err error) {
	return (&pipeS{}).pipeToReader(cb, program, args...)
}

type pipeS struct{}

func (w *pipeS) pipeToReader(cb func(stdout io.Writer), program string, args ...string) (err error) {
	fname, err := exec.LookPath(program)
	if err == nil {
		program, err = filepath.Abs(fname)
	}
	if err != nil {
		return
	}
	cmd := exec.Command(program, args...)
	// cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var in io.WriteCloser
	in, err = cmd.StdinPipe()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer in.Close()
		defer wg.Done()
		// time.Sleep(time.Second * 1)
		cb(in)
	}()
	go func() {
		defer wg.Done()
		err = cmd.Run()
	}()
	wg.Wait()
	return
}

type ws struct {
	io.Writer
}

func (s *ws) WriteString(str string) (n int, err error) {
	n, err = s.Writer.Write([]byte(str))
	return
}
