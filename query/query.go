package query

import (
	"regexp"
	"strings"
	"strconv"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Node interface {
	IsMatch(key string, val ast.Node) bool
}

type Key struct {
	value string
}

type List struct {
	Value string
	Key string
	Index *int
}

type Regex struct {
	pattern *regexp.Regexp
}

func (r *Regex) IsMatch(key string, val ast.Node) bool {
	return r.pattern.MatchString(key)
}

func (k *Key) IsMatch(key string, val ast.Node) bool {
	return k.value == key
}

func (l *List) IsMatch(key string, val ast.Node) bool {
	_, ok := val.(*ast.ListType)
	return ok && key == l.Key
}

// Matches by key literal `abc`
var keyRegex, _ = regexp.Compile(`^([\w|-]+)`)

// Matches a list `abc[]`
var listRegex, _ = regexp.Compile(`^([\w|-]+)\[(\d*)]`)

// Matches by regex
var regexRegex, _ = regexp.Compile(`^/((?:[^\\/]|\\.)*)/`)

func Parse(query string) ([]Node, error) {
	queryList := make([]Node, 0)
	query = strings.Trim(query, "\"'")
	err := parseQuery(query, 0, &queryList)
	return queryList, err
}

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
	listMatches := listRegex.FindStringSubmatch(query[i:])
	if listMatches != nil {
		list := listMatches[0]
		i += len(list)
		newNode := &List{
			Value: list,
			Key: listMatches[1],
		}
		index, err := strconv.Atoi(listMatches[2])
		if err == nil {
			newNode.Index = &index
		}
		*queue = append(*queue, newNode)
		if i >= len(query) {
			return nil
		}
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