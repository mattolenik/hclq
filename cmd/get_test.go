package cmd


import (
	"testing"
	testifyAssert "github.com/stretchr/testify/assert"
	"os/exec"
	"io"
)

func TestGet(t *testing.T) {
	assert := testifyAssert.New(t)
	out, err := run("foo = 12", "get", "foo")
	assert.Equal("\"12\"", out)
	assert.NoError(err)
}

func run(input string, args ...string) (string, error) {
	cmd := exec.Command("./hclq", args...)
	cmd.Dir = "../dist"
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	out, err := cmd.Output()
	return string(out[:]), err
}