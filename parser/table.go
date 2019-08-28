package parser

type Table struct {
	schema    string
	TableName string
	comment   string
	Columns   []*column

	ProjectPkg string

	ModelName string
	ModelPkg  string

	DateStr string
	Author  string
}
