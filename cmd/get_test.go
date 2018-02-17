package cmd

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

type Args = []string

var getTests = []struct {
	input    string
	expected string
	args     Args
	errText  string
}{
	{`a = 12`, `12`, Args{"get", "a"}, ""},
	{`a = [12]`, `[12]`, Args{"get", "a[]"}, ""},
	{`a = [12]`, `12`, Args{"get", "a[0]"}, ""},
	{`a = [1, 2, 3]`, `[1,2,3]`, Args{"get", "a[]"}, ""},
	{`a = [1, 2, 3]`, `1`, Args{"get", "a[0]"}, ""},
	{`a = [1, 2, 3]`, `2`, Args{"get", "a[1]"}, ""},
	{`a = [1, 2, 3]`, `3`, Args{"get", "a[2]"}, ""},
	{`a = []`, `[]`, Args{"get", "a[]"}, ""},
	{`a { b = "2" }`, `"2"`, Args{"get", "a.b"}, ""},
	{`a { b = "2" }`, `"2"`, Args{"get", "-r", "a.b"}, ""},
	{`a { b = "2a" }`, `"2a"`, Args{"get", "-r", "a.b"}, ""},
	{`a "b" { c = [1] }`, `[1]`, Args{"get", "a.b.c[]"}, ""},
}

func TestGet(t *testing.T) {
	for _, test := range getTests {
		t.Run(strings.Join(test.args, " "), func(tt *testing.T) {
			assert := testifyAssert.New(tt)

			cmd := exec.Command(os.Getenv("HCLQ_BIN"), test.args...)
			stdin, _ := cmd.StdinPipe()
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, test.input)
			}()
			outBytes, err := cmd.Output()
			output := string(outBytes[:])
			if test.errText != "" {
				err, ok := err.(*exec.ExitError)
				if !ok {
					tt.Fatalf("Expected ExitError, got %+v", err)
				}
				stderr := string(err.Stderr)
				assert.Contains(stderr, test.errText)
			} else {
				assert.NoError(err)
			}
			assert.Equal(test.expected, output)
		})
	}
}
