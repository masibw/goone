![test_and_lint](https://github.com/masibw/go_one/workflows/test_and_lint/badge.svg)

# go_one
go_one finds N+1(strictly speaking call SQL in a for loop) query in go 

## Example
```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
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

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Name,&person.JobID); err != nil {
			log.Fatal(err)
		}

		var job Job

        // This is N+1 query
		if err := cnn.QueryRow("SELECT job_id, name FROM Jobs WHERE job_id = ?",person.JobID).Scan(&job.JobID,&job.Name); err != nil { 
			log.Fatal(err)
		}
		fmt.Println(person.Name,job.Name)
	}

}
```

## output
```
./hoge.go:38:13: this query is called in a loop
```

# Install
```
go get github.com/masibw/goone/cmd/go_one
```

# Usage

## bash
```
go vet -vettool=`which go_one` ./...
```

## fish
```
go vet -vettool=(which go_one) ./...
```


## CI
### Github Actions
```
- name: install goone
    run: go get -u github.com/masibw/goone/cmd/go_one
- name: run goone
    run: go vet -vettool=`which go_one` -go_one.configPath="$PWD/go_one.yml" ./...
```

# Library Support
- sql
- sqlx
- gorp
- gorm

You can add types to detect sql query.

# Options
You can use the `-go_one.configPath` option at runtime to determine if you want to use a specified types.

## Example

If go_one.yml exists in the directory where the command was executed
```
go vet -vettool=`which go_one` -go_one.configPath="$PWD/go_one.yml" ./...
```

You can also detect the case where an interface is in between by writing below. example [project](https://github.com/masibw/go_todo)
```yaml:go_one.yml
package:
  - pkgName: 'github.com/masibw/go_todo/cmd/go_todo/infrastructure/api/handler'
    typeNames:
      - typeName: '*todoHandler'
```

# Contribute
You're welcome to build an Issue or create a PR and be proactive!

# Caution
This tool does not support calls to functions from other packages in a for loop

