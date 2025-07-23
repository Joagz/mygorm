package main

import (
	"fmt"
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

func Init() orm.Connector {
	conn := orm.NewConnector("mysql", "root", "123456", "db_name", "localhost:3306")
	return conn
}

func TableModel(conn orm.Connector) (orm.Model, error) {
	return conn.ModelOf(Table{}, "test")
}

func main() {
	conn := Init()

	err := conn.Open()
	if err != nil {
		fmt.Printf("error connecting: '%s'", err.Error())
		return
	}
	defer conn.Close()
	
	model, err :=  TableModel(conn)
	if err != nil {
		fmt.Printf("error creating model: %s\n", err.Error())
		return 
	}
	model.Print()
	result, err := model.FindAll()
	if err != nil {
		fmt.Printf("ERROR : %s\n", err.Error())
		return
	}

	fmt.Printf("result: len() = %d\n", len(result))

	for _, v := range result {
		if arr, ok := v.([]interface{}); ok {
			for i, item := range arr {
				switch val := item.(type) {
				case []byte:
					fmt.Printf("[%d] string: %s\n", i, string(val))
				default:
					fmt.Printf("[%d] value: %v (type %T)\n", i, val, val)
				}
			}
			fmt.Println("----")
		}
	}

}
