package models

type File struct {
	Name  string
	Size  int64
	Date  string
	Type  string
	Url   string
	IsDir bool
}
