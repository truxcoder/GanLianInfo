package controllers

import (
	"GanLianInfo/models"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

type Account struct {
	ID           string    `json:"userId"`
	PersonnelId  int64     `json:"id,string"`
	Name         string    `json:"realName"`
	Username     string    `json:"userName"`
	IdCode       string    `json:"idCode"`
	OrganID      string    `json:"organID"`
	DepartmentID string    `json:"deptId"`
	UserType     int8      `json:"userType"`
	DataStatus   int8      `json:"dataStatus"`
	Sort         int       `json:"sort"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

func AccountSync(c *gin.Context) {
	data := GetPersonnelDataFromInterfaceForAccount()
	var p, added, updated, deleted PerSlice
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
	}
	sort.Sort(p)
	for i := 0; i < len(p); i++ {
		v := p[i]
		var id string
		var account models.Account
		if v.IdCode == "" {
			continue
		}
		result := db.Model(models.Account{}).Where("id", v.UserID).Limit(1).Find(&account)
		isFound := result.RowsAffected == 1
		isValid := (v.UserType == 1 || v.UserType == 12) && v.DataStatus == 0
		isUpdated := v.UpdateTime.After(account.UpdateTime)
		if !isFound && !isValid {
			continue
		}
		if isFound && !isUpdated {
			continue
		}
		id = getOrganIdFromDepartmentId(v.DepartmentID)
		v.OrganID = id
		// 过滤掉泸州所和攀枝花所
		if v.OrganID == "c84c0a0ae2e54c5baf8c9d8c86fc9761" || v.OrganID == "6a8f659d05a74ee582c4880083ed606d" {
			continue
		}
		v.ID = getIdFromIdCode(v.IdCode)

		if !isFound && isValid {
			added = append(added, v)
		} else if isFound && isValid && isUpdated {
			updated = append(updated, v)
		} else if isFound && isUpdated {
			// TODO: 这里的删除逻辑需要大数据中心开放身份证验证接口，否则无法实现
			//v.ID = account.ID
			deleted = append(deleted, v)
		}
	}
	if len(added) == 0 && len(updated) == 0 && len(deleted) == 0 && rdb != nil {
		//res, _ := rdb.Exists(ctx, "updateTime").Result()
		now := time.Now()
		rdb.Set(ctx, "accountUpdateTime", now, time.Hour*2400)

	}
	r := gin.H{"code": 20000, "add": &added, "update": &updated, "delete": &deleted}
	c.JSON(200, r)
}

func AccountSure(c *gin.Context) {
	var r gin.H
	method := c.Query("method")
	var a []Account
	//if method != "delete" {
	//	if err := c.ShouldBindJSON(&a); err != nil {
	//		log.Error(err)
	//		r = GetError(CodeBind)
	//		c.JSON(200, r)
	//		return
	//	}
	//}

	if err := c.ShouldBindJSON(&a); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	if method == "add" {
		if result := db.Table("accounts").Create(a); result.Error != nil {
			r = GetError(CodeAdd)
			c.JSON(200, r)
			return
		}
		r = gin.H{"code": 20000, "message": "添加成功!"}
		c.JSON(200, r)
		return
	}
	if method == "update" {
		for _, v := range a {
			//log.Successf("v:%+v\n", v)
			db.Table("accounts").Where("id = ?", v.ID).Updates(&v)
		}
		r = gin.H{"code": 20000, "message": "更新成功!"}
		c.JSON(200, r)
		return
	}
	if method == "delete" {
		//var id struct {
		//	Id []string `json:"id"`
		//}
		//if err := c.ShouldBindJSON(&id); err != nil {
		//	log.Error(err)
		//	r = GetError(CodeBind)
		//	c.JSON(200, r)
		//	return
		//}

		var ids []string
		for _, v := range a {
			ids = append(ids, v.ID)
		}
		//db.Where("id in ?", &ids).Delete(models.Account{})
		db.Delete(models.Account{}, &ids)
		r = gin.H{"code": 20000, "message": "删除成功!"}
		c.JSON(200, r)
		return
	}
}
