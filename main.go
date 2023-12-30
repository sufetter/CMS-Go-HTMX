package main

import (
	"encoding/json"
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

type Block struct {
	Title string
	Body  string
}

type Triad struct {
	Subject           string `json:"subject"`
	Predicate         string `json:"predicate"`
	AdditionalMembers string `json:"additionalMembers"`
}

func main() {
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./static/images"))))
	http.Handle("/sounds/", http.StripPrefix("/sounds", http.FileServer(http.Dir("./static/sounds"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/dist/", http.StripPrefix("/dist", http.FileServer(http.Dir("./dist"))))
	http.Handle("/Build/", http.StripPrefix("/Build", http.FileServer(http.Dir("./static/Unity/Build"))))
	http.Handle("/TemplateData/", http.StripPrefix("/TemplateData", http.FileServer(http.Dir("./static/Unity/TemplateData"))))
	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir("./static/data"))))

	http.HandleFunc("/", getMain)
	http.HandleFunc("/struct", getStruct)
	http.HandleFunc("/preview", getPreview)
	http.HandleFunc("/stand", getStand)
	http.HandleFunc("/knowledge", getKnowledge)

	port := getPort()
	log.Printf("Server started: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	port := defaultPort
	if len(os.Args[1:]) > 0 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil && p > 1000 && p < 80 {
			port = os.Args[1]
		} else {
			log.Printf("Wrong port, use default: %s\n", port)
		}
	}
	return port
}

func getMain(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/index.html", "./static/header.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPreview(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/preview.html", "./static/infoBlock.html", "./static/header.html"))
	jsonFile, err := os.Open("./static/data/info.json")
	if err != nil {
		http.Error(w, "Failed to open JSON file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			http.Error(w, "Failed to close JSON file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		http.Error(w, "Failed to read JSON file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var blocks []Block
	err = json.Unmarshal(byteValue, &blocks)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, blocks)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getStruct(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/water.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
func getStand(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/stand.html", "./static/header.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getKnowledge(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/data.html", "./static/dataBlock.html", "./static/header.html"))
	jsonFile, err := os.Open("./static/data/data_normal.json")
	if err != nil {
		log.Fatal("Failed to open JSON file: ", err.Error())
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal("Failed to close JSON file: ", err.Error())
			return
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Failed to read JSON file: ", err.Error())
		return
	}

	var triads []Triad
	err = json.Unmarshal(byteValue, &triads)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON: ", err.Error())
		return
	}
	err = tmpl.Execute(w, triads)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
