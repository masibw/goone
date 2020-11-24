package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type Person struct {
	Name  string`db:"name"`
	JobID int`db:"job_id"`
}

type Job struct {
	JobID int`db:"job_id"`
	Name  string`db:"name"`
}

func forStmt() {

	cnn, _ := sqlx.Connect("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Queryx("SELECT name, job_id FROM persons")

	for rows.Next() {
		var person Person
		if err := rows.StructScan(&person); err != nil {
			log.Fatal(err)
		}

		var job Job
		// This is N+1 Query
		err := cnn.Get(&job,"SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID)//want "this query is called in a loop"
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(person.Name, job.Name)
	}

}

func rangeStmt() {
	cnn, _ := sqlx.Connect("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Queryx("SELECT name, job_id FROM persons")

	var persons []Person
	for rows.Next() {
		var person Person
		if err := rows.StructScan(&person); err != nil {
			log.Fatal(err)
		}
		persons = append(persons,person)
	}
	for _,person := range persons{
		var job Job
		// This is N+1 Query
		err := cnn.Get(&job,"SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID)//want "this query is called in a loop"
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(person.Name, job.Name)
	}
}
