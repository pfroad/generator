package parser

type table struct {
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
