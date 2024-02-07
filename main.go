package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultPort = "3000"
)

// JSON structs for responses

type BaseResponse struct {
	Answer    string
	Question  string
	IsImage   bool
	ImagePath string
}

type NotFoundResponse struct {
	Question string
	Error    string
}

var blacklist = []string{"замена", "замены", "атрибут", "маршрут", "член", "нет"}

var data *Data

var templates *template.Template

func main() {
	handleStaticFiles()
	http.HandleFunc("/", getMain)
	http.HandleFunc("/struct", getStruct)
	http.HandleFunc("/preview", getPreview)
	http.HandleFunc("/stand", getStand)
	http.HandleFunc("/knowledge", getKnowledge)
	http.HandleFunc("/api/knowledge", getBaseData)

	var err error
	templates, err = templates.ParseGlob("./static/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
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

	//Also, Unity and Adobe don't like HTMX, and I was forced to stop using it on the relevant pages

	for path, dir := range staticFiles {
		http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
	}
}

func getPort() string {
	portPtr := flag.String("port", defaultPort, "port to listen on")
	flag.Parse()
	if p, err := strconv.Atoi(*portPtr); err == nil && p > 1024 && p < 49151 {
		return *portPtr
	}
	log.Printf("Invalid port provided, using default: %s\n", defaultPort)
	return defaultPort
}

func getMain(w http.ResponseWriter, _ *http.Request) {
	
        //Attempts to more "deeply" track the work of the handler
	//pc, _, _, _ := runtime.Caller(1)
	//funcName := runtime.FuncForPC(pc).Name()
	//log.Printf("Request received from %s", r.RemoteAddr)
	//log.Printf("/ path is called by %s", funcName)
	//for name, headers := range r.Header {
	//	for _, h := range headers {
	//		fmt.Printf("%v: %v\n", name, h)
	//	}
	//}
	
	err := executeTemplate(w, []string{"index", "headerBlock"}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPreview(w http.ResponseWriter, _ *http.Request) {
	var blocks [][]string
	err := readJSON("./static/data_json/info.json", &blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = executeTemplate(w, []string{"preview", "headerBlock"}, blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getStruct(w http.ResponseWriter, _ *http.Request) {
	err := executeTemplate(w, []string{"water"}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getStand(w http.ResponseWriter, _ *http.Request) {
	err := executeTemplate(w, []string{"stand"}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getKnowledge(w http.ResponseWriter, _ *http.Request) {
	var triads [][]string
	err := readJSON("./static/data_json/data.json", &triads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = executeTemplate(w, []string{"data", "headerBlock"}, triads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		data := NotFoundResponse{
			Question: question,
			Error:    "Ответ не найден",
		}
		err := executeTemplate(w, []string{"notFoundBlock"}, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	case "Сказуемое не найдено":
		data := NotFoundResponse{
			Question: question,
			Error:    "Сказуемое не найдено",
		}
		err = executeTemplate(w, []string{"notFoundBlock"}, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	default:
		data := BaseResponse{
			Answer:    answer,
			Question:  question,
			IsImage:   isImage,
			ImagePath: imagePath,
		}
		err = executeTemplate(w, []string{"answerBlock"}, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}

func executeTemplate(w http.ResponseWriter, tmplNames []string, data interface{}) error {
	for _, name := range tmplNames {
		if !strings.Contains(name, "Block") {
			name += ".html"
		}
		tmpl := templates.Lookup(name)
		if tmpl == nil {
			return fmt.Errorf("template %s not found", name)
		}
		err := tmpl.ExecuteTemplate(w, name, data)
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}
	}
	return nil
}
