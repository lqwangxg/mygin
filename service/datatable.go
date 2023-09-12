package service

import (
	"database/sql"
)

type DataTable struct {
	Name     string
	Columns  []DataColumn
	Rows     []DataRow
	ColCount int
	RowCount int
	Comment  string
	Owner    string
}
type DataRow []interface{}

type DataColumn struct {
	Index     int
	Name      string
	Type      string
	Length    int64
	ScanType  string
	Comment   string
	Nullable  bool
	Precision int64
	Scale     int64
}

func NewDataTable() *DataTable {
	return &DataTable{Name: "datatable"}
}
func NewDataTableWithName(name string) *DataTable {
	return &DataTable{Name: name}
}

// ===============================================
// public methods
func (dataTable *DataTable) Fill(db *sql.DB, strSQL string) (*DataTable, error) {
	if dataTable == nil {
		dataTable = NewDataTable()
	}
	//====================================
	// query
	rs, err := db.Query(strSQL)
	if err != nil {
		return dataTable, err
	}
	defer rs.Close()

	//====================================
	// fetch data result to datatable
	var cols []string
	if cols, err = rs.Columns(); err != nil {
		return dataTable, err
	}
	var colTypes []*sql.ColumnType
	if colTypes, err = rs.ColumnTypes(); err != nil {
		return dataTable, err
	}
	dataTable.InitColumns(cols, colTypes)
	for rs.Next() {
		address := dataTable.NewRow()
		if err := rs.Scan(address...); err != nil {
			return dataTable, err
		}
	}
	dataTable.RowCount = len(dataTable.Rows)
	return dataTable, nil
}

// PutFields ...
func (dt *DataTable) InitColumns(columns []string, colTypes []*sql.ColumnType) {
	dt.ColCount = len(columns)
	dt.Columns = make([]DataColumn, dt.ColCount)
	for idx, name := range columns {
		dt.Columns[idx].Index = idx
		dt.Columns[idx].Name = name
		ctype := colTypes[idx]
		dt.Columns[idx].Type = ctype.DatabaseTypeName()
		if length, ok := ctype.Length(); ok {
			dt.Columns[idx].Length = length
		}
		if prec, scale, ok := ctype.DecimalSize(); ok {
			dt.Columns[idx].Precision = prec
			dt.Columns[idx].Scale = scale
		}
	}
}
func (dt *DataTable) NewRow() []interface{} {
	cells := make([]interface{}, dt.ColCount)
	addres := make([]interface{}, dt.ColCount)
	//address point to value
	for i := range addres {
		addres[i] = &cells[i]
	}
	if dt.Rows == nil {
		dt.Rows = make([]DataRow, 0)
	}
	dt.Rows = append(dt.Rows, cells)
	return addres
}
