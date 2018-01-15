package cmd

import (
    "testing"
    testifyAssert "github.com/stretchr/testify/assert"
    "os/exec"
    "io"
    "os"
    "strings"
)

type Args = []string

var getTests = []struct {
    input    string
    expected string
    args     Args
    err      error
}{
    {`a = 12`,          `"12"`,           Args{"get", "a"},    nil},
    {`a = [12]`,        `["12"]`,         Args{"get", "a[]"},  nil},
    {`a = [12]`,        `"12"`,           Args{"get", "a[0]"}, nil},
    {`a = [1, 2, 3]`,   `["1","2","3"]`,  Args{"get", "a[]"},  nil},
    {`a = [1, 2, 3]`,   `"1"`,            Args{"get", "a[0]"}, nil},
    {`a = [1, 2, 3]`,   `"2"`,            Args{"get", "a[1]"}, nil},
    {`a = [1, 2, 3]`,   `"3"`,            Args{"get", "a[2]"}, nil},
    {`a = []`,          `[]`,             Args{"get", "a[]"}, nil},
    {`a = []`,          `[]`,             Args{"get", "a[0]"}, nil},
}

func TestGet(t *testing.T) { for _, test := range getTests {
    t.Run(strings.Join(test.args, " "), func(tt *testing.T) {
        assert := testifyAssert.New(tt)

        cmd := exec.Command(os.Getenv("HCLQ_BIN"), test.args...)
        stdin, _ := cmd.StdinPipe()
        go func() {
            defer stdin.Close()
            io.WriteString(stdin, test.input)
        }()
        out, err := cmd.Output()
        if err != nil {
            tt.Fatal(err)
        }
        output := string(out[:])

        assert.Equal(test.expected, output)
        assert.NoError(err)
    })
}}