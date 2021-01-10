package another_package

import (
	"database/sql"
	"fmt"

	"github.com/masibw/goone_test/pkg/db"

	"log"
)

var cnn *sql.DB

func main() {

	cnn, _ = sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person db.Person
		if err := rows.Scan(&person.Name, &person.JobID); err != nil {
			log.Fatal(err)
		}
		job := db.GetJobInsideCnn(person) //want "this query is called in a loop"
		job =  db.GetJob2(cnn, person)  //want "this query is called in a loop"
		fmt.Println(person.Name, job.Name)
	}

}
