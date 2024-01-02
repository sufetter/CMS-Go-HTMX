package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

func parseAnswer(question string) (string, error) {
	var endings [][]string
	err := readJSON("./static/data_json/endings.json", &endings)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file: %w", err)
	}
	var data [][]string
	err = readJSON("./static/data_json/data.json", &data)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file: %w", err)
	}
	cleanQuestion := cleanInput(question)
	answer := getAnswer(cleanQuestion, endings, blacklist, data)
	log.Printf(answer)
	return answer, nil
}

func getAnswerBySubject(subjectWords []string, data [][]string) string {
	for _, entry := range data {
		for _, subject := range subjectWords {
			for _, wordInEntry := range strings.Fields(entry[0]) {
				if levenshtein(wordInEntry, subject) <= 2 {
					firstLetter := []rune(entry[0])
					firstLetter[0] = unicode.ToUpper(firstLetter[0])
					return string(firstLetter) + " " + entry[1] + " " + entry[2]
				}
			}
		}
	}
	return ""
}

func getAnswer(question string, pseudoEndings [][]string, blacklist []string, data [][]string) string {
	predicate, index := findPredicate(question, pseudoEndings, blacklist)
	if index == -1 {
		return "Сказуемое не найдено"
	}
	words := strings.Fields(question)
	subjectWords := strings.Fields(strings.Join(words[:index], " "))
	additionalWords := strings.Fields(strings.Join(words[index+1:], " "))

	for _, entry := range data {
		if levenshtein(entry[1], predicate) <= 2 {
			subjectMatch := false
			for _, subject := range subjectWords {
				for _, wordInEntry := range strings.Fields(entry[0]) {
					if levenshtein(wordInEntry, subject) <= 5 {
						subjectMatch = true
						break
					}
				}
				if subjectMatch {
					break
				}
			}
			additionalMatch := false
			for _, additional := range additionalWords {
				for _, wordInEntry := range strings.Fields(entry[2]) {
					if levenshtein(wordInEntry, additional) <= 5 {
						additionalMatch = true
						break
					}
				}
				if additionalMatch {
					break
				}
			}
			if subjectMatch || additionalMatch {
				firstLetter := []rune(entry[0])
				firstLetter[0] = unicode.ToUpper(firstLetter[0])
				return string(firstLetter) + " " + entry[1] + " " + entry[2]
			}
		}
	}
	return getAnswerBySubject(subjectWords, data)
}

func findPredicate(sentence string, pseudoEndings [][]string, blacklist []string) (string, int) {
	words := strings.Fields(sentence)
	for i, word := range words {
		if contains(blacklist, word) {
			continue
		}
		for _, pair := range pseudoEndings {
			re := regexp.MustCompile(pair[1])
			if re.MatchString(word) {
				return word, i
			}
		}
	}
	return "", -1
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func cleanInput(input string) string {
	reg, err := regexp.Compile("[^a-zA-Zа-яА-Я0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	processedString := reg.ReplaceAllString(input, " ")
	processedString = strings.ToLower(processedString)

	return processedString
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	for i := 1; i <= len(a); i++ {
		matrix[i][0] = i
	}

	for j := 1; j <= len(b); j++ {
		matrix[0][j] = j
	}

	for j := 1; j <= len(b); j++ {
		for i := 1; i <= len(a); i++ {
			if a[i-1] == b[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(
					matrix[i-1][j]+1,
					matrix[i][j-1]+1,
					matrix[i-1][j-1]+1,
				)
			}
		}
	}

	return matrix[len(a)][len(b)]
}

//func checkEndings() error {
//	var triads [][]string
//	err := readJSON("./static/data_json/data.json", &triads)
//	if err != nil {
//		return fmt.Errorf("failed to read triads JSON file: %w", err)
//	}
//
//	var endings [][]string
//	err = readJSON("./static/data_json/endings.json", &endings)
//	if err != nil {
//		return fmt.Errorf("failed to read endings JSON file: %w", err)
//	}
//
//	for _, triad := range triads {
//		predicate := triad[1]
//		pseudoEnding := predicate[len(predicate)-2:]
//
//		found := false
//		for _, ending := range endings {
//			if strings.Contains(ending[1], pseudoEnding) {
//				found = true
//				break
//			}
//		}
//
//		if !found {
//			fmt.Printf("Псевдоокончание '%s' предиката '%s' не найдено в массиве псевдоокончаний. Добавляем его.\n", pseudoEnding, predicate)
//			endings = append(endings, []string{pseudoEnding, pseudoEnding})
//		}
//	}
//
//	return nil
//}
