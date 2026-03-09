package main

func quicksort(arr []int, low int, high int) {
	if low < high {
		p := partition(arr, low, high)
		quicksort(arr, low, p-1)
		quicksort(arr, p+1, high)
	}
}

func partition(arr []int, low int, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}

func sortArray(arr []int) []int {
	if len(arr) <= 1 {
		return append([]int(nil), arr...)
	}

	sorted := append([]int(nil), arr...)
	quicksort(sorted, 0, len(sorted)-1)
	return sorted
}
