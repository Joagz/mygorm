package main

import (
	"fmt"
	"mygorm/src/orm"
	"os"
	"text/tabwriter"
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
	Column1 string     `col:"other_col_1"`
	Column2 string     `col:"other_col_2"`
	ID      int        `col:"id" props:"primary-key"`
}

func Init() orm.Connector {
	conn := orm.NewConnector("mysql", "root", "123456", "db_name", "localhost:3307")
	return conn
}

func TableModel(conn orm.Connector) (orm.Model, error) {
	return conn.ModelOf(Table{}, "test")
}

func main() {
	conn := Init()
		
	model, err :=  TableModel(conn)
	if err != nil {
		fmt.Printf("error creating model: %s\n", err.Error())
		return 
	}
	
	result, err := model.FindAll()
	if err != nil {
		fmt.Printf("ERROR : %s\n", err.Error())
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "Row\tID\tCol_1\tCol_2\tOther_ID\tOther_Col_1\tOther_Col_2\tExtra_ID\tExtra_Col_1")

	for i, v := range result {
		if arr, ok := v.([]any); ok {
			fmt.Fprintf(w, "%d\t%d\t%s\t%s\t%d\t%s\t%s\t%d\t%s\n",
				i,
				arr[0],
				arr[1],
				arr[2],
				arr[3],
				arr[4],
				arr[5],
				arr[6],
				arr[7],
			)
		}
	}

	w.Flush()
}
