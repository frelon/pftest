package pkg

type TakeFunc func(token string) bool

func Take(tokens []string, f TakeFunc) []string {
	for i, t := range tokens {
		if f(t) {
			return tokens[i+1:]
		}
	}

	return []string{}
}

func TakeUntil(tokens []string, f TakeFunc) []string {
	for i, t := range tokens {
		if f(t) {
			return tokens[i+1:]
		}
	}

	return []string{}
}
