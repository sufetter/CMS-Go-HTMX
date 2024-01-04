package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Data struct {
	Endings [][]string
	Data    [][]string
}

type Result struct {
	Data interface{}
	Err  error
}

func init() {
	var err error
	data, err = NewData()
	if err != nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}
}

func NewData() (*Data, error) {
	d := &Data{}
	resChan := make(chan Result, 2)

	go func() {
		err := readJSON("./static/data_json/endings.json", &d.Endings)
		resChan <- Result{Data: d.Endings, Err: err}
	}()

	go func() {
		err := readJSON("./static/data_json/data.json", &d.Data)
		resChan <- Result{Data: d.Data, Err: err}
	}()

	for i := 0; i < 2; i++ {
		res := <-resChan
		if res.Err != nil {
			return nil, fmt.Errorf("failed to read JSON file: %w", res.Err)
		}
	}

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
