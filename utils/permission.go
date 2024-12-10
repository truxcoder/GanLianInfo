package utils

import (
	"github.com/casbin/casbin/v2"
	"strings"
)

func HasBothPermission(enforcer *casbin.Enforcer, sub string, obj string, act string) (has bool, err error) {
	var (
		roles []string
		acts  []string
	)
	if roles, err = enforcer.GetRolesForUser(sub); err != nil || len(roles) == 0 {
		return
	}
	acts = strings.Split(act, "_")
	if len(acts) == 0 {
		return
	}
	for _, role := range roles {
		has = true
		for _, a := range acts {
			if !enforcer.HasPermissionForUser(role, obj, a) {
				has = false
				break
			}
		}
		if has {
			return
		}
	}
	return
}
