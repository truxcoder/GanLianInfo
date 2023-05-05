package controllers

import (
	"GanLianInfo/dao"
	"bytes"
	"context"
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	log "github.com/truxcoder/truxlog"

	"gorm.io/gorm"
)

const (
	accountOrder = `(case when length(d.level_code)>=3 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,3)) else null end) desc,
(case when length(d.level_code)>=6 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,6)) else null end) desc,
(case when length(d.level_code)>=9 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,9)) else null end) desc,
(case when length(d.level_code)>=12 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,12)) else null end) desc,
(case when length(d.level_code)>=15 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,15)) else null end) desc, 
accounts.sort desc nulls first`
)

var (
	db        *gorm.DB
	rdb       *redis.Client
	ctx       = context.Background()
	enforcer  *casbin.Enforcer
	reConnect bool
)

// 列表页排序字典
var orderMap = map[string]string{
	"appraisals": "years,season",
	"accounts":   accountOrder,
	"leaders":    "organ_id,sort",
}

var detailOrderMap = map[string]string{
	"appraisals": "years,season",
	"posts":      "start_day desc",
}

func init() {
	db = dao.Connect()
	FixDB()
	casbinInit()
	redisInit()
	//SetPersonOrganMap()
}

type IdStruct struct {
	Id []int64 `json:"id"`
}

type ID struct {
	ID int64 `json:"id,string"`
}

func (i *IdStruct) UnmarshalJSON(data []byte) error {
	var id struct {
		Id []string `json:"id"`
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	for _, v := range id.Id {
		_v, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		i.Id = append(i.Id, int64(_v))
	}
	return nil
}

func redisInit() {
	var err error
	if rdb, err = dao.RedisConnect(); err != nil {
		log.Error(err)
	}
}

func casbinInit() {
	var err error
	var a *gormadapter.Adapter
	gormadapter.TurnOffAutoMigrate(db)
	a, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Error(err)
	}
	dir, _ := os.Getwd()
	//path := filepath.Dir(dir)
	//log.Infof("path:%s\n", dir)
	model := filepath.Join(dir, "model.conf")

	enforcer, err = casbin.NewEnforcer(model, a)
	if err != nil {
		log.Error(err)
	}
}

func GetDb() *gorm.DB {
	return db
}

// FixDB 断线重连
func FixDB() bool {
	if reConnect {
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		reConnect = true
		log.Errorf("数据库Ping错误,  正在重连...\n")
		for {
			db = dao.Connect()
			sqlDB, err = db.DB()
			_err := sqlDB.Ping()
			if _err == nil {
				reConnect = false
				// 断线修复后初始化casbin模型
				casbinInit()
				log.Successf("数据库恢复, 时间为: %s\n", time.Now().Format(timeFormat))
				break
			}
			time.Sleep(5 * time.Minute)
		}
	}
	return true
}

func getPageData(c *gin.Context) (size int, offset int) {
	currenPage := c.Query("currentPage")
	pageSize := c.Query("pageSize")
	page, _ := strconv.Atoi(currenPage)
	size, _ = strconv.Atoi(pageSize)
	offset = (page - 1) * size
	return
}

// buildWhere 根据查询参数构建where语句。
//
// 可以查询organId，organIds，category和零值。零值查询需要前端在查询参数里传入一个key为zero，值为field_name$value_type的参数。
// 用$符号隔开。field_name:数据库字段名。value_type:数据类型，如time,int,string
func buildWhere(c *gin.Context) (string, []interface{}) {
	var result bytes.Buffer
	var params []interface{}
	organId := c.Query("organId")
	organIds := c.QueryArray("organId[]")
	category := c.Param("category")
	zero := c.Query("zero")
	result.WriteString(" 1=1 ")
	if organId != "" {
		result.WriteString(" and personnel_id in (select id from personnels where organ_id = ?) ")
		params = append(params, organId)
	}
	if len(organIds) > 0 {
		result.WriteString(" and personnel_id in (select id from personnels where organ_id in ?) ")
		params = append(params, organIds)
	}
	if category != "" {
		result.WriteString(" and category = ? ")
		params = append(params, category)
	}
	if zero != "" {
		pair := strings.Split(zero, "$")
		if len(pair) == 2 {
			result.WriteString(" and " + pair[0] + " = ? ")
			switch pair[1] {
			case "time":
				var t time.Time
				params = append(params, t)
			case "int":
				params = append(params, 0)
			case "string":
				params = append(params, "")
			}
		}
	}
	return result.String(), params
}

