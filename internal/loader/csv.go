package loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"topRepos/internal/models"
)

func LoadCSV(filePath string) ([]models.Commit, error) {
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
