package main

import "fmt"
import "textprocessing"
import "io/ioutil"
import "strings"

//get sentenceIndices and Text from corpus
func preprocess(fileName string) (sentenceIndices []int, sourceText *textprocessing.Text) {
	sourceText = &textprocessing.Text{}
	rawText, _ := ioutil.ReadFile(fileName)
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
	return
}

func main() {
	
	sentenceIndicesSource, sourceText := preprocess("/home/lea/corpora/german_test")

	//get suffixArrays(source + target), LCP(source)
	suffixArraySource := textprocessing.NewSuffixArray(sourceText)
	sentenceIndicesTarget, targetText := preprocess("/home/lea/corpora/english_test")
	lcp := suffixArraySource.Lcp()
	starts, occs := suffixArraySource.GetSubstrings(lcp, 5, 3)
	textIndices := suffixArraySource.GetSlicesTokenIndices(starts, occs)


//	get alignments:
//	1) get sentence indices for found substrings
	substringPositions := make([][]int, len(textIndices))
	for idx, tokenIndices := range textIndices {
		substringPositions[idx] = make([]int, len(tokenIndices))
		for index, _ := range tokenIndices {
			substringPositions[idx][index] = sentenceIndicesSource[tokenIndices[index]]
		}
	}
	
/*	fmt.Println(suffixArraySource)
	fmt.Println(substringPositions)
	fmt.Println(textIndices)*/
//	fmt.Println(substringPositions)

	//look at an example for every matching substring + #occurences
 /*	for idx, _ := range starts {
		textIndex := suffixArraySource.GetTokenIndex(starts[idx])
		fmt.Println((*sourceText)[textIndex:textIndex+3], occs[idx])
	}*/
	
	//recover sentences
	//output: {substr, sentence_german, sentence_english}
	for idx, tokenIndices := range textIndices {
		for index, start := range tokenIndices {
			var sourceSentence []textprocessing.Word
			var targetSentence []textprocessing.Word
			for id, el := range sentenceIndicesSource {
				if el == substringPositions[idx][index] {
					sourceSentence = append(sourceSentence, (*sourceText)[id])
				}
			}
			for id, el := range sentenceIndicesTarget {
				if el == substringPositions[idx][index] {
					targetSentence = append(targetSentence, (*targetText)[id])
				}
			}
			fmt.Println("substring", (*sourceText)[start:start+5], occs[idx])
			fmt.Println("german", sourceSentence)
			fmt.Println("english", targetSentence,"\n\n")
		}

	}

}