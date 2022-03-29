package controllers

import (
	"fmt"
	"net/http"

	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Infof("form:%v\n", file)
	log.Info(file.Filename)
	dst := "http://30.29.2.67/static/upload"

	// 上传文件至指定目录
	if err := c.SaveUploadedFile(file, dst); err != nil {
		log.Errorf("上传文件发生错误:%s\n", err.Error())
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
