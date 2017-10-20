package main

// TODO: rename
type QueryNode interface {
	Value() string
}

type Key struct {
	value string
}

func(l *Key) Value() string {
	return l.value
}
