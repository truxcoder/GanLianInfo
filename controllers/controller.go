package controllers

import (
	"GanLianInfo/dao"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	log "github.com/truxcoder/truxlog"

	"gorm.io/gorm"
)

var db *gorm.DB

//var rdb *redis.Client
//var ctx = context.Background()
var enforcer *casbin.Enforcer

var orderMap = map[string]string{
	"appraisals": "years,season",
}

var reConnect bool

func init() {
	db = dao.Connect()
	//if rdb, err = dao.RedisConnect(); err != nil {
	//	log.Error(err)
	//}
	FixDB()
	casbinInit()
	//SetPersonOrganMap()
}

type IdStruct struct {
	Id []int64 `json:"id"`
}

type ID struct {
	ID int64 `json:"id"`
}

func casbinInit() {
	a, _ := gormadapter.NewAdapterByDB(db)
	dir, _ := os.Getwd()
	//path := filepath.Dir(dir)
	//log.Infof("path:%s\n", dir)
	model := filepath.Join(dir, "model.conf")
	var err error
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
		log.Errorf("数据库Ping错误,  正在重连.\n")
		for {
			db = dao.Connect()
			sqlDB, err = db.DB()
			_err := sqlDB.Ping()
			if _err == nil {
				reConnect = false
				// 断线修复后初始化casbin模型
				casbinInit()
				//if rdb, err = dao.RedisConnect(); err != nil {
				//	log.Error(err)
				//}
				log.Successf("数据库恢复, 时间为: %s\n", time.Now())
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

func buildWhere(c *gin.Context) string {
	var where string
	organId := c.Query("organId")
	if organId == "" {
		where = " 1 = 1 "
	} else {
		where = "personnel_id in (select id from personnels where organ_id = '" + organId + "')"
	}
	return where
}

func getList(c *gin.Context, table string, mo interface{}, mos interface{}, selectStr *string, joinStr *string) {
	var where string
	var r gin.H
	var err error
	var count int64                     //总记录数
	queryMeans := c.Query("queryMeans") //请求方式，是前端分页还是后端分页
	sort, ok := orderMap[table]
	if !ok {
		sort = "id desc"
	}

	if err = c.BindJSON(mo); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	where = buildWhere(c)
	//后端分页情况
	if queryMeans == "backend" {
		size, offset := getPageData(c)
		//先查询数据总量并返回到前端
		if err = db.Table(table).Where(mo).Where(where).Count(&count).Error; err != nil {
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
		if count == 0 {
			r = Errors.NoData
			c.JSON(200, r)
			return
		}
		result := db.Table(table).Select(*selectStr).Joins(*joinStr).Where(mo).Where(where).Limit(size).Offset(offset).Order(sort).Find(mos)
		err = result.Error
		if err != nil {
			r = Errors.ServerError
		} else {
			r = gin.H{"code": 20000, "data": mos, "count": count}
		}
		c.JSON(200, r)
		return
	}
	//前端分页情况
	result := db.Table(table).Select(*selectStr).Joins(*joinStr).Where(mo).Where(where).Order(sort).Find(mos)
	err = result.Error
	if err != nil {
		r = Errors.ServerError
	} else if result.RowsAffected == 0 {
		r = Errors.NoData
	} else {
		r = gin.H{"code": 20000, "data": mos}
	}
	c.JSON(200, r)
	return
}

func getDetail(c *gin.Context, table string, mos interface{}, selectStr *string, joinStr *string) {
	var r gin.H
	var result *gorm.DB
	var id struct {
		ID string `json:"id"`
	}
	sort, ok := orderMap[table]
	if !ok {
		sort = "id desc"
	}
	if c.ShouldBindJSON(&id) != nil {
		r = Errors.ServerError
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
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": mos}
	c.JSON(200, r)
}
