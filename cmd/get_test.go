package cmd


import (
	"testing"
	testifyAssert "github.com/stretchr/testify/assert"
	"os/exec"
	"io"
	"os"
)

type Args = []string

var getTests = []struct {
	input    string
	expected string
	args     Args
}{
	{`foo = 12`,		`"12"`, 	Args{"get", "foo"}},
	{`foo = [12]`,		`["12"]`, 	Args{"get", "'foo[]'"}},
}

func TestGet(t *testing.T) {
	assert := testifyAssert.New(t)
	for _, test := range getTests {
		out, err := run(test.input, test.args...)
		assert.Equal(test.expected, out)
		assert.NoError(err)
	}
}

func run(input string, args ...string) (string, error) {
	hclqBin := os.Getenv("HCLQ_BIN")
	//dlvBin := os.Getenv("DLV_BIN")
	//args = append([]string{"--listen=:2345", "--headless=true", "--api-version=2", "exec", hclqBin}, args...)
	//cmd := exec.Command(dlvBin, args...)
	cmd := exec.Command(hclqBin, args...)
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	out, err := cmd.Output()
	return string(out[:]), err
}