func getList(c *gin.Context, table string, mo interface{}, mos interface{}, selectStr *string, joinStr *string) {
	var (
		where  string
		params []interface{}
		r      gin.H
		err    error
		count  int64 //总记录数
		_map   map[string]interface{}
		_sql   string
	)

	queryMeans := c.Query("queryMeans") //请求方式，是前端分页还是后端分页
	sort, ok := orderMap[table]
	if !ok {
		sort = "id desc"
	}
	if err = c.BindJSON(mo); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	// 处理一个字段多值的情况。转成map，gorm会自动判断是否用"IN"查询语句
	_map = structToSlice(mo)
	_sql = intercept(mo)
	where, params = buildWhere(c)

	//后端分页情况
	if queryMeans == "backend" {
		size, offset := getPageData(c)
		//先查询数据总量并返回到前端
		if err = db.Table(table).Where(mo).Where(where, params...).Where(_map).Where(_sql).Count(&count).Error; err != nil {
			//r = Errors.ServerError
			r = GetError(CodeServer)
			c.JSON(200, r)
			return
		}
		if count == 0 {
			//r = GetError(CodeNoData)
			r = GetResponse(ResNoData)
			c.JSON(200, r)
			return
		}
		var result *gorm.DB
		if *joinStr == "" {
			result = db.Table(table).Select(*selectStr).Where(mo).Where(where, params...).Where(_map).Where(_sql).Limit(size).Offset(offset).Order(sort).Find(mos)
		} else {
			result = db.Table(table).Select(*selectStr).Joins(*joinStr).Where(mo).Where(where, params...).Where(_map).Where(_sql).Limit(size).Offset(offset).Order(sort).Find(mos)
		}
		err = result.Error
		if err != nil {
			r = GetError(CodeServer)
		} else {
			r = gin.H{"code": 20000, "data": mos, "count": count}
		}
		c.JSON(200, r)
		return
	}
	//前端分页情况
	var result *gorm.DB
	if *joinStr == "" {
		result = db.Table(table).Select(*selectStr).Where(mo).Where(where, params...).Where(_map).Where(_sql).Order(sort).Find(mos)
	} else {
		result = db.Table(table).Select(*selectStr).Joins(*joinStr).Where(mo).Where(where, params...).Where(_map).Where(_sql).Order(sort).Find(mos)
	}

	err = result.Error
	if err != nil {
		r = GetError(CodeServer)
	} else if result.RowsAffected == 0 {
		r = GetResponse(ResNoData)
	} else {
		r = gin.H{"code": 20000, "data": mos}
	}
	c.JSON(200, r)
	return
}

func getDetail(c *gin.Context, table string, mos interface{}, selectStr *string, joinStr *string) {
	var r gin.H
	var result *gorm.DB
	var err error
	var id struct {
		ID int64 `json:"id,string"`
	}
	sort, ok := detailOrderMap[table]
	if !ok {
		sort = "id desc"
	}
	if err = c.ShouldBindJSON(&id); err != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	if *selectStr == "" {
		*selectStr = "*"
	}
	if *joinStr == "" {
		result = db.Table(table).Select(*selectStr).Where("personnel_id = ?", id.ID).Order(sort).Find(mos)
	} else {
		result = db.Table(table).Select(*selectStr).Joins(*joinStr).Where("personnel_id = ?", id.ID).Order(sort).Find(mos)
	}
	if result.Error != nil {
		log.Error(result.Error)
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": mos}
	c.JSON(200, r)
}

// 如果查询结构体里有包含slice的字段，则将这些字段提取出来生成一个map返回，供gorm构建"IN"查询语句
func structToSlice(model interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	T := reflect.Indirect(reflect.ValueOf(model)).Type()
	V := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < T.NumField(); i++ {
		p := T.Field(i)
		v := V.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			var tag string
			var ok bool
			if v.IsZero() {
				continue
			}

			if tag, ok = p.Tag.Lookup("query"); !ok {
				continue
			}
			if p.Type.Kind() == reflect.Slice {
				if v.Len() == 0 {
					continue
				}
				result[tag] = v.Interface()
				if tag, ok = p.Tag.Lookup("conv"); !ok {
					switch tag {
					case "atoi":
						value := v.Interface()
						var temp []int64
						if val, yes := value.([]string); yes {
							for _, j := range val {
								_v, _ := strconv.ParseInt(j, 10, 64)
								temp = append(temp, _v)
							}
							result[tag] = temp
						}
					}
				}
			}
			v.Set(reflect.New(p.Type).Elem())
		}
	}
	return result
}

func intercept(model interface{}) string {
	var result = ""
	var value string
	T := reflect.Indirect(reflect.ValueOf(model)).Type()
	V := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < T.NumField(); i++ {
		p := T.Field(i)
		v := V.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			var ok bool
			if v.IsZero() {
				continue
			}
			if _, ok = p.Tag.Lookup("sql"); !ok {
				continue
			}
			if value, ok = v.Interface().(string); !ok {
				log.Error("error: 标记tag为sql的字段必须为字符串")
				continue
			}
			if result != "" {
				result += " and "
			}
			result += value
			v.Set(reflect.New(p.Type).Elem())
		}
	}
	return result
}
