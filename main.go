package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort = "3000"
)

// JSON structs for templates
type infoBlock struct {
	Title string
	Body  string
}

type pseudoEndingBlock struct {
	Ending string
	Regex  string
}

var blacklist []string = []string{"замена", "замены", "атрибут", "маршрут", "член", "нет"}

func main() {
	//parseAnswer("уга бугагде беккерель производиться ыаы")
	handleStaticFiles()
	http.HandleFunc("/", getMain)
	http.HandleFunc("/struct", getStruct)
	http.HandleFunc("/preview", getPreview)
	http.HandleFunc("/stand", getStand)
	http.HandleFunc("/knowledge", getKnowledge)
	http.HandleFunc("/api/knowledge", getBaseData)
	port := getPort()
	log.Printf("Server started: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleStaticFiles() {
	staticFiles := map[string]string{
		"/static/":       "./static/",
		"/dist/":         "./dist",
		"/images/":       "./static/images",
		"/js/":           "./static/js",
		"/sounds/":       "./static/sounds",
		"/unity/":        "./static/unity",
		"/Build/":        "./static/unity/Build",
		"/TemplateData/": "./static/unity/TemplateData",
	}

	//Paths are written explicitly due to problems
	//with implicit access to them (Unity and Adobe libraries working only with explicit paths)

	for path, dir := range staticFiles {
		http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
	}
}

func getPort() string {
	port := defaultPort
	if len(os.Args[1:]) > 0 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil && p > 1024 && p < 49151 {
			port = os.Args[1]
		} else {
			log.Printf("Wrong port, use default: %s\n", port)
		}
	}
	return port
}

func getMain(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, []string{"./static/index.html", "./static/data_block/header.html"}, nil)
}

func getPreview(w http.ResponseWriter, r *http.Request) {
	var blocks []infoBlock
	err := readJSON("./static/data_json/info.json", &blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	executeTemplate(w, []string{"./static/preview.html", "./static/data_block/infoBlock.html", "./static/data_block/header.html"}, blocks)
}

func getStruct(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, []string{"./static/water.html"}, nil)
}

func getStand(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, []string{"./static/stand.html", "./static/data_block/header.html"}, nil)
}

func getKnowledge(w http.ResponseWriter, r *http.Request) {
	var triads [][]string
	err := readJSON("./static/data_json/data.json", &triads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	executeTemplate(w, []string{"./static/data.html", "./static/data_block/dataBlock.html", "./static/data_block/header.html"}, triads)
}

func getBaseData(w http.ResponseWriter, r *http.Request) {
	answer, err := parseAnswer(r.FormValue("answer"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if answer == "" {
		executeTemplate(w, []string{"./static/data_block/notFoundBlock.html"}, nil)
		return

	}
	executeTemplate(w, []string{""}, answer)
}

func executeTemplate(w http.ResponseWriter, files []string, data interface{}) {
	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
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
