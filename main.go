package main

import (
	"fmt"
)

type Table struct {
	ID      int        `col:"id" props:"primary-key"`
	Column1 string     `col:"col_1"`
	Column2 string     `col:"col_2"`
	Other   OtherTable `col:"other_id" props:"foreign-key" ref:"other"`
	Extra   ExtraTable `col:"extra_id" props:"foreign-key" ref:"extra"`
}
type ExtraTable struct {
	ID      int        `col:"id" props:"primary-key"`
	Column1 string     `col:"extra_col_1"`
}
type OtherTable struct {
	ID      int        `col:"id" props:"primary-key"`
	Column1 string     `col:"other_col_1"`
	Column2 string     `col:"other_col_2"`
}

/* convenient implementation */
func TableModel() (model.Model, error) {
	return model.Of(Table{}, "table")
}

func main() {
	m, err := TableModel()
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	m.Print() // print model
	rows, err := m.FindAll() // example usage
	fmt.Printf("rows: %v\n", rows)
}
