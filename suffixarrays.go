package main

import "fmt"
import "textprocessing"
import "io/ioutil"
import "strings"
import "os"
import "runtime/pprof"

//get sentenceIndices and Text from corpus
func preprocess(fileName string) (wordIndices []int, sentenceIndices [][2]int, sourceText *textprocessing.Text) {
	sourceText = &textprocessing.Text{}
	rawText, _ := ioutil.ReadFile(fileName)
	sentences := strings.Split(string(rawText), "\n", -1)
	sentenceCount := len(sentences)
	sentenceTexts := make([]textprocessing.Text, sentenceCount)
	wordCount := 0
	for idx, sentence := range sentences {
		sentenceTexts[idx] = *textprocessing.NewText(sentence)
		wordCount += sentenceTexts[idx].Length()
	}
	wordIndices = make([]int, 0, wordCount)
	sentenceIndices = make([][2]int, 0, sentenceCount)
	input := make([]textprocessing.Word, 0, wordCount)
	sentenceIndex := 0
	for idx, sentence := range sentenceTexts {
		for ii := 0; ii < sentence.Length(); ii++ {
			wordIndices = append(wordIndices, idx)
		}
		sentenceIndices = append(sentenceIndices, [2]int{sentenceIndex, sentenceIndex+sentence.Length()})
		sentenceIndex = sentenceIndex + sentence.Length()
		input = append(input, sentence...)
	}
	*sourceText = textprocessing.Text(input)
	return
}

func GetSubstringSentenceNumbers(textIndices [][]int, wordIndices []int) (substringPositions [][]int) {
	substringPositions = make([][]int, len(textIndices))
	for idx, tokenIndices := range textIndices {
		substringPositions[idx] = make([]int, len(tokenIndices))
		for index, _ := range tokenIndices {
			substringPositions[idx][index] = wordIndices[tokenIndices[index]]
		}
	}
	return
}

func main() {
	
	f, _ := os.Create("cpuprofile")

	//get suffixArrays(source + target), LCP(source)
	fmt.Println("preprocessing source...")
	wordIndicesSource, sentenceIndicesSource, sourceText := preprocess("./german")
	fmt.Println("done")
	fmt.Println("preprocessing target...")
	_, sentenceIndicesTarget, targetText := preprocess("./english")
	fmt.Println("done")
	fmt.Println("constructing suffix array...")
	suffixArraySource := textprocessing.NewSuffixArray(sourceText)
	fmt.Println("done")
	lcp := suffixArraySource.Lcp()
	
	pprof.StartCPUProfile(f)
	// print systematically
	for length :=3; length >=3; length-- {
		for minOcc :=3; minOcc >=3; minOcc-- {
			fmt.Println(length, minOcc)
			file, _ := os.Create(fmt.Sprintf("L%03dF%03d.dat",length,minOcc))
			starts, occs := suffixArraySource.GetSubstrings(lcp, length, minOcc)
			//get starting point for each substring
			substringStartIndices := suffixArraySource.GetSlicesTokenIndices(starts, occs)
			//get corresp sentence number for each matching substring
			sentencesOfSubstringOccurrence := GetSubstringSentenceNumbers(substringStartIndices, wordIndicesSource)
			for idx, substringType := range sentencesOfSubstringOccurrence {
				for index, start := range substringType {
					startPosition := substringStartIndices[idx][index]
					fmt.Fprint(file, "string", (*sourceText)[startPosition:startPosition+length], occs[idx], "\n")
					fmt.Fprint(file, "german", (*sourceText)[sentenceIndicesSource[start][0]:sentenceIndicesSource[start][1]], "\n")
					fmt.Fprint(file, "english", (*targetText)[sentenceIndicesTarget[start][0]:sentenceIndicesTarget[start][1]], "\n\n")
				}
			}
		}
	}
	defer pprof.StopCPUProfile()
}
