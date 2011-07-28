package main

import "fmt"
import "textprocessing"


func main() {
// 	sourceFile,_ := ioutil.ReadFile("../corpora/test")
// 	sentences := strings.Split(string(sourceFile), "\n", -1)
// 	for _, element := range sentences {
// 		words := strings.Split(element, " ", -1)
// 		suffixArray := NewSuffixArray(words)
// 		fmt.Println(suffixArray)
// 	}
	input := textprocessing.NewText("a b r a c a d a b r a")
	suffixArray := textprocessing.NewSuffixArray(input)
	fmt.Println(suffixArray)
}