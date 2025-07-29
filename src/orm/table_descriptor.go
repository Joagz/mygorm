package orm

import (
	"fmt"
)

// todo: modify reference in model
type reference struct {
	RefTable      string   // name of the referenced table
	RefColumns    []column // columns of the referenced table
	LocalColumn   string   // name of the column in the current table referencing ForeignColumn
	ForeignColumn string   // name of the column referenced by LocalColumn
}

type column struct {
	Name             string
	Properties       []string
	References       string
	ReferencedColumn string
	IsPrimary        bool
	IsAutoIncrement  bool
	IsForeign        bool
}

type model struct {
	config     datasourceConfig
	Table      string
	Columns    []column
	PrimaryKey string
	References map[string]reference
}

func (m model) GetColumns() []column {
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
		fmt.Printf("col [%d]: %v ", i, s)
	}
	fmt.Println("\nreferences:")
	for _, v := range m.References {
		fmt.Printf("\ttable: %v\n\t", v)
		for i, s := range v.RefColumns {
			fmt.Printf("col [%d]: %v ", i, s)
		}
		fmt.Println("")
	}
}
