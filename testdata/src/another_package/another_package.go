package another_package

import (
	"database/sql"
	"fmt"
	"github.com/masibw/go_one/testdata/src/separated"
	"log"
)

var cnn *sql.DB

func main() {

	cnn, _ = sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person separated.Person
		if err := rows.Scan(&person.Name, &person.JobID); err != nil {
			log.Fatal(err)
		}
		job := separated.GetJob2(person) //want "this query might be causes bad performance"
		fmt.Println(person.Name, job.Name)
	}

}
