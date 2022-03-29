package controllers

import (
	"GanLianInfo/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type Role struct {
	PersonnelId string `json:"personnelId"`
	Role        string `json:"role"`
}

type PerRole struct {
	Id         int64  `json:"id,string"`
	Name       string `json:"name"`
	PoliceCode string `json:"policeCode"`
	OrganID    string `json:"organId"`
	Role       string `json:"role"`
}

func RoleList(c *gin.Context) {
	//var role Role
	//var err error
	var r gin.H
	var roles []PerRole
	result := enforcer.GetGroupingPolicy()
	for _, v := range result {
		var perRole PerRole
		_v, _ := strconv.Atoi(v[0])
		perRole.Id = int64(_v)
		db.Model(&models.Personnel{}).Select("id", "name", "police_code", "organ_id").Where(&perRole).First(&perRole)
		perRole.Role = v[1]
		roles = append(roles, perRole)
	}
	r = gin.H{"code": 20000, "data": &roles}
	c.JSON(200, r)
	return
}

func RoleAdd(c *gin.Context) {
	var role Role
	var err error
	var r gin.H
	var added bool
	if err = c.BindJSON(&role); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	added, err = enforcer.AddGroupingPolicy(role.PersonnelId, role.Role)
	if err != nil {
		log.Error(err)
		r = gin.H{"code": 500, "message": err}
		c.JSON(200, r)
		return
	}
	if !added {
		r = gin.H{"code": 500, "message": "添加失败！"}
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "message": "添加成功！"}
	c.JSON(200, r)
	return
}

func RoleUpdate(c *gin.Context) {
	var role struct {
		Old []string `json:"old"`
		New []string `json:"new"`
	}
	var r gin.H
	var err error
	if err = c.BindJSON(&role); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	success, err := enforcer.UpdateGroupingPolicy(role.Old, role.New)
	if err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if !success {
		r = gin.H{"code": 20000, "message": "更新失败！"}
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "message": "更新成功！"}
	c.JSON(200, r)
	return
}

func RoleDelete(c *gin.Context) {
	//var id IdStruct
	var group []struct {
		ID   string `json:"id"`
		Role string `json:"role"`
	}
	var failedList []string
	var r gin.H
	if err := c.ShouldBindJSON(&group); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}

	for _, v := range group {
		//removed, _ := enforcer.RemoveFilteredGroupingPolicy(0, v.ID)
		removed, _ := enforcer.RemoveGroupingPolicy(v.ID, v.Role)
		if !removed {
			failedList = append(failedList, v.ID)
		}
	}

	if len(failedList) > 0 {
		message := fmt.Sprintf("成功删除%d条数据,%d条数据删除失败！", len(group)-len(failedList), len(failedList))
		r = gin.H{"message": message, "code": 20000, "failedList": failedList}
		c.JSON(200, r)
		return
	}

	message := fmt.Sprintf("成功删除%d条数据！", len(group))
	r = gin.H{"message": message, "code": 20000}
	c.JSON(200, r)
	return

}

func RoleDictList(c *gin.Context) {
	var r gin.H
	var rd []models.RoleDict
	result := db.Select("id", "name", "title").Find(&rd)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &rd}
	}
	c.JSON(200, r)
}

func RoleDictAdd(c *gin.Context) {
	var rd models.RoleDict
	var r gin.H
	if c.ShouldBindJSON(&rd) != nil {
		r = Errors.ServerError
	} else {
		db.Create(&rd)
		r = gin.H{"message": "添加成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func RoleDictUpdate(c *gin.Context) {
	var rd models.RoleDict
	var r gin.H
	if c.ShouldBindJSON(&rd) != nil {
		r = Errors.ServerError
	} else {
		db.Model(&rd).Updates(&rd)
		r = gin.H{"message": "更新成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func RoleDictDelete(c *gin.Context) {
	var dict models.RoleDict
	var r gin.H
	var err error
	if err = c.ShouldBindJSON(&dict); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	// 验证是角色是否分配了权限
	policy := enforcer.GetFilteredPolicy(0, dict.Name)
	if len(policy) > 0 {
		r = map[string]interface{}{
			"message": "失败! 该角色分配了权限，请先解除权限!",
			"code":    80405,
		}
		c.JSON(200, r)
		return
	}
	// 验证是否将用户分配给该角色
	user, _ := enforcer.GetUsersForRole(dict.Name)
	if len(user) > 0 {
		r = map[string]interface{}{
			"message": "失败! 有人员属于这个角色，请将相关人员取消该角色绑定!",
			"code":    80405,
		}
		c.JSON(200, r)
		return
	}

	result := db.Delete(&models.RoleDict{}, &dict.ID)
	err = result.Error
	if err != nil {
		log.Error(err)
		r = Errors.ServerError
	} else {
		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
		r = gin.H{"message": message, "code": 20000}
	}
	c.JSON(200, r)
	return
}
