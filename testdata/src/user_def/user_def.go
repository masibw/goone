package user_def

import (
	"context"
)

type Person struct {
	Name  string
	JobID int
}

type Job struct {
	JobID int
	Name  string
}

func ForStmt() {
	ctx, cancel := context.WithTimeout(context.Background(), 9)
	defer cancel()
	for {
		ctx.Done() //want "this query is called in a loop"
	}
}
