package main

import "fmt"
import "strings"
import "sort"
// import "regexp"
// import "math"

type SuffixArray struct{
        i []int
        s []string
}

func NewSuffixArray(str ... string) (suffixArray *SuffixArray) {
        suffixArray = new(SuffixArray)
        suffixArray.s=str
        suffixArray.i = make([]int, len(str))
        for idx := range suffixArray.i {
                suffixArray.i[idx] = idx
        }
        sort.Sort(suffixArray)
        return
}

func (suffixArray *SuffixArray) Len() (int) {
        return len(suffixArray.i)
}

func (suffixArray *SuffixArray) Less(m, n int) (bool) {
        midx:=suffixArray.i[m]
        nidx:=suffixArray.i[n]
        for{
                switch {
                        case suffixArray.s[midx]<suffixArray.s[nidx] :
                                return true
                        case suffixArray.s[midx]>suffixArray.s[nidx] :
                                return false
                        case suffixArray.s[midx]==suffixArray.s[nidx] :
                                midx++
                                nidx++
                                if nidx == len(suffixArray.i) {
                                        return false
                                } else if midx == len(suffixArray.i) {
                                        return true
                                }
                }
        }
        panic("unreachable")
}


func (suffixArray *SuffixArray) Swap(m, n int) {
        suffixArray.i[m], suffixArray.i[n] = suffixArray.i[n], suffixArray.i[m]
}

func (suffixArray *SuffixArray) String() (string){
        array := make([]string, len(suffixArray.i))
        for idx, element := range suffixArray.i {
                array[idx] = strings.Join(suffixArray.s[element:], " ")
        }
        return strings.Join(array, "\n")
}

func main() {
                suffixArray := NewSuffixArray("a","b","r","a","c","a","d","a","b","r","a")
                fmt.Println(suffixArray)
}
