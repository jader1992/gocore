package util

import "fmt"

// PrettyPrint 美观输出数组
func PrettyPrint(arr [][]string) {
	if len(arr) == 0 {
		return
	}
	rows := len(arr)
	cols := len(arr[0])

	// 统计二维数组每个元素的长度
	lens := make([][]int, rows)
	for i := 0; i < rows; i++ {
		lens[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			lens[i][j] = len(arr[i][j])
		}
	}

	// 统计每一列最长长度
	colMaxs := make([]int, cols)
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			if colMaxs[j] < lens[i][j] {
				colMaxs[j] = lens[i][j]
			}
		}
	}

	// 格式化输出
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Print(arr[i][j])
			padding := colMaxs[j] - lens[i][j] + 2 // 补充的空白数
			for p := 0; p < padding; p++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}