package service

type DBRequest struct {
	Driver string   `form:"driver"`
	Host   string   `form:"host"`
	Port   int      `form:"port"`
	Name   string   `form:"name"`
	User   string   `form:"user"`
	Pass   string   `form:"pass"`
	Schema string   `form:"schema"`
	Sql    string   `form:"sql"`
	Sqls   []string `form:"sqls[]"`
	Trans  bool     `form:"trans"`
}
type SqlIndex struct {
	Index int
	Sql   string
}

// // ===============================================
// // ===============================================
// // public methods
// func (req *DBRequest) Exec() *DataSet {

// 	ds := NewDataSet()
// 	ds.Exec(req)
// 	return ds
// }
