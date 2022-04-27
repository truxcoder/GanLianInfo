package controllers

type ErrorCode int

const (
	CodeUnique ErrorCode = iota + 1000
	CodePassword
	CodeInvalid
	CodeAdd
	CodeUpdate
	CodeDelete
	CodeNoData
	CodeServer
	CodeBind
	CodeDatabase
)

func GetError(code ErrorCode) map[string]interface{} {
	var message string
	switch code {
	case CodeUnique:
		message = "添加失败！请检查编号或其他唯一性字段是否重复"
		break
	case CodePassword:
		message = "帐号或密码错误"
		break
	case CodeInvalid:
		message = "无权限"
		break
	case CodeAdd:
		message = "添加失败"
		break
	case CodeUpdate:
		message = "更新失败"
		break
	case CodeDelete:
		message = "删除失败"
		break
	case CodeNoData:
		message = "未查询到数据"
		break
	case CodeServer:
		message = "服务器错误"
		break
	case CodeBind:
		message = "后端服务器执行数据绑定时发生错误, 请检查数据合法性"
		break
	case CodeDatabase:
		message = "后端数据库连接错误, 请通知管理员检查网络连接状况"
		break
	default:
		message = "后端服务器发生错误"
		break
	}
	return map[string]interface{}{
		"code":    int(code),
		"message": message,
	}
}

//type ErrorStruct struct {
//	Unique            map[string]interface{}
//	PasswordIncorrect map[string]interface{}
//	GetUser           map[string]interface{}
//	Insert            map[string]interface{}
//	Update            map[string]interface{}
//	Delete            map[string]interface{}
//	NoData            map[string]interface{}
//	ServerError       map[string]interface{}
//	DatabaseError     map[string]interface{}
//}
//
//var Errors = ErrorStruct{
//	Unique: map[string]interface{}{
//		"message": "添加失败！请检查编号或其他唯一性字段是否重复!",
//		"code":    80405,
//	},
//	PasswordIncorrect: map[string]interface{}{
//		"code":    60204,
//		"message": "帐号或密码错误！",
//	},
//	GetUser: map[string]interface{}{
//		"code":    50008,
//		"message": "登录失败，无法取得用户信息！",
//	},
//	Insert: map[string]interface{}{
//		"message": "添加失败！",
//		"code":    80404,
//	},
//	Update: map[string]interface{}{
//		"message": "更新失败！",
//		"code":    80404,
//	},
//	Delete: map[string]interface{}{
//		"message": "删除失败！",
//		"code":    80404,
//	},
//	NoData: map[string]interface{}{
//		"message": "未查询到数据！",
//		"code":    40400,
//		"count":   0,
//	},
//	ServerError: map[string]interface{}{
//		"message": "服务器错误！！",
//		"code":    80503,
//	},
//	DatabaseError: map[string]interface{}{
//		"message": "后端数据库连接错误, 请通知管理员检查网络连接状况!",
//		"code":    80503,
//	},
//}
