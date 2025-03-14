package models

type Commit struct {
	Timestamp  string
	Username   string
	Repository string
	Files      string
	Additions  string
	Deletions  string
}
