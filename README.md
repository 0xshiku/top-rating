# Top Scoring Repositories

This project aims to demonstrate top scoring repositories from a csv included in the project.

The have an idea of the csv representation here's an example:

| Field        | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `timestamp`  | The unix timestamp of the commit                             |
| `user`       | The GitHub username of the commit author. This data is unreliable and blank indicates the author's identity is unknown. |
| `repository` | The repository name which the commit was pushed to.          |
| `files`      | The number of files changes by the commit                    |
| `additions`  | The number of line additions in this commit                  |
| `deletions`  | The number of deletions in this commit                       |


It is assumed that the definition of activity score is the sum of additions and deletions, since these dictate the number of lines being changed over time.

## Setting up the project

I decided to create a small web boilerplate for this project. This will give an idea of how I like to organize my code and will give some structure for us to reason about later.

To run the project, make sure you have docker installed and simply run:

```
go run ./cmd/main.go
```

There are no external dependencies, just the go standard library.

## Algorithm

To solve this challenge, the first step was to plan and deconstruct the problem in hand.
The csv contains a lot of data and it can repeat since this is a dataset that represent commits over time. So first things first, a **total scoring system** needed to be developed. This system would accumulate scores (additions plus deletions) per repo. The idea here is to each repo be unique in our dataset.

Now that we have quality data in our hands, the ranking system needs to be developed.

So effectively what we need is to divide this problem in two smaller problems and tackle both.

### Get Activity Score Per Repo

```
func getActivityScorePerRepo(data []models.Commit) (map[string]int, error) {
    // 1
	dataArray := make(map[string]int)
    // 2
	for _, item := range data {
        // 3
		additions, err := strconv.Atoi(item.Additions)
		if err != nil {
			return nil, fmt.Errorf("Invalid additions value: %w", err)
		}
        // 3
		deletions, err := strconv.Atoi(item.Deletions)
		if err != nil {
			return nil, fmt.Errorf("Invalid deletions value: %w", err)
		}
        // 4
		total := additions + deletions
		dataArray[item.Repository] += total
	}
    //5
	return dataArray, nil
}
```
Let's explain this first part of the algorithm by steps:

1 - Creates a map to store our end result

2- Iterate over the commits received (data)

3 - Transforms string values that we need as integers

4 - Adds the additions and deletions in terms of line of code changed to the current iteration. After that we add up to the same key, to guarantee uniqueness on our dataset

5 - returns the results

This algorithm represents a time complexity of linear time (O(n)).

### Get Top N Repositories

So far so good, But we also need to get top ranking system as part of our algorithm.
Since this logic is somehow separated I decided to create a separated function for it.

```
func getTopNRepositories(repoScores map[string]int, topN int) []models.RepoScore {
	// 1
	type entry struct {
		repo  string
		score int
	}

	// 2
	entries := make([]entry, 0, len(repoScores))
	for repo, score := range repoScores {
		entries = append(entries, entry{repo, score})
	}

	// 3
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	// 4
	topResults := make([]models.RepoScore, 0, topN)
	for i := 0; i < topN && i < len(entries); i++ {
		topResults = append(topResults, models.RepoScore{
			Name:  entries[i].repo,
			Score: entries[i].score,
		})
	}
    // 5
	return topResults
}
```

Here are the steps:

1 - Creates a temporary struct to make sorting easier

2 - Creates a slice of entry structs. When created it has a size of 0. But capacity is the length of repoScores. We then iterate over the slice to append the different entries

3 - Here we use sort.Slice to sort everything in the order that is desirable

4 - Creates a slice with 0 size and capacity of the maximum we want, and passed as a parameter. In this case is 10. Then iterates those same number of elements to append it to topResults

5 - Return the top results

Time complexity for this algorithm is O(n).
