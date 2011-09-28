package main

import "fmt"
import "textprocessing"
import "io/ioutil"
import "strings"
import "os"
import "sync"
import "runtime"
import "sort"

//get sentenceIndices and Text from corpus
func preprocess(fileName string) (wordIndices []int, sentenceIndices [][2]int, text *textprocessing.Text) {
        text = &textprocessing.Text{}
        rawText, _ := ioutil.ReadFile(fileName)
        sentences := strings.Split(string(rawText), "\n")
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
        *text = textprocessing.Text(input)
        return
}

func GetSubstringSentenceNumbers(substringIndicesText [][]int, mapWord2Sentence []int) (substringPositions [][]int) {
        substringPositions = make([][]int, len(substringIndicesText))
        for idx, tokenIndices := range substringIndicesText {
                substringPositions[idx] = make([]int, len(tokenIndices))
                for index, _ := range tokenIndices {
                        substringPositions[idx][index] = mapWord2Sentence[tokenIndices[index]]
                }
        }
        return
}

func Uniquify(slice []int) ([]int) {
	sort.Ints(slice)
	var writeIdx, lastvalue int
	for idx,value := range slice {
		slice[writeIdx]=value
		if value!=lastvalue || idx==0 {writeIdx++}
			lastvalue=value
	}
	return slice[0:writeIdx]
}


func GetSubstringIndicesText(suffixIndices, lengths []int, suffixArray *textprocessing.SuffixArray) (tokenIndices [][]int) {
        tokenIndices = make([][]int, len(suffixIndices))
        for idx, _ := range suffixIndices {
                tokenIndices[idx] = suffixArray.GetTokenIndices(suffixIndices[idx], lengths[idx])
        }
        return
}

func countCommonItems(set1, set2 []int) (overlap int) {
	set1 = Uniquify(set1)
	set2 = Uniquify(set2)
	overlap = 0
	for _,val1 := range set1 {
		for _, val2 := range set2 {
			if val1==val2 {
				overlap++
			}
		}
	}
	return
}

func computeCorrelations(vecSource, vecTarget [][]int, stringLengthsSource, stringLengthsTarget []int) (source2targetTypes, target2sourceTypes []int, source2targetScores, target2sourceScores []float64, translationFreqSource, translationFreqTarget [][2]int) {
//compute correlations between all sets for German substrings with all sets for English substrings
//return 1 best correlating pair
	source2targetTypes = make([]int, len(vecSource))
	source2targetScores = make([]float64, len(vecSource))
	target2sourceTypes = make([]int, len(vecTarget))
	target2sourceScores = make([]float64, len(vecTarget))
	translationFreqSource = make([][2]int, len(vecSource))
	translationFreqTarget = make([][2]int, len(vecTarget))
	
	for substringTypeSource, set1 := range vecSource {
		for substringTypeTarget, set2 := range vecTarget {
			score := countCommonItems(set1, set2)
			normalizedScoreSource := float64(score)/float64(len(set1))
			normalizedScoreTarget := float64(score)/float64(len(set2))
			if source2targetScores[substringTypeSource] <= normalizedScoreSource && stringLengthsTarget[substringTypeTarget] >= stringLengthsTarget[source2targetTypes[substringTypeSource]]{
				source2targetScores[substringTypeSource] = normalizedScoreSource
				source2targetTypes[substringTypeSource] = substringTypeTarget
				translationFreqSource[substringTypeSource] = [2]int{score, len(set1)}
			}
			if target2sourceScores[substringTypeTarget] <= normalizedScoreTarget && stringLengthsSource[substringTypeSource] >= stringLengthsSource[target2sourceTypes[substringTypeTarget]] {
				target2sourceScores[substringTypeTarget] = normalizedScoreTarget
				target2sourceTypes[substringTypeTarget] = substringTypeSource
				translationFreqTarget[substringTypeTarget] = [2]int{score, len(set2)}
			}
		}
	}
	return
}


var mapWord2SentenceGer, mapWord2SentenceEng []int
//var mapSentence2WordGer, mapSentence2WordEng [][2]int
var ger, eng *textprocessing.Text

var suffixArrayGer, suffixArrayEng *textprocessing.SuffixArray
var lcpGer, lcpEng []int

var substringSentenceNumbersEng, substringSentenceNumbersGer, substringIndicesTextEng, substringIndicesTextGer [][]int
var stringLengthsEng, stringLengthsGer []int

