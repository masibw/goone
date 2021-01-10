package separated

import (
	"database/sql"
	"fmt"
	"log"
)

type Person struct {
	Name  string
	JobID int
}

type Job struct {
	JobID int
	Name  string
}

func getJob(person Person) Job {
	var job Job
	// This is N+1 Query
	err := cnn.QueryRow("SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID).Scan(&job.JobID, &job.Name)
	if err != nil {
		log.Fatal(err)
	}
	return job
}

var cnn *sql.DB

func main() {

	cnn, _ = sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Name, &person.JobID); err != nil {
			log.Fatal(err)
		}
		job := getJob(person) //want "this query is called in a loop"
		job = getJob(person)  //want "this query is called in a loop"
		job = GetJob2(person) //want "this query is called in a loop"
		job = NotCallQuery()  // This should be not detected
		fmt.Println(person.Name, job.Name)
	}

}
