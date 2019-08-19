package main

import (
	"fmt"
	"math"
)

//1 = I
//5 = V
//10 = X
//50 = L
//100 = C
//500 = D
//1000 = M
//var roman_list = [7]string{"I", "V", "X", "L", "C", "D", "M"}
var roman_list = [5]string{"I", "V", "X", "L", "C"}
const const_max_num = 100
func main() {

	for i := 1; i <= const_max_num; i++ {
		fmt.Println("Value = ", i, " Roman = ", convert2roman(i))
	}
}

func convert2roman(x int) string {
	//	var remainder, quotient float64
	var count int
	result := ""
	
	if x == const_max_num {
		return roman_list[len(roman_list)-1]
	}
	quotient := float64(x)
	for ; quotient > 0; count += 2 {
		temp_result :=""
		remainder := math.Mod(quotient, 10)
		quotient = (quotient - remainder) / 10
		for i := 1; float64(i) <= remainder; i++ {
			switch {
			case i == 4:
				temp_result = roman_list[count] + roman_list[count+1]
			case i == 5:
				temp_result = roman_list[count+1]
			case i == 9:
				temp_result = roman_list[count] + roman_list[count+2]
			default:
				temp_result += roman_list[count]
			}
		}
		result = temp_result + result
	}
	return result
}
