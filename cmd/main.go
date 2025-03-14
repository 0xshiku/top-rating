package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"topRepos/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}
type templateData struct {
	RepoScores []models.RepoScore
}

func loadCSV(filePath string) ([]models.Commit, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	header := records[0]
	var dataArray []models.Commit

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) != len(header) {
			return nil, fmt.Errorf("row %d has different number of columns than header", i+1)
		}

		username := "unknown"
		if len(row[1]) > 0 {
			username = row[1]
		}

		commit := models.Commit{
			Timestamp:  row[0],
			Username:   username,
			Repository: row[2],
			Files:      row[3],
			Additions:  row[4],
			Deletions:  row[5],
		}

		dataArray = append(dataArray, commit)
	}

	return dataArray, nil
}

func getActivityScorePerRepo(data []models.Commit) (map[string]int, error) {
	dataArray := make(map[string]int)

	for _, item := range data {

		additions, err := strconv.Atoi(item.Additions)
		if err != nil {
			return nil, fmt.Errorf("Invalid additions value: %w", err)
		}

		deletions, err := strconv.Atoi(item.Deletions)
		if err != nil {
			return nil, fmt.Errorf("Invalid deletions value: %w", err)
		}

		total := additions + deletions
		dataArray[item.Repository] += total
	}

	return dataArray, nil
}

func getTopNRepositories(repoScores map[string]int, topN int) []models.RepoScore {
	// temporary struct to make sorting easier
	type entry struct {
		repo  string
		score int
	}

	// Creates a slice of entry structs
	// When created it has a size of 0
	// But capacity is the length of repoScores
	entries := make([]entry, 0, len(repoScores))
	// Iterates throught repoScores map and map it to entries
	for repo, score := range repoScores {
		entries = append(entries, entry{repo, score})
	}

	// Here we use sort.Slice to sort everything
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	// creates a slice with 0 size and capacity of the maximum we want
	topResults := make([]models.RepoScore, 0, topN)
	// Then iterate trought that maximum and append to TopResults
	for i := 0; i < topN && i < len(entries); i++ {
		topResults = append(topResults, models.RepoScore{
			Name:  entries[i].repo,
			Score: entries[i].score,
		})
	}

	return topResults
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

	filePath := "./commits.csv"

	data, err := loadCSV(filePath)
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

	infoLog.Printf("Starting server on %s", addr)

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
