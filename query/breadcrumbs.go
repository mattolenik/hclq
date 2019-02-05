package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Breadcrumbs represents the dot.style.query that the user passes in.
type Breadcrumbs struct {
	Parts  []Crumb
	Length int
}

func (q *Breadcrumbs) Slice(low int) *Breadcrumbs {
	return &Breadcrumbs{Parts: q.Parts[low:], Length: q.Length}
}

type Crumb interface {
	IsMatch(key string, val ast.Node) (bool, error)
	Key() string
}

type IndexedCrumb interface {
	IsMatch(key string, val ast.Node) (bool, error)
	Key() string
	Index() *int
}

type Key struct {
	value string
}

type List struct {
	value string
	index *int
	key   string
}

type Regex struct {
	pattern *regexp.Regexp
	index   *int
}

type Wildcard struct {
}

func (w *Wildcard) IsMatch(key string, val ast.Node) (bool, error) {
	return true, nil
}

func (w *Wildcard) Key() string {
	return "*"
}

func (r *Regex) IsMatch(key string, val ast.Node) (bool, error) {
	return r.pattern.MatchString(key), nil
}

func (r *Regex) Key() string {
	return ""
}

func (r *Regex) Index() *int {
	return r.index
}

func (k *Key) IsMatch(key string, val ast.Node) (bool, error) {
	if key == k.value {
		_, isList := val.(*ast.ListType)
		if isList {
			return true, fmt.Errorf("key '%s' found but is of wrong type, query requested key/literal, found list", key)
		}
		return true, nil
	}
	return false, nil
}

func (k *Key) Key() string {
	return k.value
}

func (l *List) IsMatch(key string, val ast.Node) (bool, error) {
	if key == l.Key() {
		_, ok := val.(*ast.ListType)
		if !ok {
			return false, fmt.Errorf("key '%s' found but is of wrong type, query requested list", key)
		}
		return true, nil
	}
	return false, nil
}

func (l *List) Key() string {
	return l.key
}

func (l *List) Index() *int {
	return l.index
}

// Matches by key literal `abc`
var keyRegex, _ = regexp.Compile(`^([\w|-]+)`)

// Matches a list `abc[]` or `abc[123]`
var listRegex, _ = regexp.Compile(`^([\w|-]+)\[(-?\d*)]`)

// Matches by regex `/someRegex/` with optional indexer, e.g. `/someRegex/[]`
var regexRegex, _ = regexp.Compile(`/((?:[^\\/]|\\.)*)/\[(\d*)]?`)

func ParseBreadcrumbs(queryString string) (*Breadcrumbs, error) {
	queryString = strings.Trim(queryString, "\"'")
	query := &Breadcrumbs{Parts: []Crumb{}}
	err := parseBreadcrumbs(queryString, 0, &query.Parts)
	query.Length = len(query.Parts)
	return query, err
}

func parseBreadcrumbs(query string, i int, queue *[]Crumb) error {
	if i >= len(query) {
		return nil
	}
	char := query[i : i+1]
	if char == "." {
		return parseBreadcrumbs(query, i+1, queue)
	}
	if char == "*" {
		newCrumb := &Wildcard{}
		*queue = append(*queue, newCrumb)
		return parseBreadcrumbs(query, i+1, queue)
	}
	regexMatches := regexRegex.FindStringSubmatch(query[i:])
	if len(regexMatches) > 1 {
		pattern, err := regexp.Compile(regexMatches[1])
		if err != nil {
			return err
		}
		newCrumb := &Regex{
			pattern: pattern,
		}
		index, err := strconv.Atoi(regexMatches[1])
		if err == nil {
			newCrumb.index = &index
		}
		*queue = append(*queue, newCrumb)
		i += len(regexMatches[0])
		return parseBreadcrumbs(query, i, queue)
	}
	listMatches := listRegex.FindStringSubmatch(query[i:])
	if listMatches != nil {
		list := listMatches[0]
		i += len(list)
		newCrumb := &List{
			value: list,
			key:   listMatches[1],
		}
		index, err := strconv.Atoi(listMatches[2])
		if err == nil {
			newCrumb.index = &index
		}
		*queue = append(*queue, newCrumb)
		if i >= len(query) {
			return nil
		}
		return parseBreadcrumbs(query, i, queue)
	}
	key := keyRegex.FindString(query[i:])
	if key != "" {
		i += len(key)
		newCrumb := &Key{
			value: key,
		}
		*queue = append(*queue, newCrumb)
		if i >= len(query) {
			return nil
		}
		return parseBreadcrumbs(query, i, queue)
	}
	return nil
}
