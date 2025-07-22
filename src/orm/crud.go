package orm

/*
* We implement these for different drivers
*/
func (model) FindAll() (resultSet []any, err error) 		{ panic("model: function not implemented") }
func (model) FindById(id int) (result any, err error)           { panic("model: function not implemented") }
func (model) FindBy(params ...any) (resultSet []any, err error) { panic("model: function not implemented") }
func (model) NumRows() (n int)                                  { panic("model: function not implemented") }
func (model) Insert(data any) error                             { panic("model: function not implemented") }
func (model) UpdateById(data any, id int) error                 { panic("model: function not implemented") }
func (model) UpdateBy(data any, params ...any) error            { panic("model: function not implemented") }

type mySqlModel struct {
	model
}

func (mySqlModel) FindAll() (resultSet []any, err error) 		{ panic("mySqlModel: function not implemented") }
func (mySqlModel) FindById(id int) (result any, err error)           	{ panic("mySqlModel: function not implemented") }
func (mySqlModel) FindBy(params ...any) (resultSet []any, err error) 	{ panic("mySqlModel: function not implemented") }
func (mySqlModel) NumRows() (n int)                                  	{ panic("mySqlModel: function not implemented") }
func (mySqlModel) Insert(data any) error                             	{ panic("mySqlModel: function not implemented") }
func (mySqlModel) UpdateById(data any, id int) error                 	{ panic("mySqlModel: function not implemented") }
func (mySqlModel) UpdateBy(data any, params ...any) error            	{ panic("mySqlModel: function not implemented") }



