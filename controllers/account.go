package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountDept struct {
	models.Account
	DepartmentName      string `json:"departmentName"`
	DepartmentShortName string `json:"departmentShortName"`
	OrganName           string `json:"organName"`
	OrganShortName      string `json:"organShortName"`
}

func AccountList(c *gin.Context) {
	var mos []AccountDept
	var mo models.Account
	selectStr := "accounts.*,d.name as department_name,d.short_name as department_short_name," +
		"o.name as organ_name,o.short_name as organ_short_name"
	joinStr := "left join departments as d on accounts.department_id = d.id " +
		"left join departments as o on accounts.organ_id = o.id"
	getList(c, "accounts", &mo, &mos, &selectStr, &joinStr)
}

func AccountBaseList(c *gin.Context) {
	var (
		err    error
		r      gin.H
		result *gorm.DB
	)

	var mos []struct {
		ID             string `json:"id" gorm:"autoIncrement:false;primaryKey"`
		PersonnelId    int64  `json:"personnelId,string"`
		Name           string `json:"name" gorm:"size:50"`
		Username       string `json:"username"`
		OrganShortName string `json:"organShortName"`
	}
	var mo struct {
		AccountId string `json:"accountId"`
		OrganId   string `json:"organId"`
	}
	selectStr := "accounts.id,accounts.personnel_id,accounts.name,accounts.username,o.short_name as organ_short_name"
	joinStr := "left join departments as o on accounts.organ_id = o.id"
	if c.ShouldBindJSON(&mo) != nil {
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	canGlobal, _ := enforcer.Enforce(mo.AccountId, "Personnel", "GLOBAL")
	if canGlobal {
		result = db.Table("accounts").Select(selectStr).Joins(joinStr).Find(&mos)
	} else {
		result = db.Table("accounts").Select(selectStr).Joins(joinStr).Where("organ_id = ?", mo.OrganId).Find(&mos)
	}
	err = result.Error
	if err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	if result.RowsAffected == 0 {
		//r = Errors.NoData
		r = GetError(CodeNoData)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &mos}
	c.JSON(200, r)
}
