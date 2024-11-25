package algo

func RemoveFromSlice(slice, itemsToRemove []int) []int {
	// 创建一个 map 来存储需要删除的元素
	toRemove := make(map[int]bool)
	for _, item := range itemsToRemove {
		toRemove[item] = true
	}

	// 遍历原始 slice，将不在 toRemove 中的元素添加到新的 slice 中
	var result []int
	for _, item := range slice {
		if !toRemove[item] {
			result = append(result, item)
		}
	}
	return result
}
