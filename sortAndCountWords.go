package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func ConvertJsonToWordList(inputJ []byte) (stringL map[string]string, err error) {
	err = json.Unmarshal(inputJ, &stringL)
	if err != nil {
		fmt.Println("Error in Unmarshal")
		return
	}

	checkErrorStr, ok := stringL["text"]
	if !ok {
		err = errors.New("missing fields: text keyword, string of words")
		return
	} else if checkErrorStr == "" {
		err = errors.New("missing fields: string of words")
		return
	}
	return
}

func convertStringListToWordMap(stringL map[string]string) (wordM map[string]int) {
	//Convert the string to lower case before creating a word map out of it to avoid multiple entries for a word
	stringL["text"] = strings.ToLower(stringL["text"])

	// Convert string into list of words
	wordList := strings.Split(stringL["text"], " ")
	fmt.Println("Lowercase tokenized Word List:", wordList)

	// Map to store words as key and frequency as value
	wordM = make(map[string]int)

	//Iterate the word list to create a map of those words with corresponding frequencies
	for i := range wordList {
		wordKey := wordList[i]
		_, ok := wordM[wordKey]
		if !ok {
			fmt.Println("Creating a new Entry in Word map for", wordKey)
			wordM[wordKey] = 1
		} else {
			wordM[wordKey] = wordM[wordKey] + 1
			fmt.Printf("Found %s in the word map, updated the frequency to %d\n", wordKey, wordM[wordKey])
		}
	}
	return
}

// Structure representing mirror of each element of Json output
type jsonKeys struct {
	W string `json:"w"`
	C int    `json:"c"`
}

func convertWordmapToJson(wordM map[string]int) (outputJson []byte, err error) {
	//Create a slice of JSON data elements of type jsonKeys
	jsonData := make([]jsonKeys, len(wordM))

	// Create a temp slice of strings to store the words (fetched from word map key values) in lexicographic order
	tempWordsSlice := make([]string, 0, len(wordM))
	for word := range wordM {
		tempWordsSlice = append(tempWordsSlice, word)
	}

	sort.Strings(tempWordsSlice)
	fmt.Println("Sorted word slice:", tempWordsSlice)

	/* 1. Iterate sorted words
	   2. Fetch the frequency of that word from map using word as key
	   3. Create Json-data struture elments using word and frequency for w and c */
	for i, wordKeyForMap := range tempWordsSlice {
		jsonData[i] = jsonKeys{
			W: wordKeyForMap,
			C: wordM[wordKeyForMap],
		}
	}

	// Convert the structure data to JSON format
	outputJson, err = json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error in Marshalling the data")
	}
	return
}

func GetJsonOutputFromJsonInput(inputJson []byte) (outputJson []byte, err error) {
	stringList, err := ConvertJsonToWordList(inputJson)
	if err != nil {
		return
	}
	fmt.Println("String list:", string(stringList["text"]))
	wordMap := convertStringListToWordMap(stringList)
	fmt.Println("Word Map:", wordMap)
	outputJson, err = convertWordmapToJson(wordMap)
	//fmt.Println("Json data", string(outputJson))
	return
}

func main() {

	// Valid Json inputs
	inputJson := []byte(`{"text": "cat Mat bat Rat Cat cat Bat"}`)
	//inputJson := []byte(`{"text": "cat Mat bat Rat hat sat nat"}`)

	// Various Json strings formats to test negative scenarios
	//inputJson := []byte(`{"text": "cat" "Mat bat Rat cat Cat bat"}`)
	//inputJson := []byte(``)
	//inputJson := []byte(`{"text":}`)
	//inputJson := []byte(`{"text":""}`)
	//inputJson := []byte(`{"":""}`)

	outputJson, err := GetJsonOutputFromJsonInput(inputJson)
	if err == nil {
		fmt.Println("SUCCESS! output = ", string(outputJson))
	} else {
		fmt.Println("FAILED with err:", err)
	}
}
