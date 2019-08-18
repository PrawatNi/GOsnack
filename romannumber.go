package main

import (
	"fmt"
	"math"
)

func main() {

	for i:= 1; i<= 100; i++ {
		fmt.Println("Value = ", i , " Roman = ", convert2roman(i))
	} 
}

func convert2roman(x int) string{
	result := ""
	y := math.Mod(float64(x), 10)
	z := (float64(x)-y)/10
	
	if x == 100 {
		return "C"
	}
	
	//The ten
	for z > 0 {
		switch {
			case z == 9 :
				result += "XC"
				z -= 9
			case z >= 5:
				result += "L"
				z -= 5
			case z >=	4:
				result += "XL"
				z -= 4	
			default:
				result += "X"
				z -= 1	
		}
	}

	//The digit
	for y > 0 {	
		switch {
			case y ==	9:
				result += "IX"
				y -= 9
			case y >=	5:
				result += "V"
				y -= 5	
			case y == 4:
				result += "IV"
				y -= 4	
			default:
				result += "I"
				y -= 1
		}
	}
	return result
}