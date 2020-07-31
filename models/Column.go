package models

type Column struct {
	ColumnName    string
	ColumnType    string
	ColumnDefault string
	IsNullable    string
	Extra         string
	ColumnComment string
}
