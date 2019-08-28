package parser

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"log"
	"os/user"
	"path"
	"text/template"
	"time"

	"database/sql"
	"os"
)

var db *sql.DB

func init() {
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		"apuser", "airparking", "10.35.22.61:3306", "airparking")
	//conn := "apuser:airparking@tcp(10.35.22.61:3306)/airparking"
	//viper.GetString("db.user"),
	//viper.GetString("db.password"),
	//viper.GetString("addr"),
	//viper.GetString("schema"))
	var err error
	db, err = sql.Open("mysql", conn)

	if err != nil {
		os.Exit(0)
	}
}

func parseSchema(schema, projectPkg, modelPkg string) ([]*Table, error) {
	sql := "SELECT DISTINCT `TABLE_NAME` FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `COLUMN_NAME` = 'id'"
	rows, err := db.Query(sql, schema)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var tables []*Table
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			log.Println(err)
			return nil, err
		}
		var t *Table
		t, err = parseTable(schema, tableName, projectPkg, modelPkg)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}

func parseTable(schema, tableName, projectPkg, modelPkg string) (*Table, error) {
	sql := "SELECT `COLUMN_NAME`,`DATA_TYPE`,`COLUMN_COMMENT`,`COLUMN_KEY`,`COLUMN_TYPE`,`IS_NULLABLE` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE`TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	rows, err := db.Query(sql, schema, tableName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var cols []*column

	for rows.Next() {
		var columnName, dataType, comment, colKey, colType, nullable string
		if err = rows.Scan(&columnName, &dataType, &comment, &colKey, &colType, &nullable); err != nil {
			log.Println(err)
			return nil, err
		}

		var col *column
		col, err = newColumn(columnName, dataType, comment, colType, isPrimaryKey(colKey), isNullable(nullable))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		cols = append(cols, col)
	}

	if len(cols) == 0 {
		return nil, errors.New(fmt.Sprintf("Cannot found Table %s", tableName))
	}

	var u *user.User
	u, err = user.Current()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Table{
		schema:     schema,
		TableName:  tableName,
		Columns:    cols,
		ProjectPkg: projectPkg,
		ModelPkg:   modelPkg,
		ModelName:  strcase.ToCamel(tableName),
		Author:     u.Name,
		DateStr:    time.Now().Format("2006-01-02 15:04:05"),
	}, nil

}

func parseTemplate(tplName, tplPath string) (*template.Template, error) {
	//templateFile := path.Join(tplPath, tplName)
	return template.New(tplName).Funcs(template.FuncMap{"ToLowerCamel": strcase.ToLowerCamel}).ParseFiles(path.Join(tplPath, tplName))
	//return template.ParseFiles(path.Join(tplPath, tplName))
}
