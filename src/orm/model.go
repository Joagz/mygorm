package orm 

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

const (
	CommaSeparatedRegex       = `^(\s*[\w\s-]+\s*)(,\s*[\w\s-]+\s*)*$`
	PrimaryKeyPropertyName    = "primary-key"
	AutoIncrementPropertyName = "auto-increment"
	ForeignKeyPropertyName    = "foreign-key"
	ReferenceKeyPropertyName  = "ref"
	ColumnKeyPropertyName     = "col"
	PropertyListName          = "props"
)


type Model interface {
	/*
	 * SQL Crud Methods
	 */
	FindAll() ([][]any, error)
	FindById(int) ([]any, error)
	FindBy(map[string]any) ([][]any, error)
	NumRows() int
	Insert(any) error
	UpdateById(any, int) error
	UpdateBy(any, ...any) error
	
	/*
	 * Other methods
	 */
	Print()
	GetColumns() []string
	GetPrimary() string
}
// We have to inject some configuration here
func (c connector) ModelOf(d any, tablename string) (Model, error) {
	var t reflect.Type
	if k, ok := d.(reflect.Type); !ok {
		t = reflect.TypeOf(d)
	} else {
		t = k
	}

	expr, err := regexp.Compile(CommaSeparatedRegex)
	if err != nil {
		panic("model.Of: INVALID CommaSeparatedRegex")
	}

	var m model

	var cols []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("col")

		hasPk := false
		hasFk := false 

		if props := field.Tag.Get(PropertyListName); props != "" {
			ok := expr.MatchString(props)
			if !ok {
				return nil, fmt.Errorf("regular expression checking failed for %s", props)
			}

			// treat as comma-separated string, must use a struct for this later 
			properties := strings.Split(props, ",")
			if slices.Contains(properties, PrimaryKeyPropertyName) {
				m.PrimaryKey = col
				// always append pk first
				hasPk = true

			}

			if slices.Contains(properties, ForeignKeyPropertyName) {
				ref := field.Tag.Get("ref")
				if ref == "" {
					return nil, fmt.Errorf("please set `ref` for column %s", col)
				}

				foreign := field.Type
				fmodel, err := c.ModelOf(foreign, ref)

				if err != nil {
					return nil, fmt.Errorf("could not get model for %s", ref)
				}

				if m.References == nil {
					m.References = make(map[string]reference)
				}

				r := reference {
					RefTable: ref,
					RefColumns: fmodel.GetColumns(),
					LocalColumn: col,
					ForeignColumn: fmodel.GetPrimary(),
				}

				m.References[ref] = r

				hasFk = true
			}
		}

		if hasPk {
			m.Columns = append(m.Columns, col)
		} else if !hasFk {
			cols = append(cols, col)
		}
	}

	// manually copy the rest of the columns 
	for _, v := range cols {
		m.Columns = append(m.Columns, v)
	}
	
	m.Table = tablename

	switch c.Engine {
	case "mysql":
		return &mySqlModel{
			model: m,
			mySqlConnector: mySqlConnector{connector: c, DB: nil}, 
		}, nil
	default:
		panic("please select a valid engine")

	}


}
