package model

import (
	"fmt"
)

type model struct {
	Table       string
	Columns     []string
	PrimaryKey  string
	ForeignKeys []string
	References  map[string][]string
}

func (m model) GetReference() map[string][]string {
	ref := make(map[string][]string)
	ref[m.Table] = m.Columns
	return ref
}

func (m model) GetColumns() []string {
	return m.Columns
}

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

func (model) FindAll() ([]any, error) {
	panic("no sql engine for model.FindAll")
}


