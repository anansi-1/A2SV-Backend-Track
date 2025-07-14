package main

import (
	"fmt"
)

func main() {
	input1 := "hello there"
	input2 := "Never a foot too far, even."

	fmt.Println("Input 1:", input1)
	fmt.Println("Word Frequency:", WordFrequencyCounter(input1))
	fmt.Println("Is Palindrome:", PalindromeCheck(input1))
	fmt.Println()

	fmt.Println("Input 2:", input2)
	fmt.Println("Word Frequency:", WordFrequencyCounter(input2))
	fmt.Println("Is Palindrome:", PalindromeCheck(input2))
}
