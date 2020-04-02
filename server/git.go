package main

import (
	"log"
	"encoding/csv"
	"os"
)

type GitShas struct {
	WebShellSha string
	ScriptSha string
}

func loadGitShas(path string) *GitShas {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open csv")
	}

	r := csv.NewReader(file)
	record, err := r.Read()
	if err != nil {
		log.Fatal("CSV Parse Error")
	}

	return &GitShas{
		WebShellSha: record[0],
		ScriptSha: record[1],
	}
}