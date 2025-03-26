package loaders

import (
	"encoding/csv"
	"fmt"
	"github.com/0xshiku/toprating/internal/models"
	"os"
)

func LoadCSV(filePath string) ([]models.Commit, error) {
	// First we try to open the file with the os library
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer file.Close()

	// Then we use the csv module to read the file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// We start to assimilate values like the header
	header := records[0]

	// Declare a list of commits as defined in the models
	var dataArray []models.Commit

	// Iterate through each row of records
	// We start at 1 because 0 is the header
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

	// After appending to data array we return the collection
	return dataArray, nil
}
