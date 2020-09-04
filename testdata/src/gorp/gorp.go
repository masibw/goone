package main

import (
	"database/sql"
	"fmt"
	"log"

	"gopkg.in/gorp.v1"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Name  string`db:"name"`
	JobID int`db:"job_id"`
}

type Job struct {
	JobID int`db:"job_id"`
	Name  string`db:"name"`
}

func main() {
	db, err := sql.Open( "mysql","user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err.Error())
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	defer dbmap.Db.Close()

	// Read
	persons := []Person{}
	_, err = dbmap.Select(&persons,"SELECT * FROM persons")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, person := range persons {
		var job Job
		err = dbmap.SelectOne(&job,"SELECT * FROM jobs where job_id=?",person.JobID)//want "this query called in loop"
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(person.Name, job.Name)
	}

}

