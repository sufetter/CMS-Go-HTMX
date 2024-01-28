package main

import (
	"regexp"
	"strings"
)

func (d *Data) ParseAnswer(question string) (string, error) {
	cleanQuestion, err := cleanInput(question)
	if err != nil {
		return "", err
	}
	answer := d.GetAnswer(cleanQuestion, blacklist)
	//log.Print(answer)
	return answer, nil
}

//The attempt to implement the use of goroutines here completely failed;
//too many additional operations need to be implemented for them to work correctly.
//It was decided to optimize the algorithmic component.
//Due to the lack of enough time, I’m not sure that the problem was solved optimally, BUT IT WORKS!

func (d *Data) GetAnswer(question string, blacklist []string) string {
	predicate, index := findPredicate(question, d.Endings, blacklist)
	if index == -1 {
		return "Сказуемое не найдено"
	}
	words := strings.Fields(question)
	var subjectWords []string
	var additionalWords []string

	if len(words) > index+1 {
		subjectWords = strings.Fields(strings.Join(words[index+1:index+2], " "))
		additionalWords = strings.Fields(strings.Join(words[index+2:], " "))
	}

	for _, entry := range d.Data {
		if levenshtein(entry[1], predicate) <= 2 {
			subjectMatch := false
			for _, subject := range subjectWords {
				for _, wordInEntry := range append(strings.Fields(entry[0]), strings.Fields(entry[2])...) {
					if levenshtein(wordInEntry, subject) <= 5 {
						subjectMatch = true
						break
					}
				}
				if subjectMatch {
					return formatAnswer(entry)
				}
			}
		}
	}

	for _, entry := range d.Data {
		if levenshtein(entry[1], predicate) <= 2 {
			additionalMatch := false
			for _, additional := range additionalWords {
				for _, wordInEntry := range strings.Fields(entry[2]) {
					if levenshtein(wordInEntry, additional) <= 5 {
						additionalMatch = true
						break
					}
				}
				if additionalMatch {
					return formatAnswer(entry)
				}
			}
		}
	}

	for _, entry := range d.Data {
		if levenshtein(entry[1], predicate) <= 2 {
			return formatAnswer(entry)
		}
	}

	return d.GetAnswerBySubject(subjectWords)
}

func (d *Data) GetAnswerBySubject(subjectWords []string) string {
	for _, entry := range d.Data {
		for _, subject := range subjectWords {
			for _, column := range []int{0, 2} {
				for _, wordInEntry := range strings.Fields(entry[column]) {
					if levenshtein(wordInEntry, subject) <= 2 {
						return formatAnswer(entry)
					}
				}
			}
		}
	}
	return ""
}

func findPredicate(sentence string, pseudoEndings [][]string, blacklist []string) (string, int) {
	words := strings.Fields(sentence)
	for i, word := range words {
		if contains(blacklist, word) {
			continue
		}
		for _, pair := range pseudoEndings {
			re := regexp.MustCompile(pair[1] + "$")
			if re.MatchString(word) {
				//if pair[0] == pair[1] {
				//	log.Printf("%v", words[i+1])
				//
				//	if i+1 < len(words) {
				//		return words[i+1], i + 1
				//	}
				//} else {
				//	return word, i
				//}
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

func levenshtein(a, b string) int {
	ar := []rune(a)
	br := []rune(b)

	if len(ar) == 0 {
		return len(br)
	}
	if len(br) == 0 {
		return len(ar)
	}

	matrix := make([][]int, len(ar)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(br)+1)
	}

	for i := 1; i <= len(ar); i++ {
		matrix[i][0] = i
	}

	for j := 1; j <= len(br); j++ {
		matrix[0][j] = j
	}

	for j := 1; j <= len(br); j++ {
		for i := 1; i <= len(ar); i++ {
			if ar[i-1] == br[j-1] {
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

	return matrix[len(ar)][len(br)]
}

//Turned out that according to the technical specifications, this is not necessary...
//I’ll leave it here just in case

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
