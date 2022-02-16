package controllers

import (
	"GanLianInfo/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type InfoParam struct {
	Token string `form:"token"`
}

type UserRole struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PoliceCode   string   `json:"policeCode"`
	OrganID      string   `json:"organ"`
	DepartmentID string   `json:"departmentId"`
	Roles        []string `json:"roles"`
}

// Login 单点登录认证
func Login(c *gin.Context) {
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
	if id == "" {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	var count int64
	db.Model(&models.Personnel{}).Where("id = ?", id).Count(&count)
	if err != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if count == 0 {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "token": id, "isValid": true}
	c.JSON(200, r)
}

func Logout(c *gin.Context) {
	r := gin.H{"code": 20000, "data": "success"}
	c.JSON(200, r)
}

func UserInfo(c *gin.Context) {
	id := c.Query("id")
	p := GetUserRoles(id)
	r := gin.H{"code": 20000, "data": p}
	c.JSON(200, r)
}

//func GetUserInfo(c *gin.Context) {
//	u := GetUserRoles("3f024fb3c9494a7292b3f2368cd402a9")
//	r := gin.H{"code": 20000, "data": u}
//	c.JSON(200, r)
//}

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
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if count == 0 {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	p := GetUserRoles(id)
	r = gin.H{"code": 20000, "data": p, "isValid": true}
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

func GetUserRoles(id string) *UserRole {
	var p UserRole
	p.ID = id
	db.Model(&models.Personnel{}).Select("id", "name", "police_code", "organ_id", "department_id").First(&p, "id=?", id)
	roles := enforcer.GetFilteredGroupingPolicy(0, id)
	if len(roles) == 0 {
		p.Roles = append(p.Roles, "normal")
		return &p
	}
	for _, v := range roles {
		p.Roles = append(p.Roles, v[1])
	}
	return &p
}

//func GetUserList (c *gin.Context) {
//	query := c.Request.URL.Query()
//	params := map[string]interface{}{}
//	for key, value := range query {
//		params[key] = value[0]
//	}
//	fmt.Println("params:",params)
//	db := dao.Connect()
//	var m []models.User
//	var r gin.H
//		result := db.Preload("Department").Where(params).Find(&m)
//		err := result.Error
//		if err!=nil {
//			r = Errors.ServerError
//		} else {
//			r = gin.H{"code": 20000,"data": &m}
//		}
//
//	c.JSON(200, r)
//}
//func UserAdd(c *gin.Context) {
//	var user models.User
//	db := dao.Connect()
//	var r gin.H
//	if c.ShouldBindJSON(&user) != nil {
//		r = Errors.ServerError
//	} else {
//		user.CreateTime = time.Now()
//		db.Omit("Department").Create(&user)
//		r = gin.H{"message": "添加成功！", "code": 20000}
//	}
//	c.JSON(200, r)
//}
//func UserUpdate(c *gin.Context) {
//	var user models.User
//	db := dao.Connect()
//	var r gin.H
//	if c.ShouldBindJSON(&user) != nil {
//		r = Errors.ServerError
//	} else {
//		db.Omit("Department","CreateTime").Model(&user).Updates(&user)
//		r = gin.H{"message": "更新成功！", "code": 20000}
//	}
//	c.JSON(200, r)
//}
//
//func UserDelete(c *gin.Context) {
//	var id IdStruct
//	var r gin.H
//	db := dao.Connect()
//	if c.ShouldBindJSON(&id) != nil {
//		r = Errors.Delete
//	} else {
//		result := db.Where(&id.Id).Delete(&models.User{})
//		err := result.Error
//		if err != nil {
//			fmt.Println("err:",err)
//			r = Errors.ServerError
//		} else {
//			message := fmt.Sprintf("成功删除%d条数据",result.RowsAffected)
//			r = gin.H{"message": message, "code": 20000}
//		}
//	}
//	c.JSON(200, r)
//}
