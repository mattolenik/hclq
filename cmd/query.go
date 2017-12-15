package cmd

import "regexp"

// A segment in a query string. Given 'a.b.c', a b and c are query parts
type Node interface {
	Value() string
}

type Key struct {
	value string
}

func (l *Key) Value() string {
	return l.value
}

var keyRegex, _ = regexp.Compile(`(\w+)`)

func Parse(query string) []Node {
	queryList := make([]Node, 0)
	parseQuery(query, 0, &queryList)
	return queryList
}

// .simple.struct.value
// .foo.bar.items[*]
func parseQuery(query string, i int, queue *[]Node) {
	if i >= len(query) {
		return
	}
	char := query[i : i+1]
	if char == "." {
		parseQuery(query, i+1, queue)
		return
	}
	key := keyRegex.FindString(query[i:])
	if key != "" {
		i += len(key)
		newNode := &Key{
			value: key,
		}
		*queue = append(*queue, newNode)
		if i >= len(query) {
			return
		}
		parseQuery(query, i, queue)
		return
	}
}
