package main

import (
	"fmt"
	"github.com/0xshiku/toprating/internal/loaders"
	"github.com/0xshiku/toprating/internal/models"
	"html/template"
	"log"
	"net/http"
)

type templateData struct {
	RepoScores []models.RepoScore
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"../ui/html/base.gohtml",
		"../ui/html/partials/nav.gohtml",
		"../ui/html/pages/home.gohtml",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	filePath := "../commits.csv"

	data, err := loaders.LoadCSV(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	activityScore, err := getActivityScorePerRepo(data)
	if err != nil {
		fmt.Println("Error getting acitivity score")
	}

	top := getTopNRepositories(activityScore, 10)
	records := templateData{RepoScores: top}

	err = templates.ExecuteTemplate(w, "base", records)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
