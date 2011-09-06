package main

import "fmt"
import "textprocessing"
import "io/ioutil"
import "strings"
import "os"

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

func GetSubstringSentenceNumbers(substringIndicesText [][]int) (substringPositions [][]int) {
        substringPositions = make([][]int, len(substringIndicesText))
        for idx, tokenIndices := range substringIndicesText {
                substringPositions[idx] = make([]int, len(tokenIndices))
                for index, _ := range tokenIndices {
                        substringPositions[idx][index] = mapWord2SentenceSource[tokenIndices[index]]
                }
        }
        return
}


func GetSubstringIndicesText(suffixIndices, lengths []int, suffixArray *textprocessing.SuffixArray) (tokenIndices [][]int) {
        // For all substring types: execute the above function
        // same as above for different slice lengths
        tokenIndices = make([][]int, len(suffixIndices))
        for idx, _ := range suffixIndices {
                tokenIndices[idx] = suffixArray.GetTokenIndices(suffixIndices[idx], lengths[idx])
        }
        return
}


var mapWord2SentenceSource []int
var mapSentence2WordSource, mapSentence2WordTarget [][2]int
var sourceText, targetText *textprocessing.Text

func main() {

        //get suffixArrays(source), LCP(source)
        fmt.Println("preprocessing source...")
        mapWord2SentenceSource, mapSentence2WordSource, sourceText = preprocess("./german")
        fmt.Println("done")
        fmt.Println("preprocessing target...")
        _, mapSentence2WordTarget, targetText = preprocess("./english")
        fmt.Println("done")
        fmt.Println("constructing suffix array...")
        suffixArraySource := textprocessing.NewSuffixArray(sourceText)
        fmt.Println("done")
        lcpSource := suffixArraySource.Lcp()
        
        // print systematically
        for length :=12; length >=8; length-- {
                for minOccurrences :=120; minOccurrences >=50; minOccurrences=minOccurrences-10 {
                        fmt.Println(length, minOccurrences)
                        file, _ := os.Create(fmt.Sprintf("L%03dF%03d.dat",length,minOccurrences))
                        firstOccInSuffixArray, occs := suffixArraySource.GetSubstrings(lcpSource, length, minOccurrences)
                        //get starting point for each substring
                        substringIndicesText := GetSubstringIndicesText(firstOccInSuffixArray, occs, suffixArraySource)
                        //get corresp sentence number for each matching substring
                        substringSentenceNumbers := GetSubstringSentenceNumbers(substringIndicesText)
                        for typeIndex, _ := range substringSentenceNumbers {
                                for sentenceNumberIndex, sentenceNumber := range substringSentenceNumbers[typeIndex] {
                                        startPosition := substringIndicesText[typeIndex][sentenceNumberIndex]
                                        fmt.Fprint(file, "string", (*sourceText)[startPosition:startPosition+length], occs[typeIndex], "\n")
                                        fmt.Fprint(file, "german", (*sourceText)[mapSentence2WordSource[sentenceNumber][0]:mapSentence2WordSource[sentenceNumber][1]], "\n")
                                        fmt.Fprint(file, "english", (*targetText)[mapSentence2WordTarget[sentenceNumber][0]:mapSentence2WordTarget[sentenceNumber][1]], "\n\n")
                                }
                        }
                }
        }
}
 
