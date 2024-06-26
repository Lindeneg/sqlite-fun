package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xwb1989/sqlparser"
)

// https://www.sqlite.org/fileformat.html

var t int64
var timing bool = false

func main() {
	if len(os.Args) < 3 {
		log.Fatal("please provide arguments: file command")
	}
	if len(os.Args) > 3 && os.Args[3] == "-t" {
		timing = true
		t = time.Now().UnixMilli()
	}
	databaseFile := os.Args[1]
	cmd := os.Args[2]
	db, err := newDatabaseFile(databaseFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.File.Close()
	switch cmd {
	case ".dbinfo":
		fmt.Printf("database page size: \t%v\n", db.Header.PageSize)
		fmt.Printf("number of tables: \t%v\n", len(db.Tables))
		break
	case ".tables":
		fmt.Println(strings.Join(db.TableNames(), " "))
	case ".roots":
		fmt.Println(db)
	default:
		stmt, err := sqlparser.Parse(cmd)
		if err != nil {
			log.Fatal("unknown command/query: " + cmd)
		}
		switch stmt := stmt.(type) {
		case *sqlparser.Select:
			HandleSelect(NewSelectCtx(stmt), db)
		}
	}
	if timing {
		diff := float64(time.Now().UnixMilli() - t)
		fmt.Println(diff/1000, "seconds")
	}

}
