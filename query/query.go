package query

import (
	"regexp"
	"strings"
)

// A segment in a query string. Given 'a.b.c', a b and c are query parts
type Node interface {
	IsMatch(value string) bool
}

type Key struct {
	value string
}

type Regex struct {
	pattern *regexp.Regexp
}

func (r *Regex) IsMatch(value string) bool {
	return r.pattern.FindString(value) != ""
}

func (k *Key) IsMatch(value string) bool {
	return k.value == value
}

var keyRegex, _ = regexp.Compile(`^([\w|-]+)`)

var regexRegex, _ = regexp.Compile(`^/((?:[^\\/]|\\.)*)/`)

func Parse(query string) ([]Node, error) {
	queryList := make([]Node, 0)
	query = strings.Trim(query, "\"'")
	err := parseQuery(query, 0, &queryList)
	return queryList, err
}

// .simple.struct.value
// .foo.bar.items[*]
func parseQuery(query string, i int, queue *[]Node) error {
	if i >= len(query) {
		return nil
	}
	char := query[i: i+1]
	if char == "." {
		return parseQuery(query, i+1, queue)
	}
	regexMatches := regexRegex.FindStringSubmatch(query[i:])
	if len(regexMatches) > 1 {
		pattern, err := regexp.Compile(regexMatches[1])
		if err != nil {
			return err
		}
		newNode := &Regex{
			pattern: pattern,
		}
		*queue = append(*queue, newNode)
		i += len(regexMatches[0])
		return parseQuery(query, i, queue)
	}
	key := keyRegex.FindString(query[i:])
	if key != "" {
		i += len(key)
		newNode := &Key{
			value: key,
		}
		*queue = append(*queue, newNode)
		if i >= len(query) {
			return nil
		}
		return parseQuery(query, i, queue)
	}
	return nil
}
