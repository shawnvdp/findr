package main

func Max(numbers ...int) int {
	max := numbers[0]

	for _, number := range numbers[1:] {
		if max >= number {
			continue
		}
		max = number
	}

	return max
}

func Min(numbers ...int) int {
	min := numbers[0]

	for _, number := range numbers[1:] {
		if min <= number {
			continue
		}
		min = number
	}

	return min
}
