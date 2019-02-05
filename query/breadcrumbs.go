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

// Crumb represents an individual portion of a user's query. For example, given
// the query 'data.*.bar.id', each of data, *, bar, and id are each crumbs.
// Crumbs implement IsMatch to determine whether or not it can be used to match
// a given HCL key.
type Crumb interface {
	IsMatch(key string, val ast.Node) (bool, error)
	Key() string
}

// IndexedCrumb is a type of Crumb used for indexable elements such as arrays
// or lists. Its index is optional and will be nil if no index was specified.
// For example, in the query 'foo.bar[1]', the bar[1] portion will be represented
// by an IndexedCrumb with an index of 1. If an index is not present, such as
// with 'bar[]', the index will be nil.
type IndexedCrumb interface {
	IsMatch(key string, val ast.Node) (bool, error)
	Key() string
	Index() *int
}

// Key is a literal breadcrumb that matches based on its exact value.
type Key struct {
	value string
}

// List is a breadcrumb representing a list or array, either the entire list or
// just an individual item.
type List struct {
	value string
	index *int
	key   string
}

// Regex is a breadcrumb that matches based on a regex. It can take an optioanl indexer.
type Regex struct {
	pattern *regexp.Regexp
	index   *int
}

// Wildcard is the literal breadcrumb '*' that matches anything.
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

// Matches a list `abc[]` or `abc[123]` or `a[-1]`
var listRegex, _ = regexp.Compile(`^([\w|-]+)\[(-?\d*)]`)

// Matches by regex `/someRegex/` with optional indexer, e.g. `/someRegex/[]`
var regexRegex, _ = regexp.Compile(`/((?:[^\\/]|\\.)*)/(\[(\d*)\])?`)

// ParseBreadcrumbs reads in a query string specified by the user and breaks it
// down into Crumb instances that can be matched against HCL keys with IsMatch.
func ParseBreadcrumbs(queryString string) (*Breadcrumbs, error) {
	queryString = strings.Trim(queryString, `"'`)
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
