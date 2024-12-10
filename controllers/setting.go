package controllers

import (
	"GanLianInfo/models"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

// SettingList 获取设置信息
func SettingList(c *gin.Context) {
	var mos []models.Setting
	var err error
	var r gin.H
	mos, err = getSetting()
	if err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &mos}
	c.JSON(200, r)
}

func SettingUpdate(c *gin.Context) {
	var (
		r     gin.H
		err   error
		model models.Setting
	)
	if err = c.ShouldBindJSON(&model); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Model(&model).Updates(&model)
	if err = setSetting(); err != nil {
		r = GetError(CodeDataWrite)
		c.JSON(200, r)
		return
	}
	r = gin.H{"message": "更新成功！", "code": 20000}
	c.JSON(200, r)
}

func getSetting() ([]models.Setting, error) {
	var s []models.Setting
	var temp string
	var err error
	var result []models.Setting
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "setting").Result()
		if exist == 0 {
			if err = setSetting(); err != nil {
				return nil, err
			}
		}
		if temp, err = rdb.Get(ctx, "setting").Result(); err != nil {
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
	s = getSettingFromDB()
	return s, nil
}

func setSetting() error {
	var result string
	var err error
	s := getSettingFromDB()
	if len(s) > 0 {
		if result, err = jsoniter.MarshalToString(s); err != nil {
			return err
		}
	}
	if rdb != nil && result != "" {
		rdb.Del(ctx, "setting")
		rdb.Set(ctx, "setting", result, expiration)
	}
	return nil
}

func getSettingFromDB() []models.Setting {
	var (
		s []models.Setting
	)
	db.Omit("created_at", "updated_at", "version").Find(&s)
	return s
}
