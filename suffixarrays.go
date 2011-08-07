package main

import "fmt"
import "textprocessing"
import "io/ioutil"
import "strings"

func main() {
	var sentenceIndices []int
	var sourceText *textprocessing.Text = &textprocessing.Text{}
	{
		rawText, _ := ioutil.ReadFile("/home/lea/corpora/test")
		sentences := strings.Split(string(rawText), "\n", -1)
		sentenceTexts := make([]textprocessing.Text, len(sentences))
		corpusLength := 0
		for idx, sentence := range sentences {
			sentenceTexts[idx] = *textprocessing.NewText(sentence)
			corpusLength += sentenceTexts[idx].Length()
		}
		sentenceIndices = make([]int, 0, corpusLength)
		input := make([]textprocessing.Word, 0, corpusLength)
		for idx, sentence := range sentenceTexts {
			for ii := 0; ii < sentence.Length(); ii++ {
				sentenceIndices = append(sentenceIndices, idx)
			}
			input = append(input, sentence...)
		}
		*sourceText = textprocessing.Text(input)
	}
	suffixArray := textprocessing.NewSuffixArray(sourceText)
	lcp := suffixArray.Lcp()
	fmt.Println(suffixArray)
	fmt.Println(lcp)
	fmt.Println(sentenceIndices)
}