package controllers

import (
	"GanLianInfo/models"
	"io/ioutil"
	"net/http"
	"strconv"

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
	ID           int64    `json:"id,string"`
	IdCode       string   `json:"idCode"`
	Name         string   `json:"name"`
	PoliceCode   string   `json:"policeCode"`
	OrganID      string   `json:"organ"`
	DepartmentID string   `json:"departmentId"`
	Roles        []string `json:"roles" gorm:"-"`
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
		log.Errorf("从大数据中心获取数据发生错误:%v\n", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("读取大数据中心返回数据发生错误:%v\n", err)
	}
	// 解析大数据中心返回的数据
	//id := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "user").ToString()
	//id := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "attributes", "data", 4, "idCode").ToString()
	var id string
	temp := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "attributes", "data").ToString()
	var tmp []map[string]string
	err = jsoniter.UnmarshalFromString(temp, &tmp)
	if err != nil {
		log.Error(err)
	}
	// 遍历找到身份证号idCode
	for _, v := range tmp {
		if value, ok := v["idCode"]; ok {
			id = value
		}
	}
	log.Successf("idCode:%s\n", id)
	if id == "" {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	var p struct {
		ID int64
	}
	// 根据身份证号在数据库里找对应人员
	result := db.Model(&models.Personnel{}).Where("id_code = ?", id).Find(&p)
	if result.RowsAffected == 0 {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	// 把int64型的id转换成string返回前端
	r = gin.H{"code": 20000, "token": strconv.FormatInt(p.ID, 10), "isValid": true}
	c.JSON(200, r)
}

func UserInfo(c *gin.Context) {
	var id struct {
		ID int64 `json:"id,string"`
	}
	if err := c.ShouldBindJSON(&id); err != nil {
		r := Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}

	p := GetUserRoles(id.ID)
	r := gin.H{"code": 20000, "data": p}
	c.JSON(200, r)
}

func GetUserRoles(id int64) *UserRole {
	var p UserRole
	//p.IdCode = id
	//db.Model(&models.Personnel{}).Select("id", "name", "police_code", "organ_id", "department_id").First(&p, "id=?", id)
	db.Model(&models.Personnel{}).Select("id", "name", "id_code", "police_code", "organ_id", "department_id").Where("id = ?", id).Limit(1).Find(&p)
	roles := enforcer.GetFilteredGroupingPolicy(0, strconv.FormatInt(id, 10))
	if len(roles) == 0 {
		p.Roles = append(p.Roles, "normal")
		return &p
	}
	for _, v := range roles {
		p.Roles = append(p.Roles, v[1])
	}
	return &p
}

func GetPersonOrganId(c *gin.Context) {
	var id struct {
		ID int64 `json:"id,string"`
	}
	var r gin.H
	if c.ShouldBindJSON(&id) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	var p struct {
		OrganId string
	}

	if id.ID == 0 {
		r = gin.H{"code": 20000, "data": ""}
		c.JSON(200, r)
		return
	}
	db.Table("personnels").Select("organ_id").Where("id = ?", id.ID).Limit(1).Find(&p)
	r = gin.H{"code": 20000, "data": p.OrganId}
	c.JSON(200, r)
	return
}

// GetPersonOrgans 获取所有用户的id
func GetPersonOrgans(c *gin.Context) {
	var r gin.H
	var p []struct {
		ID      int64
		OrganId string
	}
	_map := make(map[string]string)
	db.Table("personnels").Select("id,organ_id").Find(&p)
	for _, v := range p {
		_map[strconv.FormatInt(v.ID, 10)] = v.OrganId
	}
	r = gin.H{"code": 20000, "data": _map}
	c.JSON(200, r)
	return
}

// SetPersonOrganMap 将用户id与organ_id的map写入redis
//func SetPersonOrganMap() {
//	var p []struct {
//		ID      string
//		OrganId string
//	}
//	_map := make(map[string]string)
//	db.Table("personnels").Select("id,organ_id").Find(&p)
//	for _, v := range p {
//		_map[v.ID] = v.OrganId
//
//	}
//	if rdb != nil {
//		res, _ := rdb.Exists(ctx, "personOrganMap").Result()
//		log.Infof("res: %d\n", res)
//		if res == 0 {
//			rdb.HSet(ctx, "personOrganMap", _map)
//		}
//
//	}
//}
