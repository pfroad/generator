/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	. "github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

// gogenCmd represents the gogen command
var gogenCmd = &cobra.Command{
	Use:   "gogen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
generator gogen --schema=airparking --tables=ap_user
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gogen called")
		ts := strings.Split(tables, ",")
		for _, t := range ts {
			goGen(schema, t)
		}
	},
}

var db *sql.DB

func goGen(schema, table string) error {
	f := NewFilePath(fmt.Sprintf("%s/%s/model", "pay", "wxpay"))
	f.Comment(fmt.Sprintf("%s ", strcase.ToCamel(table)))

	u, err := user.Current()
	if err != nil {
		log.Println(err)
		return err
	}
	f.Comment(fmt.Sprintf("Created by %s at %s", u.Name, time.Now().Format("2006-01-02 15:04:05")))
	cols, err := parseTable(schema, table)
	if err != nil {
		return err
	}

	var fields []Code

	for _, v := range cols {
		fields = append(fields, v.toField())
	}
	f.Type().Id(strcase.ToCamel(table)).Struct(fields...)

	fmt.Printf("%#v\n", f)

	ioutil.WriteFile(path.Join(output, "model"))
	return nil
}

type column struct {
	columnName   string
	dataType     string
	comment      string
	colType      string
	isNullable   bool
	isPrimaryKey bool
}

func newColumn(columnName, dataType, comment, colType string, isPrimaryKey, isNullable bool) *column {
	return &column{
		columnName:   columnName,
		dataType:     dataType,
		comment:      comment,
		colType:      colType,
		isPrimaryKey: isPrimaryKey,
		isNullable:   isNullable,
	}
}

func (self *column) toField() *Statement {
	s := Id(strcase.ToCamel(self.columnName))
	switch self.dataType {
	case "bigint":
		s.Int64()
	case "char", "varchar", "text", "nvarchar", "nchar", "mediumtext", "json", "longtext":
		s.String()
	case "blob":
		s.Index().Byte()
	case "int", "tinyint", "smallint", "mediumint":
		s.Int()
	case "float":
		s.Float32()
	case "double":
		s.Float64()
	case "time", "datetime", "timestamp", "date":
		s.Op("*").Qual("time", "time")
	case "bit", "boolean":
		s.Bool()
	default:
		fmt.Errorf("unsupport data type %s", self.dataType)
	}

	var jsonTag string
	if !self.isNullable {
		jsonTag = self.columnName
	} else {
		jsonTag = fmt.Sprintf("%s,emitempty", self.columnName)
	}
	s.Tag(map[string]string{"gorm": fmt.Sprintf("column:%s", self.columnName), "json": jsonTag})
	s.Comment(self.comment)
	return s
}

func parseTable(schema, tableName string) ([]*column, error) {
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
		col = newColumn(columnName, dataType, comment, colType, isPrimaryKey(colKey), isNullable(nullable))
		cols = append(cols, col)
	}

	if len(cols) == 0 {
		return nil, errors.New(fmt.Sprintf("Cannot found table %s", tableName))
	}

	return cols, err
}

func isPrimaryKey(colKey string) bool {
	return colKey == "PRI"
}

func isNullable(nullable string) bool {
	return nullable == "YES"
}
