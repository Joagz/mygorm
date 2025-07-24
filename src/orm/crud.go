package orm

import (
	"database/sql"
	"fmt"

	"golang.org/x/exp/maps"
)

/*
* We implement these for different drivers
 */
func (model) FindAll() (resultSet []any, err error) {
	panic("please select a sql engine")
}
func (model) FindById(id int) (result any, err error) {
	panic("please select a sql engine")
}
func (model) FindBy(map[string]any) (resultSet []any, err error) {
	panic("please select a sql engine")
}
func (model) NumRows() (n int) {
	panic("please select a sql engine")
}
func (model) Insert(data any) error {
	panic("please select a sql engine")
}
func (model) UpdateById(data any, id int) error {
	panic("please select a sql engine")
}
func (model) UpdateBy(data any, params ...any) error {
	panic("please select a sql engine")
}

/*
 * MySQL Model Drivers
 */

type mySqlModel struct {
	model
	mySqlConnector
}

func (m mySqlModel) makeMySqlSelect() (result string, columnRefLength int) {

	cols := arrayToCommaSeparatedTable(m.Columns, m.Table)
	joins := ""

	if len(m.References) > 0 {
		keys := maps.Keys(m.References)
		for i,k := range keys {
			ref := m.References[k]
			columnRefLength += len(ref)
			cols += "," + arrayToCommaSeparatedTable(ref, k)
			joins += fmt.Sprintf("JOIN %s ON %s = %s ", k, m.ForeignKeys[i], fmt.Sprintf("%s.%s", k, ref[0]))
		}
	}
	result = fmt.Sprintf("SELECT %s FROM %s %s", cols, m.Table, joins)
	return result, columnRefLength
}

func (m mySqlModel) makeMySqlSelectWhere(conditions ...string) (result string, columnRefLength int) {

	selectStr, columnRefLength := m.makeMySqlSelect()
	selectStr = fmt.Sprintf("%s WHERE ", selectStr)
	
	for i, v := range conditions {
		if i > 0 {
			selectStr += "AND "
		}
		selectStr += fmt.Sprintf("%s = ? ", m.Table + "." + v)
	}	
	return selectStr, columnRefLength
}

func storeMySqlResultSet(resultSetPtr *[]any, rows *sql.Rows, readLength int) {
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

		*resultSetPtr = append(*resultSetPtr, rawResult)
	}
}

func (m mySqlModel) FindAll() (resultSet []any, err error) {
	err = m.Open()
	if err != nil {
		return nil, err
	}
	defer m.Close()

	selectStr, columnRefLength := m.makeMySqlSelect()
	
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
	
	storeMySqlResultSet(&resultSet, rows, readLength)
	
	return resultSet, nil
}

func (m mySqlModel) FindById(id int) (any, error) {
	err := m.Open()
	if err != nil {
		return nil, err
	}
	defer m.Close()

	pk := m.PrimaryKey

	selectStr, columnRefLength := m.makeMySqlSelectWhere(pk)

	fmt.Printf("selectStr: %v\n", selectStr)	
	
	stmt, err := m.DB.Prepare(selectStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	readLength := len(m.Columns) + columnRefLength
	result := make([]any, readLength)
	dest := make([]any, readLength)

	for i := range dest {
		dest[i] = &result[i]
	}

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	return result, err
}

func (m mySqlModel) FindBy(params map[string]any) (resultSet []any, err error) {
	err = m.Open()
	if err != nil {
		return nil, err
	}
	defer m.Close()

	keys := maps.Keys(params)
	values := maps.Values(params)

	selectStr, columnRefLength := m.makeMySqlSelectWhere(keys...)

	stmt, err := m.DB.Prepare(selectStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}

	readLength := len(m.Columns) + columnRefLength

	storeMySqlResultSet(&resultSet, rows, readLength)

	return resultSet, nil
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
