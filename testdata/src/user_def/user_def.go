package user_def

import (
	"dummy_type"
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
	dummy := dummy_type.Dummy{}
	for  {
		dummy.DummyFunc()//want "this query is called in a loop"
	}
}
