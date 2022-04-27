package controllers

import (
	"GanLianInfo/models"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type Role struct {
	AccountId string `json:"accountId"`
	Role      string `json:"role"`
}

type AccountRole struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	OrganID  string `json:"organId"`
	Roles    string `json:"role"`
}

func RoleList(c *gin.Context) {
	var (
		err   error
		r     gin.H
		roles []AccountRole
		where string
		args  []interface{}
	)
	var mo struct {
		AccountId string `json:"accountId"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		OrganId   string `json:"organId"`
	}
	var selectStr = "casbin_rule.v0 as id, casbin_rule.v1 as roles, accounts.name, accounts.username, accounts.organ_id"
	var joinStr = "left join accounts on accounts.id = casbin_rule.v0"
	if err = c.BindJSON(&mo); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	where = "1=1"
	if mo.AccountId != "" {
		where += " and v0 = ? "
		args = append(args, mo.AccountId)
	}
	if mo.Username != "" {
		var users []string
		var username string
		db.Model(&models.Account{}).Where("username = ?", mo.Username).Pluck("id", &users)
		if len(users) > 0 {
			username = users[0]
		}
		where += " and v0 = ? "
		args = append(args, username)
	}
	if mo.Role != "" {
		where += " and v1= ? "
		args = append(args, mo.Role)
	}
	if mo.OrganId != "" {
		var users []string
		db.Model(&models.Account{}).Where("organ_id = ?", mo.OrganId).Pluck("id", &users)
		where += " and v0 in ? "
		args = append(args, users)
	}
	db.Table("casbin_rule").Select(selectStr).Joins(joinStr).Where("ptype = ?", "g").Where(where, args...).Order("username asc").Find(&roles)
	//if mo.Role != "" {
	//	result = enforcer.GetFilteredGroupingPolicy(1, mo.Role)
	//} else {
	//	result = enforcer.GetGroupingPolicy()
	//}
	//// TODO: 等帐号系统建立后增加角色查询功能
	//result = enforcer.GetGroupingPolicy()
	//for _, v := range result {
	//	var perRole PerRole
	//	_v, _ := strconv.Atoi(v[0])
	//	perRole.Id = int64(_v)
	//	db.Model(&models.Personnel{}).Select("id", "name", "police_code", "organ_id").Where(&perRole).First(&perRole)
	//	perRole.Role = v[1]
	//	roles = append(roles, perRole)
	//}
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
		//r = Errors.ServerError
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}

	added, err = enforcer.AddGroupingPolicy(role.AccountId, role.Role)
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
		//r = Errors.ServerError
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	success, err := enforcer.UpdateGroupingPolicy(role.Old, role.New)
	if err != nil {
		log.Error(err)
		//r = Errors.ServerError
		r = GetError(CodeServer)
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
		//r = Errors.ServerError
		r = GetError(CodeBind)
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
		//r = Errors.ServerError
		r = GetError(CodeServer)
	} else {
		r = gin.H{"code": 20000, "data": &rd}
	}
	c.JSON(200, r)
}

func RoleDictAdd(c *gin.Context) {
	var rd models.RoleDict
	var r gin.H
	if c.ShouldBindJSON(&rd) != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
	} else {
		db.Create(&rd)
		r = gin.H{"message": "添加成功！", "code": 20000}
	}
	c.JSON(200, r)
}

// RoleDictUpdate 修改角色
func RoleDictUpdate(c *gin.Context) {
	var rd, one models.RoleDict
	var r gin.H
	var newPolicies [][]string
	if c.ShouldBindJSON(&rd) != nil {
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	db.First(&one, rd.ID)
	if one.Name != rd.Name {
		// 验证是角色是否分配了权限, 如分配权限，修改对应的policy
		policy := enforcer.GetFilteredPolicy(0, one.Name)
		if len(policy) > 0 {
			for _, v := range policy {
				newPolicies = append(newPolicies, []string{rd.Name, v[1], v[2]})
			}
			ok, err := enforcer.UpdatePolicies(policy, newPolicies)
			if err != nil {
				r = gin.H{"code": 1403, "message": err.Error()}
				c.JSON(200, r)
				return
			}
			if !ok {
				r = gin.H{"code": 1403, "message": "修改失败!"}
				c.JSON(200, r)
				return
			}
		}

		// 验证是否将用户分配给该角色，如分配，则修改为新角色名
		user, _ := enforcer.GetUsersForRole(one.Name)
		if len(user) > 0 {
			var oldRules, newRules [][]string
			for _, v := range user {
				oldRules = append(oldRules, []string{v, one.Name})
				newRules = append(newRules, []string{v, rd.Name})
			}
			ok, err := enforcer.UpdateGroupingPolicies(oldRules, newRules)
			if err != nil {
				r = gin.H{"code": 1403, "message": err.Error()}
				c.JSON(200, r)
				return
			}
			if !ok {
				r = gin.H{"code": 1403, "message": "修改失败!"}
				c.JSON(200, r)
				return
			}
		}
	}
	db.Model(&rd).Updates(&rd)
	r = gin.H{"message": "更新成功！", "code": 20000}
	c.JSON(200, r)
}

func RoleDictDelete(c *gin.Context) {
	var dict models.RoleDict
	var r gin.H
	var err error
	if err = c.ShouldBindJSON(&dict); err != nil {
		r = GetError(CodeBind)
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
		//r = Errors.ServerError
		r = GetError(CodeServer)
	} else {
		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
		r = gin.H{"message": message, "code": 20000}
	}
	c.JSON(200, r)
	return
}
