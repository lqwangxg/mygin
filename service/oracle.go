package service

// type OracleDBInfo DBInfo

// func GetOracleInfo(user, dbname string) *OracleDBInfo {
// 	return &OracleDBInfo{
// 		Host: "db.dev.mbpsmartec.co.jp",
// 		Port: 1521,
// 		User: user,
// 		Pass: user,
// 		Name: dbname,
// 	}
// }

// func (di *OracleDBInfo) Open() (*sql.DB, error) {
// 	dataSourceName := go_ora.BuildUrl(di.Host, di.Port, di.Name, di.User, di.Pass, nil)
// 	db, err := sql.Open("oracle", dataSourceName)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return db, nil
// }

// func (di *OracleDBInfo) Query(strSql string) interface{} {
// 	db, err := di.Open()
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer db.Close()

// 	rows, err := db.Query(strSql)
// 	if err != nil {
// 		fmt.Println("Error running query")
// 		fmt.Println(err)
// 		return nil
// 	}
// 	defer rows.Close()

// 	var thedate string
// 	for rows.Next() {
// 		rows.Scan(&thedate)
// 	}
// 	fmt.Printf("The date is: %s\n", thedate)
// 	return thedate
// }
const (
	SysDate_Oracle = "select sysdate from dual"
)

// func (di *OracleDBInfo) SysDate() interface{} {
// 	return di.Query("select sysdate from dual")
// }
