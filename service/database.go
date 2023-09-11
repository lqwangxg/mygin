package service

import (
	"database/sql"
	"fmt"

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
func (di *DBConnectInfo) Exec(strSQL string) (int64, error) {
	//====================================
	var count int64
	if res, err := di.db.Exec(strSQL); err == nil {
		if i, err := res.LastInsertId(); err == nil {
			count = count + i
		}
		if i, err := res.RowsAffected(); err == nil {
			count = count + i
		}
		return count, nil
	} else {
		return count, err
	}
}
