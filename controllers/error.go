package controllers

import "github.com/pkg/errors"

type ErrorCode int

var (
	RedisKeyNotFoundERR = errors.New("redis key not found")
	RedisNullERR        = errors.New("redis instance not found")
)

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
	CodeValidate
	CodeExist
	CodeParse
	CodeDataWrite
)

func GetError(code ErrorCode) map[string]interface{} {
	var message string
	switch code {
	case CodeUnique:
		message = "添加失败！请检查编号或其他唯一性字段是否重复"
	case CodePassword:
		message = "帐号或密码错误"
	case CodeInvalid:
		message = "无权限"
	case CodeAdd:
		message = "添加失败"
	case CodeUpdate:
		message = "更新失败"
	case CodeDelete:
		message = "删除失败"
	case CodeNoData:
		message = "未查询到数据"
	case CodeServer:
		message = "服务器错误"
	case CodeBind:
		message = "后端服务器执行数据绑定时发生错误, 请检查数据合法性"
	case CodeDatabase:
		message = "后端数据库连接错误, 请通知管理员检查网络连接状况"
	case CodeValidate:
		message = "数据验证发生错误"
	case CodeExist:
		message = "查询到数据库中已存在该条信息，操作失败"
	case CodeParse:
		message = "数据解析错误"
	case CodeDataWrite:
		message = "写入数据时发生错误!"
	default:
		message = "后端服务器发生错误"
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
