package utils

func Ceiling(number int, significance int) int {
	if number%significance == 0 {
		return number
	}
	y := number / 500
	return 500 * (y + 1)
}
