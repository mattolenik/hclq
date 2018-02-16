foo "baz" {
	key = 7
	arr = [4, 5]
}

foo "bar" {
	key = 12
	arr = [1, 2, 3]

	obj {
		a = "abc"
	}

	objAssign = {
		a = "def"
	}
}

a {
  b {
    c = true
  }
}

k = 10
