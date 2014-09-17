package gisp

func MinInts(data ...Int) Int {
	if len(data) == 0 {
		panic("Min Ints except ints more than zero.")
	}
	min := data[0]
	for _, item := range data[1:] {
		if item < min {
			min = item
		}
	}
	return min
}

func MaxInts(data ...Int) Int {
	if len(data) == 0 {
		panic("Max Ints except ints more than zero.")
	}
	max := data[0]
	for _, item := range data[1:] {
		if max < item {
			max = item
		}
	}
	return max
}
