package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type RBAC struct {
	Sub string   `json:"sub"`
	Obj string   `json:"obj"`
	Act []string `json:"act"`
}

type PermissionModel struct {
	Sub string `json:"sub"`
	Obj string `json:"obj"`
	Act string `json:"act"`
}

type RolePermission struct {
	ModuleName     string `json:"moduleName"`
	ModuleTitle    string `json:"moduleTitle"`
	RoleName       string `json:"roleName"`
	RoleTitle      string `json:"roleTitle"`
	PermissionName string `json:"permissionName"`
}

func PermissionCheck(c *gin.Context) {
	var rbac RBAC
	var r gin.H
	var err error
	var result = make(map[string]bool)
	if err = c.BindJSON(&rbac); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if len(rbac.Act) < 1 {
		log.Error(err)
		r = Errors.NoData
		c.JSON(200, r)
		return
	}
	for _, v := range rbac.Act {
		ok, _ := enforcer.Enforce(rbac.Sub, rbac.Obj, v)
		result[v] = ok
	}
	r = gin.H{"code": 20000, "data": &result}
	c.JSON(200, r)
	return
}

func PermissionActCheck(c *gin.Context) {
	var rbac struct {
		Sub string   `json:"sub"`
		Obj []string `json:"obj"`
		Act string   `json:"act"`
	}
	var r gin.H
	var err error
	var result = make(map[string]bool)
	if err = c.BindJSON(&rbac); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if len(rbac.Obj) < 1 {
		log.Error(err)
		r = Errors.NoData
		c.JSON(200, r)
		return
	}
	for _, v := range rbac.Obj {
		ok, _ := enforcer.Enforce(rbac.Sub, v, rbac.Act)
		result[v] = ok
	}
	r = gin.H{"code": 20000, "data": &result}
	c.JSON(200, r)
	return
}

func GetRolePermission(c *gin.Context) {
	var err error
	var r gin.H
	var role struct {
		Role string `json:"role"`
	}
	if err = c.BindJSON(&role); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	result := enforcer.GetFilteredPolicy(0, role.Role)
	r = gin.H{"code": 20000, "data": &result}
	c.JSON(200, r)

}

func PermissionList(c *gin.Context) {
	//var err error
	//var r gin.H
	//var result [][]string
	//var id struct {
	//	ID string `json:"id"`
	//}
	//if err = c.BindJSON(&id); err != nil {
	//	log.Error(err)
	//	r = Errors.ServerError
	//	c.JSON(200, r)
	//	return
	//}
	//
	//if id.ID != "" {
	//	result = enforcer.GetFilteredPolicy(0, id.ID)
	//	roles := enforcer.GetFilteredGroupingPolicy(0, id.ID)
	//	for _, v := range roles {
	//		temp := enforcer.GetFilteredPolicy(0, v[1])
	//		result = append(result, temp...)
	//	}
	//} else {
	//	result = enforcer.GetPolicy()
	//}
	//r = gin.H{"code": 20000, "data": &result}
	//c.JSON(200, r)
	//return

	//var err error
	var r gin.H
	var rp []RolePermission

	selectStr := "casbin_rule.v2 as permission_name," +
		"modules.title as module_title, modules.name as module_name," +
		"role_dicts.name as role_name,role_dicts.title as role_title "
	joinStr := "left join modules on casbin_rule.v1 = modules.name " +
		"left join role_dicts on casbin_rule.v0 = role_dicts.name"
	db.Debug().Table("casbin_rule").Select(selectStr).Joins(joinStr).Where("casbin_rule.ptype = ?", "p").Find(&rp)

	r = gin.H{"code": 20000, "data": &rp}
	c.JSON(200, r)
	return
}

func PermissionManage(c *gin.Context) {
	var err error
	var r gin.H
	var addMessage, delMessage string
	added := false
	deled := false
	var data = make(map[string][][]string)
	data["add"] = [][]string{}
	data["del"] = [][]string{}
	isAdd := c.Query("add")
	isDel := c.Query("del")

	if err = c.BindJSON(&data); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if isAdd == "true" {
		added, err = enforcer.AddPolicies(data["add"])
		if err != nil {
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
	}
	if isDel == "true" {
		deled, err = enforcer.RemovePolicies(data["del"])
		if err != nil {
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
	}
	if added {
		addMessage = "成功添加权限！"
	} else {
		addMessage = "未能添加权限！"
	}
	if deled {
		delMessage = "成功删除权限！"
	} else {
		delMessage = "未能删除权限！"
	}
	if isAdd == "true" && isDel == "true" {
		r = gin.H{"code": 20000, "message": addMessage + delMessage}
	} else if isAdd == "true" {
		r = gin.H{"code": 20000, "message": addMessage}
	} else if isDel == "true" {
		r = gin.H{"code": 20000, "message": delMessage}
	} else {
		r = gin.H{"code": 20000, "message": "无任何操作！"}
	}
	c.JSON(200, r)
	return
}

func PermissionAdd(c *gin.Context) {
	var err error
	var r gin.H
	var data [][]string

	if err = c.BindJSON(&data); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	added, err := enforcer.AddPolicies(data)
	if err != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if added {
		r = gin.H{"code": 20000, "message": "成功添加权限！"}
	} else {
		r = gin.H{"code": 20000, "message": "未能添加权限！"}
	}
	c.JSON(200, r)
	return
}

func GetPolicy(c *gin.Context) {
	p := enforcer.GetPolicy()
	r := gin.H{"code": 20000, "data": p}
	c.JSON(200, r)
	return
}

func PermissionDelete(c *gin.Context) {
	var data [][]string
	var r gin.H
	if err := c.BindJSON(&data); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}

	ok, err := enforcer.RemovePolicies(data)
	if err != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if ok {
		r = gin.H{"code": 20000, "message": "成功删除所选权限！"}
	} else {
		r = gin.H{"code": 20000, "message": "未能成功删除所选权限！"}
	}
	c.JSON(200, r)
	return

}
