package algorithm

func FindOptimalTry(num int) int {
	left := 0
	right := num
	attempts := 0
	for left <= right {
		attempts++
		temp := left + (right-left)/2
		if temp == right {
			return attempts
		} else if temp < num {
			left = temp + 1
		} else {
			right = right - 1
		}
	}
	return attempts
}
