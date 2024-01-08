package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultPort = "3000"
)

// JSON structs for templates, TO DO: convert it to a two-dimensional array

type InfoBlock struct {
	Title string
	Body  string
}

type BaseResponse struct {
	Answer    string
	Question  string
	IsImage   bool
	ImagePath string
}

var blacklist = []string{"замена", "замены", "атрибут", "маршрут", "член", "нет"}

var data *Data

func main() {
	//data.ParseAnswer("уга беккурель производит что?")

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
	var blocks []InfoBlock
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
	question := r.FormValue("question")
	answer, err := data.ParseAnswer(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imagePath := ""
	isImage := false
	if strings.Contains(answer, "src") {
		re := regexp.MustCompile(`src=\{([^}]*)}`)
		match := re.FindStringSubmatch(answer)
		if len(match) > 1 {
			imagePath = match[1]
			isImage = true
			answer = strings.Replace(answer, match[0], "", -1)
		}
	}

	switch answer {
	case "":
		tmpl, _ := template.ParseFiles("./static/data_block/notFoundBlock.html")
		err = tmpl.ExecuteTemplate(w, "notFoundBlock", "Ответ не найден")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	case "Сказуемое не найдено":
		tmpl, _ := template.ParseFiles("./static/data_block/notFoundBlock.html")
		err = tmpl.ExecuteTemplate(w, "notFoundBlock", "Сказуемое не найдено")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	default:
		tmpl, _ := template.ParseFiles("./static/data_block/answerBlock.html")
		data := BaseResponse{
			Answer:    answer,
			Question:  question,
			IsImage:   isImage,
			ImagePath: imagePath,
		}
		err = tmpl.ExecuteTemplate(w, "answerBlock", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}

func executeTemplate(w http.ResponseWriter, files []string, data interface{}) {
	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
