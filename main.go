package main

import "fmt"

func main() {
	fmt.Println(TwoSum([]int{2, 8, 7, 123}, 9))
}

func TwoSum(nums []int, target int) []int {
	hashT := make(map[int]int)
	for i, num := range nums { //i — индекс, num — значение
		diff := target - num
		prevIndex, isOk := hashT[diff]

		if isOk {
			return []int{prevIndex, i}
		}

		hashT[num] = i
	}

	return nil
}
