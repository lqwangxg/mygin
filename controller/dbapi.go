package controller

import (
	"mygin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Show(c *gin.Context) {

	services := make([]service.ApiInfo, 0)
	var desc *service.ApiInfo
	var params []string
	//==================================================
	params = make([]string, 0)
	params = append(params, "host: host name or ip address")
	params = append(params, "port: port of service")
	params = append(params, "dbtype: type of databases. oracle, sqlserver, postgre")
	desc = service.NewApiDescInfo("show", params, "method: get.")
	services = append(services, *desc)
	//==================================================
	params = make([]string, 0)
	params = append(params, "sql: standard query language for relation database like oracle, sqlserver, postgres.")
	desc = service.NewApiDescInfo("query", params, "query resultSet. method: post.")
	services = append(services, *desc)
	//==================================================
	params = make([]string, 0)
	params = append(params, "sql: standard query language for relation database like oracle, sqlserver, postgres.")
	desc = service.NewApiDescInfo("update", params, "count of insert/update/delete. method: post.")
	services = append(services, *desc)
	//==================================================

	c.IndentedJSON(http.StatusOK, services)
}

func SysDate(c *gin.Context) {
	var request service.DBRequest
	c.Bind(&request)
	db := service.NewConnection(&request)
	dt, err := db.QueryTable(service.SysDate_Oracle)
	if err == nil {
		sysdate := dt.Rows[0][0]
		c.IndentedJSON(http.StatusOK, sysdate)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	//==================================================

}

func Exec(c *gin.Context) {
	var request service.DBRequest
	c.Bind(&request)
	paramPairs := c.Request.URL.Query()
	for key, values := range paramPairs {
		if service.NewRegexText(`sql\w+`, key).IsMatch() {
			request.Sqls = append(request.Sqls, values[0])
		}
	}
	result := request.Exec()
	c.IndentedJSON(http.StatusOK, result)
}

func QueryValue(c *gin.Context) {
	var request service.DBRequest
	c.Bind(&request)
	db := service.NewConnection(&request)
	dt, err := db.QueryTable(request.Sql)
	if err == nil {
		value := dt.Rows[0][0]
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	//==================================================

}
func QueryTable(c *gin.Context) {
	var request service.DBRequest
	c.Bind(&request)
	db := service.NewConnection(&request)
	dt, err := db.QueryTable(request.Sql)
	if err == nil {
		c.IndentedJSON(http.StatusOK, dt)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	//==================================================

}
