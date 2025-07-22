package main

import (
	"mygorm/src/orm"
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

func Connect() orm.Connector {
	conn := orm.NewConnector("mysql", "root", "123456", "db_name", "localhost:3306")

	return conn
}

func main() {
	conn := Connect()
	
	model, err := conn.ModelOf(Table{}, "table")

	if err != nil {
		panic("error making model")
	}

	model.Print()


}
