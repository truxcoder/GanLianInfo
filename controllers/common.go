package controllers

import (
	"GanLianInfo/models"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truxcoder/gorm-dm8/datatype"

	jsoniter "github.com/json-iterator/go"

	log "github.com/truxcoder/truxlog"
)

const (
	expiration = time.Hour * 24 * 365
)

func GetDataFromCache(obj interface{}, key string) (err error) {
	var (
		temp string
	)
	typ := reflect.TypeOf(obj)
	if typ == nil || typ.Kind() != reflect.Pointer {
		log.Errorf("[GetDataFromCache]只能传入指针\n")
		return
	}
	if rdb != nil {
		res, _ := rdb.Exists(ctx, key).Result()
		if res > 0 {
			if temp, err = rdb.Get(ctx, key).Result(); err != nil {
				log.Error(err.Error())
				return
			}
			if err = jsoniter.UnmarshalFromString(temp, obj); err != nil {
				log.Error(err.Error())
				return
			}
		} else {
			log.Errorf("[GetDataFromCache]未发现key:%s\n", key)
			err = RedisKeyNotFoundERR
			return
		}
	}
	return
}
func WriteDataToCache(data interface{}, key string) (err error) {
	var (
		temp string
	)
	if rdb != nil {
		if temp, err = jsoniter.MarshalToString(data); err != nil {
			return
		}
		rdb.Set(ctx, key, temp, expiration)
	}
	return
}

func getIdCodeMapFromDB() map[string]string {
	var p []struct {
		ID     int64
		IdCode string
	}
	_map := make(map[string]string)
	db.Table("personnels").Select("id, id_code").Find(&p)
	if len(p) > 0 {
		for _, v := range p {
			_map[v.IdCode] = strconv.FormatInt(v.ID, 10)
		}
	}
	return _map
}

func getIdFromIdCode(idCode string) int64 {
	var (
		id  int64
		err error
		p   struct {
			ID int64
		}
	)
	if idCode == "" {
		return 0
	}
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "idCodeList").Result()
		if exist == 0 {
			setIdCodeMapToCache()
		}
		if id, err = rdb.HGet(ctx, "idCodeList", idCode).Int64(); err != nil {
			return 0
		}
		return id

	}
	log.Info("redis instance is nil, now get data from database...")
	db.Table("personnels").Select("id").Where("id_code = ?", idCode).Limit(1).Find(&p)
	return p.ID
}

func getOrganIdFromDepartmentId(departmentId string) string {
	var (
		id  string
		err error
		d   struct {
			ID string
		}
	)
	if departmentId == "" {
		return ""
	}
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "departmentMap").Result()
		if exist == 0 {
			setDepartmentMapToCache()
		}
		if id, err = rdb.HGet(ctx, "departmentMap", departmentId).Result(); err != nil {
			return ""
		}
		return id
	}
	log.Info("redis instance is nil, now get data from database...")
	db.Table("departments").Select("id").Where("bus_org_code = (select bus_org_code from departments where id = ?) and dept_type = 1", departmentId).Limit(1).Find(&d)
	return d.ID
}

func getDepartmentMapFromDB() map[string]string {
	var d []struct {
		ID      string
		OrganId string
	}
	_map := make(map[string]string)
	joinStr := "left join departments as d on (d.bus_org_code = departments.bus_org_code and d.dept_type = 1)"
	selectStr := "departments.id, d.id as organ_id"
	db.Table("departments").Select(selectStr).Joins(joinStr).Find(&d)
	if len(d) > 0 {
		for _, v := range d {
			_map[v.ID] = v.OrganId
		}
	}
	return _map
}

func getDepartmentSliceFromDB() (d []models.Department, err error) {
	err = db.Omit("position").Order("sort desc").Find(&d).Error
	return
}

func getDepartmentSlice() (d []models.Department, err error) {
	if d, err = getDepartmentSliceFromCache(); err != nil {
		log.Info("redis instance is nil, now get data from database...")
		if d, err = getDepartmentSliceFromDB(); err != nil {
			log.Errorf("[getDepartmentSliceFromDB]%v\n", err.Error())
			return
		}
		return d, nil
	}
	return
	//if rdb != nil {
	//	exist, _ := rdb.Exists(ctx, "departmentSlice").Result()
	//	if exist == 0 {
	//		setDepartmentSliceToCache()
	//	}
	//	if temp, err = rdb.Get(ctx, "departmentSlice").Result(); err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	if err = jsoniter.UnmarshalFromString(temp, &result); err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	return result, nil
	//}
	//log.Info("redis instance is nil, now get data from database...")
	//d = getDepartmentSliceFromDB()
	//return d, nil
}

func WriteLog(c *gin.Context, category LogCategory, content string) {
	var (
		l models.Log
	)
	l.IP = c.ClientIP()
	l.AccountId = c.GetString("userId")
	l.Content = datatype.Clob(content)
	l.Category = int8(category)
	db.Create(&l)
}

// 计算退休年龄
func getRetireAge(p *posStruct) int {
	// 女性特殊退休年龄，一，现曾任处级领导职务。二，2019年3月1日前担任副处级非领导职务，例如副调研员。
	line := time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)
	isLeader := !p.FcStartDay.IsZero() || !p.ZcStartDay.IsZero()
	//isNonLeader := !getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 0).IsZero() && getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 0).Before(line)
	isNonLeader := getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 0).Before(line)
	if p.Gender == "女" && !(isLeader || isNonLeader) {
		return 55
	}
	return 60
}
