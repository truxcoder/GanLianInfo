package controllers

type ResponseCode int

const (
	ResNoData ResponseCode = iota
	ResAddSuccess
)

func GetResponse(code ResponseCode) map[string]interface{} {
	var message string
	switch code {
	case ResNoData:
		message = "未查询到数据"
	case ResAddSuccess:
		message = "添加成功"
	default:
		message = "操作成功"
	}
	return map[string]interface{}{
		"code":    20000,
		"message": message,
	}
}
