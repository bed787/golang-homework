package main

import "fmt"

func main() {
	fmt.Println(twoSum([]int{2, 8, 7, 123}, 9))
}

func twoSum(nums []int, target int) []int {
	hashT := make(map[int]int)
	for i, num := range nums {
		diff := target - num
		prevIndex, isOk := hashT[diff]

		if isOk {
			return []int{prevIndex, i}
		}

		hashT[num] = i
	}

	return nil
}
