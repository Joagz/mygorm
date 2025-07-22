package model

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

const (
	CommaSeparatedRegex = `^(\s*[\w\s-]+\s*)(,\s*[\w\s-]+\s*)*$` 
	PrimaryKeyPropertyName = "primary-key"
	ForeignKeyPropertyName = "foreign-key"
)

type Model interface {
	FindAll() ([]any, error)
	// FindById(int) (any, error)
	// FindBy(...any) ([]any, error)
	// NumRows() int

	Print()
	GetReference() map[string][]string
	GetColumns() []string
}

func Of(d any, tablename string) (Model, error) {
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

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("col")
		m.Columns = append(m.Columns, col)

		if props := field.Tag.Get("props"); props != "" {
			ok := expr.MatchString(props)
			if !ok {
				return nil, fmt.Errorf("regular expression checking failed for %s", props)	
			}

			// treat as comma-separated string 
			properties := strings.Split(props, ",")
			if slices.Contains(properties, PrimaryKeyPropertyName) {
				m.PrimaryKey = col
			}

			if slices.Contains(properties, ForeignKeyPropertyName) {
				ref := field.Tag.Get("ref")
				if ref == ""{
					return nil, fmt.Errorf("please set `ref` for column %s", col)	
				}

				m.ForeignKeys = append(m.ForeignKeys, col)
				foreign := field.Type
				fmodel, err := Of(foreign, ref)	
				if err != nil {
					return nil, err
				}
				if(m.References == nil) {
					m.References = make(map[string][]string)
				}
				m.References[ref] = fmodel.GetColumns()
			}
		}
	}

	m.Table = tablename

	return &m, nil
}
