package util

// Ternary 三目运算的函数
func Ternary(a bool, b, c interface{}) interface{} {
	if a {
		return b
	}
	return c
}
