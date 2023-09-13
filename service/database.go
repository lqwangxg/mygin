package service

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	go_ora "github.com/sijms/go-ora/v2"
)

type DBConnectInfo struct {
	Driver     string
	DataSource string
	db         *sql.DB
}

func (di *DBConnectInfo) Open() (*sql.DB, error) {
	db, err := sql.Open(di.Driver, di.DataSource)
	if err == nil {
		di.db = db
	}
	return db, err
}

// ===============================================
// private methods
func NewConnection(req *DBRequest) *DBConnectInfo {
	di := &DBConnectInfo{Driver: req.Driver}
	switch req.Driver {
	case "oracle":
		di.DataSource = go_ora.BuildUrl(req.Host, req.Port, req.Name, req.User, req.Pass, nil)
	case "postgres":
		di.DataSource = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			req.Host, req.Port, req.User, req.Pass, req.Name)
	case "mysql":
		//id:password@tcp(your-amazonaws-uri.com:3306)/dbname
		di.DataSource = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			req.User, req.Pass, req.Host, req.Port, req.Name)
	}
	return di
}

// ===============================================
// public methods
func (di *DBConnectInfo) QueryTable(strSQL string) (*DataTable, error) {
	//====================================
	// open database
	db, err := di.Open()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()
	//====================================
	// fetch data result to datatable
	dataTable := NewDataTable()
	dataTable.Fill(db, strSQL)

	return dataTable, nil
}

// ===============================================
// public methods
func (di *DBConnectInfo) Exec(strSQL string) *ExecResult {
	//====================================
	var count int64
	result := &ExecResult{Sql: strSQL}
	if res, err := di.db.Exec(strSQL); err == nil {
		if i, err := res.LastInsertId(); err == nil {
			count = count + i
		}
		if i, err := res.RowsAffected(); err == nil {
			count = count + i
		}
		result.Count = count
	} else {
		result.Error = err.Error()
	}
	return result
}

func (di *DBConnectInfo) batch(sql string) *ExecResult {
	if strings.TrimSpace(sql) == "" {
		return nil
	}

	rs := NewRegexText(PATTERN_SQL_SELECT, sql)
	if rs.IsMatch() {
		return di.query(sql)
	} else {
		return di.Exec(sql)
	}
}
func (di *DBConnectInfo) query(sql string) *ExecResult {
	result := &ExecResult{DataTable: *NewDataTable(), Sql: sql, IsQuery: true}
	dt, err := result.DataTable.Fill(di.db, sql)
	if err != nil {
		dt.Error = err.Error()
	}
	return result
}

// =============================================
// public methods
func (di *DBConnectInfo) ExecBatch(sqls []*SqlIndex, trans bool) *[]ExecResult {
	db, err := di.Open()
	results := make([]ExecResult, 0)
	if err != nil {
		result := &ExecResult{}
		result.Error = err.Error()
		results = append(results, *result)
		return &results
	}
	defer db.Close()

	//=============================================
	var tx *sql.Tx
	if trans {
		tx, err = db.BeginTx(context.Background(), nil)
		if err != nil {
			result := &ExecResult{}
			result.Error = err.Error()
			results = append(results, *result)
			return &results
		}
		defer tx.Rollback()
	}

	sort.SliceStable(sqls, func(i, j int) bool { return sqls[i].Index < sqls[j].Index })
	//=============================================
	for _, sql := range sqls {
		result := di.batch(sql.Sql)
		results = append(results, *result)
	}
	//=============================================
	if trans && tx != nil {
		tx.Commit()
	}
	return &results
}
