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
	"os"
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
		if _, err := os.Stat(output); os.IsNotExist(err) {
			os.MkdirAll(output, os.ModePerm)
		}

		username := "Ryan"
		u, err := user.Current()
		if err == nil {
			username = u.Name
		}

		for _, t := range ts {
			goGen(schema, t, username)
		}
	},
}

var db *sql.DB

func goGen(schema, table, username string) error {
	f := NewFilePath(fmt.Sprintf("%s/%s/model", "pay", "wxpay"))
	f.Comment(fmt.Sprintf("%s ", strcase.ToCamel(table)))

	f.Comment(fmt.Sprintf("Created by %s at %s", username, time.Now().Format("2006-01-02 15:04:05")))
	cols, err := parseTable(schema, table)
	if err != nil {
		return err
	}

	var fields []Code

	for _, v := range cols {
		fields = append(fields, v.toField())
	}
	f.Type().Id(strcase.ToCamel(table)).Struct(fields...)

	// func TableName
	f.Comment("TableName set model table name")
	f.Func().Params(
		Id("m").Op("*").Id(strcase.ToCamel(table)),
	).Id("TableName").Params().String().Block(
		Return(Lit(table)),
	)

	// ID
	f.Comment("ID return model id")
	f.Func().Params(
		Id("m").Op("*").Id(strcase.ToCamel(table)),
	).Id("ID").Params().Id("uint64").Block(
		Return(Id("m.Id")),
	)

	// BeforeUpdate
	f.Comment("BeforeUpdate set updateAt before update db in gorm")
	f.Func().Params(
		Id("m").Op("*").Id(strcase.ToCamel(table)),
	).Id("BeforeUpdate").Params().Error().Block(
		Id("now").Op(":=").Qual("time", "Now").Call(),
		Id("m.UpdatedAt").Op("=").Op("&").Id("now"),
		Return(Id("nil")),
	)

	fmt.Printf("%#v\n", f)

	ioutil.WriteFile(path.Join(output, fmt.Sprintf("%s.go", table)), []byte(f.GoString()), os.ModePerm)
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
		if strings.Contains(strings.ToLower(self.colType), "unsigned") {
			s.Uint64()
		} else {
			s.Int64()
		}
	case "char", "varchar", "text", "nvarchar", "nchar", "mediumtext", "json", "longtext":
		s.String()
	case "blob":
		s.Index().Byte()
	case "int", "smallint", "mediumint":
		if strings.Contains(strings.ToLower(self.colType), "unsigned") {
			s.Uint()
		} else {
			s.Int()
		}
	case "tinyint":
		s.Int8()
	case "float":
		s.Float32()
	case "double":
		s.Float64()
	case "time", "datetime", "timestamp", "date":
		s.Op("*").Qual("time", "Time")
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
