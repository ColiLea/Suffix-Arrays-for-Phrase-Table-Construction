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
	suffixArray := textprocessing.NewSuffixArray("a","b","r","a","c","a","d","a","b","r","a")
	fmt.Println(suffixArray)
	fmt.Println(suffixArray.SuffixOfLength(4))
}