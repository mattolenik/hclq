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
	{`a = 12`,			`"12"`, 			Args{"get", "a"}},
	{`a = [12]`,		`["12"]`, 			Args{"get", "a[]"}},
	{`a = [12]`,		`"12"`, 			Args{"get", "a[0]"}},
	{`a = [1, 2, 3]`,	`["1","2","3"]`, 	Args{"get", "a[]"}},
	{`a = [1, 2, 3]`,	`"1"`, 				Args{"get", "a[0]"}},
	{`a = [1, 2, 3]`,	`"2"`, 				Args{"get", "a[1]"}},
	{`a = [1, 2, 3]`,	`"3"`, 				Args{"get", "a[2]"}},
}

func TestGet(t *testing.T) {
	assert := testifyAssert.New(t)
	for _, test := range getTests {
		out, err := run(test.input, test.args...)
		assert.Equal(test.expected, out, "args: %s", test.args)
		assert.NoError(err, "args: %s", test.args)
	}
}

func run(input string, args ...string) (string, error) {
	hclqBin := os.Getenv("HCLQ_BIN")
	cmd := exec.Command(hclqBin, args...)
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	out, err := cmd.Output()
	return string(out[:]), err
}