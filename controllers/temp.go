package controllers

import (
	"GanLianInfo/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

func PoliceInfo(c *gin.Context) {
	ticket := c.Query("ticket")
	service := c.Query("service")
	var r gin.H
	//向大数据中心请求数据
	url := "http://30.29.2.6:8686/cas/serviceValidate"
	params := "ticket=" + ticket + "&service=" + service + "&format=json"
	uri := url + "?" + params
	resp, err := http.Get(uri)
	if err != nil {
		log.Errorf("get failed, err:%v\n", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get resp failed, err:%v\n", err)
	}
	//解析大数据中心返回的数据
	id := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "user").ToString()
	name := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "attributes", "data", 4, "realName").ToString()
	log.Infof("id:%s,name:%s\n", id, name)
	if id == "" {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	var count int64
	db.Model(&models.Personnel{}).Where("id = ?", id).Count(&count)
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	if count == 0 {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}

	c.JSON(200, r)
}

func PolicePhoto(c *gin.Context) {
	id := c.Query("id")
	b := _requestPhoto(id)
	//contentType := "image/jpeg"
	//c.Data(200, contentType, b)
	r := gin.H{"code": 20000, "data": string(b)}
	c.JSON(200, r)
}

func _requestPhoto(id string) []byte {
	url := "http://30.29.2.6:8686/unionapi/user/headshot/base64"
	// json数据
	//contentType := "application/json"
	//data := `{"name":"小王子","age":18}`
	// 表单数据
	contentType := "application/x-www-form-urlencoded"
	data := "Authorization=438019355f6940fba3b98316d97fd5f0" + "&userId=" + id
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return b
}
