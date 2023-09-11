package service

import (
	"context"
	"database/sql"
	"strings"
)

const PATTERN_SQL_SELECT = `(?i)^\s*select\s*`

type DataSet struct {
	Results *[]interface{}
}

func NewDataSet() *DataSet {
	ds := make([]interface{}, 0)
	return &DataSet{Results: &ds}
}

func (ds *DataSet) batch(db *sql.DB, sql string) {
	if strings.TrimSpace(sql) == "" {
		return
	}

	rs := NewRegexText(PATTERN_SQL_SELECT, sql)
	if rs.IsMatch() {
		ds.query(db, sql)
	} else {
		ds.exec(db, sql)
	}
}
func (ds *DataSet) query(db *sql.DB, sql string) {
	dt := NewDataTable()
	if dt, err := dt.Fill(db, sql); err == nil {
		*ds.Results = append(*ds.Results, *dt)
	} else {
		*ds.Results = append(*ds.Results, err.Error())
	}
}
func (ds *DataSet) exec(db *sql.DB, sql string) {
	//====================================
	var count int64
	if res, err := db.Exec(sql); err == nil {
		if i, err := res.LastInsertId(); err == nil {
			count = count + i
		}
		if i, err := res.RowsAffected(); err == nil {
			count = count + i
		}
		*ds.Results = append(*ds.Results, count)
	} else {
		*ds.Results = append(*ds.Results, err.Error())
	}
}

// =============================================
// public methods
func (ds *DataSet) Exec(req *DBRequest) {
	di := NewConnection(req)
	db, err := di.Open()
	if err != nil {
		*ds.Results = append(*ds.Results, err.Error())
		return
	}
	defer db.Close()
	//=============================================
	var tx *sql.Tx
	if req.Trans {
		tx, err = db.BeginTx(context.Background(), nil)
		if err != nil {
			*ds.Results = append(*ds.Results, err.Error())
			return
		}
		defer tx.Rollback()
	}
	//=============================================
	ds.batch(db, req.Sql)
	for _, sql := range req.Sqls {
		ds.batch(db, sql)
	}
	//=============================================
	if req.Trans && tx != nil {
		tx.Commit()
	}
}
