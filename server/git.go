package main

import (
	"log"
	"encoding/csv"
	"strings"
)

type GitShas struct {
	WebShellSha string
	ScriptSha string
}

func loadGitShas(dir string) *GitShas {
	r := csv.NewReader(strings.NewReader(dir))
	record, err := r.Read()
	if err != nil {
		log.Fatal("CSV Parse Error")
	}

	return &GitShas{
		WebShellSha: record[0],
		ScriptSha: record[1],
	}
}