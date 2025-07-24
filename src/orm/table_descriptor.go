package orm

import (
	"fmt"
)

// todo: modify reference in model
type reference struct {
	RefTable string 	// name of the referenced table
	RefColumns []string	// columns of the referenced table 
	LocalColumn string	// name of the column in the current table referencing ForeignColumn
	ForeignColumn string 	// name of the column referenced y LocalColumn
}

type model struct {
	config      datasourceConfig
	Table       string
	Columns     []string
	PrimaryKey  string
	ForeignKeys []string
	References  map[string][]string
}

func (m model) GetColumns() []string {
	return m.Columns
}

// This function is horrible and beautiful at the same time
func (m model) Print() {
	fmt.Printf("name: %s\n", m.Table)
	fmt.Printf("pk: %s\n", m.PrimaryKey)
	for i, s := range m.Columns {
		fmt.Printf("col [%d]: %s ",i, s)
	}
	for _, s := range m.ForeignKeys {
		fmt.Printf("fk: %s ", s)
	}
	fmt.Println("\nreferences:")
	for k, v := range m.References {	
		fmt.Printf("\ttable: %s\n\t",k)
		for i, s := range v {
			fmt.Printf("col [%d]: %s ",i, s)
		}
		fmt.Println("")
	}
}

