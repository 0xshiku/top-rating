package main

import (
	"fmt"
	"github.com/0xshiku/toprating/internal/models"
	"sort"
	"strconv"
)

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
