package controllers

import (
	"GanLianInfo/auth"
	"GanLianInfo/models"
	"io/ioutil"
	"net/http"

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
	PersonnelId  string   `json:"personnelId"`
	IdCode       string   `json:"idCode"`
	Name         string   `json:"name"`
	Username     string   `json:"username"`
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
	id := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "user").ToString()
	//id := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "attributes", "data", 4, "idCode").ToString()
	//var id string
	//temp := jsoniter.Get(b, "serviceResponse", "authenticationSuccess", "attributes", "data").ToString()
	//var tmp []map[string]string
	//err = jsoniter.UnmarshalFromString(temp, &tmp)
	//if err != nil {
	//	log.Error(err)
	//}
	// 遍历找到身份证号idCode
	//for _, v := range tmp {
	//	if value, ok := v["idCode"]; ok {
	//		id = value
	//	}
	//}
	//log.Successf("idCode:%s\n", id)
	if id == "" {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}
	var u struct {
		ID string
	}
	// 根据身份证号在数据库里找对应人员
	//result := db.Model(&models.Personnel{}).Where("status = 1 AND id_code = ?", id).Find(&p)
	//if result.RowsAffected == 0 {
	//	r = gin.H{"code": 20000, "isValid": false}
	//	c.JSON(200, r)
	//	return
	//}

	result := db.Model(&models.Account{}).Where("id = ?", id).Find(&u)
	if result.RowsAffected == 0 {
		r = gin.H{"code": 20000, "isValid": false}
		c.JSON(200, r)
		return
	}

	// 把jwt生成的token返回前端
	tokenString, _ := auth.GenToken(id)
	data := gin.H{"id": u.ID, "token": tokenString}
	r = gin.H{"code": 20000, "data": data, "isValid": true}
	c.JSON(200, r)
}

func UserInfo(c *gin.Context) {
	var id struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&id); err != nil {
		r := GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	p := GetUserRoles(id.ID)
	r := gin.H{"code": 20000, "data": p}
	c.JSON(200, r)
}

func GetUserRoles(id string) *UserRole {
	var u UserRole
	//p.IdCode = id
	//db.Model(&models.Personnel{}).Select("id", "name", "police_code", "organ_id", "department_id").First(&p, "id=?", id)
	db.Model(&models.Account{}).Select("id", "personnel_id", "name", "id_code", "username", "organ_id", "department_id").Where("id = ?", id).Limit(1).Find(&u)
	roles := enforcer.GetFilteredGroupingPolicy(0, id)
	if len(roles) == 0 {
		u.Roles = append(u.Roles, "normal")
		return &u
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, v[1])
	}
	return &u
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
