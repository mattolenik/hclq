package main

type literal struct {
	tok      string
	val      string
	nextNode queryNode
}

func(l *literal) value() string {
	return l.val
}

func(l *literal) setValue(value string) {
	l.val = value
}

func(l *literal) token() string {
	return l.tok
}

func(l *literal) setToken(value string) {
	l.tok = value
}

func(l *literal) next() queryNode {
	return l.nextNode
}

func(l *literal) setNext(next queryNode) {
	l.nextNode = next
}
