package main

type root struct {
	tokenStr string
	nextNode queryNode
}

func(l *root) value() string {
	return l.tokenStr
}

func(l *root) setValue(value string) {
	l.tokenStr = value
}

func(l *root) token() string {
	return l.tokenStr
}

func(l *root) setToken(value string) {
	l.tokenStr = value
}

func(l *root) next() queryNode {
	return l.nextNode
}

func(l *root) setNext(next queryNode) {
	l.nextNode = next
}

func(l root) String() string {
	var node queryNode
	node = &l
	result := ""
	for node != nil {
		result += "." + node.token()
		node = node.next()
	}
	return result
}
