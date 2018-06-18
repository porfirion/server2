package network

type SearchableArray []uint64

func (arr SearchableArray) indexOf(value uint64) (bool, int) {
	if len(arr) == 0 {
		return false, -1
	}
	for ind, val := range arr {
		if val == value {
			return true, ind
		}
	}
	return false, -1
}
