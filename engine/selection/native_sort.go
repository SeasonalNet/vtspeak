package selection

// paul2013SmallSortLimit is FUN_1001b5f0's insertion-sort cutoff. Its range
// check uses lastIndex-firstIndex > 16, so slices of at most 17 entries take
// the insertion path.
const paul2013SmallSortLimit = 17

// sortPaul2013Native ports FUN_1001b5f0's partition sort. It preserves the
// native median-of-three swaps, sentinel-bounded strict scans, 17-entry
// insertion cutoff, and 32:1 imbalance fallback to FUN_1001b400's heap sort.
// Those details determine equal-score permutations, so Go's standard sort is
// not an equivalent replacement here.
func sortPaul2013Native[T any](values []T, less func(left, right T) bool) {
	type interval struct{ first, last int }
	if len(values) < 2 {
		return
	}
	stack := make([]interval, 0, 32)
	first, last := 0, len(values)-1
	for {
		for last-first+1 > paul2013SmallSortLimit {
			middle := first + (last-first)/2
			if less(values[last], values[first]) {
				values[first], values[last] = values[last], values[first]
			}
			if less(values[middle], values[first]) {
				values[first], values[middle] = values[middle], values[first]
			}
			if less(values[last], values[middle]) {
				values[middle], values[last] = values[last], values[middle]
			}
			pivot := values[middle]
			left, right := first+1, last-1
			for left < right {
				for less(values[left], pivot) {
					left++
				}
				for less(pivot, values[right]) {
					right--
				}
				if left >= right {
					break
				}
				values[left], values[right] = values[right], values[left]
				left++
				right--
			}

			leftSize := right - first + 1
			rightSize := last - right
			if leftSize >= rightSize {
				if leftSize/32 > rightSize {
					sortPaul2013Heap(values[first:right+1], less)
					first = right + 1
				} else {
					stack = append(stack, interval{right + 1, last})
					last = right
				}
			} else if rightSize/32 > leftSize {
				sortPaul2013Heap(values[right+1:last+1], less)
				last = right
			} else {
				stack = append(stack, interval{first, right})
				first = right + 1
			}
		}
		sortPaul2013Small(values[first:last+1], less)
		if len(stack) == 0 {
			return
		}
		next := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		first, last = next.first, next.last
	}
}

// sortPaul2013Heap mirrors FUN_1001b400's one-based, max-heap path. The
// native helper descends through the larger child, then bubbles the saved
// value back up that path before each extraction.
func sortPaul2013Heap[T any](values []T, less func(left, right T) bool) {
	length := len(values)
	for root := length / 2; root >= 1; root-- {
		sortPaul2013HeapAdjust(values, root, length, less)
	}
	for end := length; end > 1; end-- {
		values[0], values[end-1] = values[end-1], values[0]
		sortPaul2013HeapAdjust(values[:end-1], 1, end-1, less)
	}
}

func sortPaul2013HeapAdjust[T any](values []T, root, length int, less func(left, right T) bool) {
	saved := values[root-1]
	hole := root
	for child := hole * 2; child <= length; child = hole * 2 {
		if child < length && less(values[child-1], values[child]) {
			child++
		}
		values[hole-1] = values[child-1]
		hole = child
	}
	for hole > root {
		parent := hole / 2
		if !less(values[parent-1], saved) {
			break
		}
		values[hole-1] = values[parent-1]
		hole = parent
	}
	values[hole-1] = saved
}

// sortPaul2013Small applies the native insertion path for short candidate
// lists. It shifts only strictly greater keys, preserving equal entries.
func sortPaul2013Small[T any](values []T, less func(left, right T) bool) {
	for index := 1; index < len(values); index++ {
		value := values[index]
		position := index
		for position > 0 && less(value, values[position-1]) {
			values[position] = values[position-1]
			position--
		}
		values[position] = value
	}
}
