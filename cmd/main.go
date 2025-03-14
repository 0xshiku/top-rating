package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"topRepos/internal/loader"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func main() {
	addr := ":4000"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)

	server := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	filePath := "./commits.csv"

	data, err := loader.LoadCSV(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(len(data))

	infoLog.Printf("Starting server on %s", addr)

	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
