package controllers

import (
	"GanLianInfo/models"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"

	log "github.com/truxcoder/truxlog"
)

const (
	expiration = time.Hour * 24 * 365
)

func getIdCodeMap() map[string]string {
	_map := make(map[string]string)
	var err error
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "idCodeList").Result()
		if exist == 0 {
			setIdCodeMap()
		}
		if _map, err = rdb.HGetAll(ctx, "idCodeList").Result(); err != nil {
			log.Error(err)
		}
		return _map
	}
	log.Info("redis instance is nil, now get data from database...")
	return getIdCodeMapFromDB()
}

func setIdCodeMap() {
	_map := getIdCodeMapFromDB()
	if rdb != nil {
		rdb.Del(ctx, "idCodeList")
		rdb.HSet(ctx, "idCodeList", _map)
	}
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
			setIdCodeMap()
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
			setDepartmentMap()
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

func setDepartmentMap() {
	_map := getDepartmentMapFromDB()
	if rdb != nil {
		rdb.Del(ctx, "departmentMap")
		rdb.HSet(ctx, "departmentMap", _map)
	}
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

func getDepartmentSliceFromDB() []models.Department {
	var (
		d []models.Department
	)
	db.Omit("position").Order("sort desc").Find(&d)
	return d
}

func setDepartmentSlice() {
	var result string
	var err error
	d := getDepartmentSliceFromDB()
	if len(d) > 0 {
		if result, err = jsoniter.MarshalToString(d); err != nil {
			log.Error(err)
		}
	}
	if rdb != nil && result != "" {
		rdb.Del(ctx, "departmentSlice")
		rdb.Set(ctx, "departmentSlice", result, expiration)
	}
}

func getDepartmentSlice() ([]models.Department, error) {
	var d []models.Department
	var temp string
	var err error
	var result []models.Department
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "departmentSlice").Result()
		if exist == 0 {
			setDepartmentSlice()
		}
		if temp, err = rdb.Get(ctx, "departmentSlice").Result(); err != nil {
			log.Error(err)
			return nil, err
		}
		if err = jsoniter.UnmarshalFromString(temp, &result); err != nil {
			log.Error(err)
			return nil, err
		}
		return result, nil
	}
	log.Info("redis instance is nil, now get data from database...")
	d = getDepartmentSliceFromDB()
	return d, nil
}

func GetDepSlice() []models.Department {
	result, _ := getDepartmentSlice()
	return result
}
