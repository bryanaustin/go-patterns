
import "sort"

// dedup will remove duplicates and return a sorted copy
func dedup(input []int) (output []int) {
	if len(inputs) < 1 {
		// Prevent index 0 error
		return
	}

	output = make([]int, 1, len(input))
	sortedinput := append([]int(nil), input...) //copy
	sort.Ints(sortedinput)
	output[0] = sortedinput[0]

	for i := 1; i < len(sortedinput); i++ {
		if sortedinput[i-1] == sortedinput[i] {
			continue
		}
		output = append(output, sortedinput[i])
	}
	return
}
