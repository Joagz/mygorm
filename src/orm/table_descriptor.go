package orm

import (
	"fmt"
)

// todo: modify reference in model
type reference struct {
	RefTable      string   // name of the referenced table
	RefColumns    []string // columns of the referenced table
	LocalColumn   string   // name of the column in the current table referencing ForeignColumn
	ForeignColumn string   // name of the column referenced by LocalColumn
}

type model struct {
	config      datasourceConfig
	Table       string
	Columns     []string
	PrimaryKey  string
	References  map[string]reference
}

func (m model) GetColumns() []string {
	return m.Columns
}

func (m model) GetPrimary() string {
	return m.PrimaryKey
}	

// This function is horrible and beautiful at the same time
func (m model) Print() {
	fmt.Printf("name: %s\n", m.Table)
	fmt.Printf("pk: %s\n", m.PrimaryKey)
	for i, s := range m.Columns {
		fmt.Printf("col [%d]: %s ", i, s)
	}
	fmt.Println("\nreferences:")
	for _, v := range m.References {
		fmt.Printf("\ttable: %s\n\t", v)
		for i, s := range v.RefColumns {
			fmt.Printf("col [%d]: %s ", i, s)
		}
		fmt.Println("")
	}
}
