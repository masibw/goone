package gorm

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Person struct {
	Name  string
	JobID int
}

type Tabler interface {
	TableName() string
}
// TableName overrides the table name used by User to `profiles`
func (Person) TableName() string {
	return "persons"
}

type Job struct {
	JobID int
	Name  string
}

func main() {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Read
	persons := []Person{}
	if err := db.Find(&persons).Error; err != nil {
		log.Fatal(err.Error())
	}

	for _, person := range persons {
		var job Job
		if err := db.First(&job, person.JobID).Error; err != nil {//want "this query called in loop"
			log.Fatal(err.Error())
		}
		fmt.Println(person.Name, job.Name)
	}

}

