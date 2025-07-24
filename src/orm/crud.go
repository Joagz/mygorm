package orm

import (
	"fmt"
	"maps"
)

/*
* We implement these for different drivers
 */
func (model) FindAll() (resultSet []any, err error) 		{ panic("please select a sql engine") }
func (model) FindById(id int) (result any, err error)           { panic("please select a sql engine") }
func (model) FindBy(params ...any) (resultSet []any, err error) { panic("please select a sql engine") }
func (model) NumRows() (n int)                                  { panic("please select a sql engine") }
func (model) Insert(data any) error                             { panic("please select a sql engine") }
func (model) UpdateById(data any, id int) error                 { panic("please select a sql engine") }
func (model) UpdateBy(data any, params ...any) error            { panic("please select a sql engine") }

/*
 * MySQL Model Drivers
 */

type mySqlModel struct {
	model
	mySqlConnector	
}

func (m mySqlModel) FindAll() (resultSet []any, err error) {
	err = m.Open()	
	if err != nil {
		return nil, err 
	}
	defer m.Close()

	cols := arrayToCommaSeparatedTable(m.Columns, m.Table)
	joins := ""
	columnRefLength := 0

	if len(m.References) > 0 {
		keys := maps.Keys(m.References)
		i := 0
		for k := range keys {
			ref := m.References[k]
			columnRefLength+=len(ref)
			cols += "," + arrayToCommaSeparatedTable(ref, k)
			joins += fmt.Sprintf("JOIN %s ON %s = %s ", k, m.ForeignKeys[i], fmt.Sprintf("%s.%s", k, ref[0]))
			i++
		}
	}
	selectStr := fmt.Sprintf("SELECT %s FROM %s %s %s", cols, m.Table, joins, "")

	stmt, err := m.DB.Prepare(selectStr)	
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	readLength := len(m.Columns) + columnRefLength
	for rows.Next() {
		/*
		 * We use dest to store pointers to each rawResult entry
		 */
		rawResult := make([]any, readLength)
		dest := make([]any, readLength)
		for i := range rawResult {
			dest[i] = &rawResult[i]
		}

		err := rows.Scan(dest...)
		if err != nil {
			fmt.Printf("error writing to buffer on FindAll: '%s'\n", err.Error())
			continue
		}
		
		resultSet = append(resultSet, rawResult)
	}
	return resultSet, nil
}

func (m mySqlModel) FindById(id int) (any, error) {
	err := m.Open()	
	if err != nil {
		return nil, err 
	}
	defer m.Close()

	pk := m.PrimaryKey
	cols := arrayToCommaSeparatedTable(m.Columns, m.Table)
	selectStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", cols, m.Table, pk)

	stmt, err := m.DB.Prepare(selectStr)	
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	if err := row.Err(); err != nil {
		return nil, err 
	}

	length := len(m.Columns)
	buffer := make([]any, length)
	dest := make([]any, length)

	for i := range dest {
		dest[i] = &buffer[i]	
	}

	err = row.Scan(dest...) 

	return buffer, err
}

func (mySqlModel) FindBy(params ...any) (resultSet []any, err error) {
	panic("mySqlModel: function not implemented")
}

func (mySqlModel) NumRows() (n int) {
	panic("mySqlModel: function not implemented")
}

func (mySqlModel) Insert(data any) error {
	panic("mySqlModel: function not implemented")
}

func (mySqlModel) UpdateById(data any, id int) error {
	panic("mySqlModel: function not implemented")
}

func (mySqlModel) UpdateBy(data any, params ...any) error {
	panic("mySqlModel: function not implemented")
}



