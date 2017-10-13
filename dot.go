package main

type dot struct {
	tokenStr string
	nxt queryNode
}

func(l *dot) value() string {
	return l.tokenStr
}

func(l *dot) setValue(value string) {
	l.tokenStr = value
}

func(l *dot) token() string {
	return l.tokenStr
}

func(l *dot) setToken(value string) {
	l.tokenStr = value
}

func(l *dot) next() queryNode {
	return l.nxt
}

func(l *dot) setNext(next queryNode) {
	l.nxt = next
}