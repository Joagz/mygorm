package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

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
		for _, v := range m.References {
			columnRefLength += len(v.RefColumns)
			joins += fmt.Sprintf("JOIN %s ON %s.%s = %s.%s ", v.RefTable, v.RefTable, v.ForeignColumn, m.Table, v.LocalColumn)
			cols += "," + arrayToCommaSeparatedTable(v.RefColumns, v.RefTable)
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
		selectStr += fmt.Sprintf("%s = ? ", m.Table+"."+v)
	}
	return selectStr, columnRefLength
}

func storeMySqlResult(row *sql.Row, readLength int) ([]any, error) {
	result := make([]any, readLength)
	dest := make([]any, readLength)

	for i := range dest {
		dest[i] = &result[i]
	}

	if err := row.Scan(dest...); err != nil {
		return result, err
	}

	return result, nil
}

func storeMySqlResultSet(resultSetPtr *[][]any, rows *sql.Rows, readLength int) {
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

func (m mySqlModel) FindAll() (resultSet [][]any, err error) {
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

func (m mySqlModel) FindById(id int) (result []any, err error) {
	err = m.Open()
	if err != nil {
		return nil, err
	}
	defer m.Close()

	pk := m.PrimaryKey

	selectStr, columnRefLength := m.makeMySqlSelectWhere(pk)

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
	result, err = storeMySqlResult(row, readLength)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m mySqlModel) FindBy(params map[string]any) (resultSet [][]any, err error) {
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

func appendMySqlInsertValues(field reflect.StructField, value reflect.Value, m mySqlModel) (string, error) {
	switch value.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		props := field.Tag.Get(PropertyListName)
		properties := strings.Split(props, ",")
		if slices.Contains(properties, PrimaryKeyPropertyName) && 
			slices.Contains(properties, AutoIncrementPropertyName) {
			return "", nil
		}

		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32), nil
	case reflect.String:
		return fmt.Sprintf("'%s'", value.String()), nil
	case reflect.Bool:
		val := value.Bool()

		if val {
			return "TRUE", nil
		} else {
			return "FALSE", nil
		}

	default:
		ref := field.Tag.Get(ReferenceKeyPropertyName)
		found := false
		if ref != "" {
			ftype := value.Type()
			r := m.References[ref]

			for j := 0; j < value.NumField(); j++ {
				field := ftype.Field(j)
				if tag := field.Tag.Get(ColumnKeyPropertyName); tag == r.ForeignColumn {
					found = true
					return appendMySqlInsertValues(field, value.Field(j), m)
				}
			}
		}

		if !found {
			return "", fmt.Errorf("could not append values: invalid data type")
		}
	}

	return "", nil
}

func (m mySqlModel) Insert(data any) error {

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	var cols []string 
	format := "INSERT INTO %s (%s) VALUES (%s)"

	valuesStr := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if col := field.Tag.Get(ColumnKeyPropertyName); col != "" {

			if ref := field.Tag.Get(ReferenceKeyPropertyName); ref != "" {
				r := m.References[ref]
				cols = append(cols, r.LocalColumn)
			} else {

				props := strings.Split(strings.Trim(field.Tag.Get(PropertyListName), " "), ",")
				
				if len(props) > 0 {
					if slices.Contains(props, PrimaryKeyPropertyName) &&
						slices.Contains(props, AutoIncrementPropertyName) {
						fmt.Printf("skip pk\n")
						continue
					}
				}

				cols = append(cols, col)

			}

			str, err := appendMySqlInsertValues(field, value, m)
			
			if err != nil {
				return err
			}

			valuesStr += str

			// skip last item
			if i < t.NumField()-1 {
				valuesStr += ","
			}
		}
	}

	insertStr := fmt.Sprintf(format, m.Table, arrayToCommaSeparatedTable(cols, m.Table), valuesStr)
	fmt.Printf("insertStr: %v\n", insertStr)

	err := m.Open()
	if err != nil {
		return err
	}
	defer m.Close()

	stmt, err := m.DB.Prepare(insertStr)
	if err != nil {
		return err 
	}

	if _,err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (mySqlModel) UpdateById(data any, id int) error {
	panic("mySqlModel: function not implemented")
}

func (mySqlModel) UpdateBy(data any, params ...any) error {
	panic("mySqlModel: function not implemented")
}
