package util

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	return -MaxInt(-a, -b)
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	return -MaxInt64(-a, -b)
}
