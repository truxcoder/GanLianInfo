package controllers

//type Type int
//const (
//	UNIQUEERROR Type = iota
//	PASSWORDERROR
//	GETUSERERROR
//	INSERTERROR
//	UPDATEERROR
//	DELETEERROR
//	NODATAERROR
//	SERVERERROR
//)
//
//type messageMap map[string]string
//var message map[Type]messageMap
//
//func init(){
//	message[UNIQUEERROR] = map[string]string{
//		"code": "80405",
//		"message": "添加失败！请检查编号或其他唯一性字段是否重复!",
//	}
//	message[PASSWORDERROR] = map[string]string{
//		"code": "60204",
//		"message": "帐号或密码错误！",
//	}
//	message[GETUSERERROR] = map[string]string{
//		"code": "50008",
//		"message": "登录失败，无法取得用户信息！",
//	}
//	message[INSERTERROR] = map[string]string{
//		"code": "80404",
//		"message": "添加失败！",
//	}
//	message[UPDATEERROR] = map[string]string{
//		"code": "80404",
//		"message": "更新失败！",
//	}
//	message[DELETEERROR] = map[string]string{
//		"code": "80404",
//		"message": "删除失败！",
//	}
//	message[NODATAERROR] = map[string]string{
//		"code": "80404",
//		"message": "删除失败！",
//	}
//	message[SERVERERROR] = map[string]string{
//		"code": "80503",
//		"message": "服务器错误！",
//	}
//}

type ErrorStruct struct {
	Unique            map[string]interface{}
	PasswordIncorrect map[string]interface{}
	GetUser           map[string]interface{}
	Insert            map[string]interface{}
	Update            map[string]interface{}
	Delete            map[string]interface{}
	NoData            map[string]interface{}
	ServerError       map[string]interface{}
}

var Errors = ErrorStruct{
	Unique: map[string]interface{}{
		"message": "添加失败！请检查编号或其他唯一性字段是否重复!",
		"code":    80405,
	},
	PasswordIncorrect: map[string]interface{}{
		"code":    60204,
		"message": "帐号或密码错误！",
	},
	GetUser: map[string]interface{}{
		"code":    50008,
		"message": "登录失败，无法取得用户信息！",
	},
	Insert: map[string]interface{}{
		"message": "添加失败！",
		"code":    80404,
	},
	Update: map[string]interface{}{
		"message": "更新失败！",
		"code":    80404,
	},
	Delete: map[string]interface{}{
		"message": "删除失败！",
		"code":    80404,
	},
	NoData: map[string]interface{}{
		"message": "未查询到数据！",
		"code":    40400,
		"count":   0,
	},
	ServerError: map[string]interface{}{
		"message": "服务器错误！！",
		"code":    80503,
	},
}

