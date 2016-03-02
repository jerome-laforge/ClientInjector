package main

func Pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
