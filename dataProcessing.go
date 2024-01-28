package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type Data struct {
	Endings [][]string
	Data    [][]string
}

func init() {
	var err error
	data, err = NewData()
	if err != nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}
}

// TO DO: make sure to use goroutines correctly

func NewData() (*Data, error) {
	d := &Data{}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := readJSON("./static/data_json/endings.json", &d.Endings)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := readJSON("./static/data_json/data.json", &d.Data)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
	}()

	wg.Wait()

	return d, nil
}

func readJSON(filename string, data interface{}) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer func() {
		if closeErr := jsonFile.Close(); closeErr != nil {
			log.Printf("Failed to close JSON file: %v", closeErr)
		}
	}()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	err = json.Unmarshal(byteValue, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func cleanInput(input string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Zа-яА-Я0-9]+")
	if err != nil {
		return "", err
	}

	processedString := reg.ReplaceAllString(input, " ")
	processedString = strings.ToLower(processedString)

	return processedString, nil
}

func formatAnswer(entry []string) string {
	firstLetter := []rune(entry[0])
	firstLetter[0] = unicode.ToUpper(firstLetter[0])
	return string(firstLetter) + " " + entry[1] + " " + entry[2]
}
