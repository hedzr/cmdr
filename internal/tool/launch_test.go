package tool

import (
	"os"
	"testing"
)

func TestRandomFileName(t *testing.T) {
	fn := shellEditorRandomFilename()
	t.Logf("fn = %v", fn)
}

func TestTempFileName(t *testing.T) {
	fn := TempFileName("message*.tmp", "message001.tmp")
	t.Logf("fn = %v", fn)
}

func TestLaunchEditorWithGetter(t *testing.T) {
	logfile := "/tmp/1.log"
	str := []byte("hello world")
	os.WriteFile(logfile, str, 0644)
	content, err := LaunchEditorWithGetter("echo", func() string { return logfile }, false)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(str) {
		t.Fatal("content != string(str)")
	}
}
