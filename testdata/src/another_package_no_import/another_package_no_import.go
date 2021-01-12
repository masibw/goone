package another_package_no_import

import (
	"fmt"

	"github.com/masibw/goone_test/pkg/db"
)


func main() {

	for i:=0;i<10;i++ {
		person := db.Person{Name:"a",JobID: 1}
		job := db.GetJobInsideCnn(person) //want "this query is called in a loop"
		fmt.Println(person.Name, job.Name)
	}

}
