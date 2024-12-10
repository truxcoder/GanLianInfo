package controllers

import (
	"GanLianInfo/models"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
	"reflect"
)

type Cacher interface {
	Query()
	Obj() interface{}
	Key() string
}
type Dpt struct {
	Data *[]models.Department
}

func (d Dpt) Query() {
	db.Omit("position").Order("sort desc").Find(d.Data)
}
func (d Dpt) Obj() interface{} {
	return d.Data
}
func (d Dpt) Key() string {
	return "departmentSlice2"
}

func ReadCache(c Cacher) (err error) {
	var (
		temp string
	)
	typ := reflect.TypeOf(c.Obj())
	if typ == nil || typ.Kind() != reflect.Pointer {
		log.Errorf("[GetDataFromCache]只能传入指针\n")
		return
	}
	if rdb == nil {
		log.Errorf("redis实例为空")
		return RedisNullERR
	}
	res, _ := rdb.Exists(ctx, c.Key()).Result()
	if res <= 0 {
		log.Errorf("[GetDataFromCache]未发现key:%s\n", c.Key())
		err = RedisKeyNotFoundERR
		return
	}
	if temp, err = rdb.Get(ctx, c.Key()).Result(); err != nil {
		log.Error(err.Error())
		return
	}
	if err = jsoniter.UnmarshalFromString(temp, c.Obj()); err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func WriteCache(c Cacher) (err error) {
	var (
		temp string
	)
	c.Query()
	if rdb != nil {
		if temp, err = jsoniter.MarshalToString(c.Obj()); err != nil {
			return
		}
		rdb.Set(ctx, c.Key(), temp, expiration)
	}
	return
}
