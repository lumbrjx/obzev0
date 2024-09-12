package main

import (
	"fmt"
	"strconv"
)

func strToFlt(str string) float32 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return 0
	}
	return float32(f)
}
