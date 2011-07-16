package main

import "fmt"
import "strings"
import "io/ioutil"
import "sort"
// import "regexp"
// import "math"

type SuffixArray struct{
	i []int
	s []string
}

func NewSuffixArray(str string) (suffixArray *SuffixArray) {
	suffixArray = new(SuffixArray)
	suffixArray.s = strings.Split(str, " ", -1)
	suffixArray.i = make([]int, len(suffixArray.s))
	for idx, _ := range suffixArray.s {
		suffixArray.i[idx] = idx
	}
	sort.Sort(suffixArray)
	return
}

func (suffixArray *SuffixArray) Len() (int) {
	return len(suffixArray.i)
}

func (suffixArray *SuffixArray) Less(m, n int) (b bool) {
	switch {
		case suffixArray.s[suffixArray.i[m]] < suffixArray.s[suffixArray.i[n]] :
			b = true
			fmt.Println(1)
		case suffixArray.s[suffixArray.i[m]] > suffixArray.s[suffixArray.i[n]] :
			b = false
			fmt.Println(2)
		case suffixArray.s[suffixArray.i[m]] == suffixArray.s[suffixArray.i[n]] :
			if n == len(suffixArray.i)-1 {
				b = false
				fmt.Println(3)
			} else if m == len(suffixArray.i)-1 {
				b = true
				fmt.Println(4)
			} else {
				b = suffixArray.Less(m+1, n+1)
				fmt.Println(5)
			}
	}
	return b
}


func (suffixArray *SuffixArray) Swap(m, n int) {
	suffixArray.i[m], suffixArray.i[n] = suffixArray.i[n], suffixArray.i[m]
}

func main() {
	sourceFile,_ := ioutil.ReadFile("./corpora/test")
	sentences := strings.Split(string(sourceFile), "\n", -1)
	for _, element := range sentences {
		suffixArray := NewSuffixArray(element)
		fmt.Println(suffixArray)
	}
}