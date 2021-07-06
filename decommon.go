
import "sort"

// decommon will remove common elements from two slices and return a sorted copy of unique elements
func decommon(ai, bi []int) (ao, bo []int) {
	ao = append([]int(nil), ai...) // copy
	bo = append([]int(nil), bi...) // copy
	sort.Ints(ao)
	sort.Ints(bo)

	var i, j int
	for i < len(ao) && j < len(bo) {
		if ao[i] == bo[j] {
			ao = append(ao[:i], a[i+1:]...) // remove i
			bo = append(bo[:j], a[j+1:]...) // remove j
			continue
		}
		if ao[i] < bo[j] {
			i++
			continue
		}
		j++
	}
}