package main

import "fmt"
import "textprocessing"


func main() {
	input := textprocessing.NewText("a b r a c a d a b r a")
	suffixArray := textprocessing.NewSuffixArray(input)
	fmt.Println(suffixArray)
}