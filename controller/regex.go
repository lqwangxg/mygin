package controller

import (
	"mygin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Match(c *gin.Context) {
	regex := &service.RegexText{
		Pattern: c.PostForm("pattern"),
		Content: c.PostForm("content"),
	}
	if regex.Pattern == "" && regex.Content == "" {
		Index(c)
		return
	}
	result := regex.GetMatchResult()
	c.HTML(http.StatusOK, "match.html", gin.H{
		"Pattern": regex.Pattern,
		"Content": regex.Content,
		"result":  result,
	})
}
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		//"user_agent": c.GetHeader("User-Agent"),
	})
}
