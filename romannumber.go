package main

import (
	"fmt"
	"math"
)

func main() {

	for i:= 0; i<= 100; i++ {
		fmt.Println("Value = ", i , " Roman = ", convert2roman(i))
	} 
}

func convert2roman(x int) string{
	
	y := math.Mod(float64(x), 10)
	z := (float64(x)-y)/10
	a := ""
	b := ""
	
	switch z {
		case 1:
			a = "X"
		case 2:
			a = "XX"		
		case 3:
			a = "XXX"
		case 4:
			a = "XL"
		case 5:
			a = "L"
		case 6:
			a = "LX"		
		case 7:
			a = "LXX"
		case 8:
			a = "LXXX"
		case 9:
			a = "XC"				
		case 10:
			a = "C"					
	}
	switch y {
		case 1:
			b = "I"
		case 2:
			b = "II"		
		case 3:
			b = "III"
		case 4:
			b = "IV"
		case 5:
			b = "V"
		case 6:
			b = "VI"		
		case 7:
			b = "VII"
		case 8:
			b = "VIII"
		case 9:
			b = "IV"				
		case 10:
			b = "X"					
	}
		
	return a+b
}