var occsEng, occsGer, firstOccInSuffixArrayEng, firstOccInSuffixArrayGer []int

func main() {
	
        occurrences := 100
	minLength := 10
	maxLength := 15
	
	fmt.Println("**NEW** Parameter: min: ", minLength, "max: ", maxLength, "occurrences: ", occurrences)
	
	
	runtime.GOMAXPROCS(5)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		mapWord2SentenceGer, _, ger = preprocess("../german")
		suffixArrayGer = textprocessing.NewSuffixArray(ger)
		lcpGer = suffixArrayGer.Lcp()
		//get starting point for each substring
		firstOccInSuffixArrayGer, occsGer, stringLengthsGer = suffixArrayGer.GetSubstrings(lcpGer, occurrences, minLength, maxLength)
		//get corresp sentence number for each matching substring
		substringIndicesTextGer = GetSubstringIndicesText(firstOccInSuffixArrayGer, occsGer, suffixArrayGer)
        	//compute correlations
		substringSentenceNumbersGer = GetSubstringSentenceNumbers(substringIndicesTextGer, mapWord2SentenceGer)
		wg.Done()
	}()
	go func() {
		mapWord2SentenceEng, _, eng = preprocess("../english")
        	suffixArrayEng = textprocessing.NewSuffixArray(eng)
        	lcpEng = suffixArrayEng.Lcp()
        	//get starting point for each substring
        	firstOccInSuffixArrayEng, occsEng, stringLengthsEng = suffixArrayEng.GetSubstrings(lcpEng, occurrences, minLength, maxLength)
        	//get corresp sentence number for each matching substring
        	substringIndicesTextEng = GetSubstringIndicesText(firstOccInSuffixArrayEng, occsEng, suffixArrayEng)
        	//compute correlations
        	substringSentenceNumbersEng = GetSubstringSentenceNumbers(substringIndicesTextEng, mapWord2SentenceEng)
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("Finished preprocessing and array construction.")
	
	ger2engTypes, eng2gerTypes, ger2engScores, eng2gerScores, translationFreqGer, translationFreqEng := computeCorrelations(substringSentenceNumbersGer, substringSentenceNumbersEng, stringLengthsGer, stringLengthsEng)
 

	fmt.Println("Printing...")
	ger2eng, _ := os.Create(fmt.Sprintf("ger2eng.dat"))
	eng2ger, _ := os.Create(fmt.Sprintf("eng2ger.dat"))
	for sourceIndex, targetIndex := range ger2engTypes {
			startPositionGer := substringIndicesTextGer[sourceIndex][0] //0 weil brauch ja nur einen
			startPositionEng := substringIndicesTextEng[targetIndex][0]
			fmt.Fprint(ger2eng,ger2engScores[sourceIndex], "@")
			if startPositionGer+stringLengthsGer[sourceIndex] < len(*ger) {
				fmt.Fprint(ger2eng,(*ger)[startPositionGer:startPositionGer+stringLengthsGer[sourceIndex]], startPositionGer+stringLengthsGer[sourceIndex]-startPositionGer, "@")
			}
			if startPositionEng+stringLengthsEng[targetIndex] < len(*eng) {
				fmt.Fprint(ger2eng,(*eng)[startPositionEng:startPositionEng+stringLengthsEng[targetIndex]], startPositionEng+stringLengthsEng[targetIndex]-startPositionEng, "@")
			}
			fmt.Fprint(ger2eng,translationFreqGer[sourceIndex], "\n")
        }
        for sourceIndex, targetIndex := range eng2gerTypes {
			startPositionEng := substringIndicesTextEng[sourceIndex][0]
			startPositionGer := substringIndicesTextGer[targetIndex][0]
			fmt.Fprint(eng2ger,eng2gerScores[sourceIndex], "@")
			if startPositionEng+stringLengthsEng[sourceIndex] < len(*eng) {
				fmt.Fprint(eng2ger,(*eng)[startPositionEng:startPositionEng+stringLengthsEng[sourceIndex]], startPositionEng+stringLengthsEng[sourceIndex]-startPositionEng, "@")
			}
			if startPositionGer+stringLengthsGer[targetIndex] < len(*ger) {
				fmt.Fprint(eng2ger,(*ger)[startPositionGer:startPositionGer+stringLengthsGer[targetIndex]], startPositionGer+stringLengthsGer[targetIndex]-startPositionGer, "@")
			}
			fmt.Fprint(eng2ger,translationFreqEng[sourceIndex], "\n")
        }
}


