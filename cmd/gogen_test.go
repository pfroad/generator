package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGoGen(t *testing.T) {
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		"apuser", "airparking", "10.35.22.61:3306", "airparking")
	var err error
	db, err = sql.Open("mysql", conn)

	if err != nil {
		os.Exit(0)
	}
	// GoGen()
}
