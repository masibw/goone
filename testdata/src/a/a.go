package main

import (
	"database/sql"
	"fmt"
	"log"
)


type Person struct {
	Name string
	JobID int
}

type Job struct {
	JobID int
	Name string
}

func main(){

	cnn, _ := sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Query("SELECT name, job_id FROM person")

	defer rows.Close()

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Name,&person.JobID); err != nil {
			log.Fatal(err)
		}

		var job Job
		// This is N+1 Query
		if err := cnn.QueryRow("SELECT job_id, name FROM Job WHERE job_id = ?",person.JobID).Scan(&job.JobID,&job.Name); err != nil { //want "N+1 query"
			log.Fatal(err)
		}
		fmt.Println(person.Name,job.Name)
	}

}
