package controllers

import (
	"GanLianInfo/models"
	"fmt"

	log "github.com/truxcoder/truxlog"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func TrainingList(c *gin.Context) {
	var mos []models.Training
	var mo models.Training
	var r gin.H
	var err error
	var count int64 //总记录数
	whereTitle := "1 = 1"

	if err = c.BindJSON(&mo); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	//因为title要用模糊查询like,所以这里拦截后端查询数据，对title进行处理
	if mo.Title != "" {
		whereTitle = "title like '%" + mo.Title + "%'"
		mo.Title = ""
	}

	size, offset := getPageData(c)

	//先查询数据总量并返回到前端
	if err = db.Table("trainings").Where(&mo).Where(whereTitle).Count(&count).Error; err != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	if count == 0 {
		r = Errors.NoData
		c.JSON(200, r)
		return
	}
	result := db.Table("trainings").Where(&mo).Where(whereTitle).Limit(size).Offset(offset).Order("start_time desc").Find(&mos)
	err = result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": mos, "count": count}
	}
	c.JSON(200, r)
	return
}

func TrainingDetail(c *gin.Context) {
	var mos []models.Training
	var r gin.H
	var result *gorm.DB
	var id struct {
		ID string `json:"id"`
	}
	if c.ShouldBindJSON(&id) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	result = db.Table("trainings").Where("id in (?)", db.Table("person_trains").Select("train_id").Where("personnel_id = ?", id.ID)).Order("start_time desc").Find(&mos)
	if result.Error != nil {
		log.Error(result.Error)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &mos}
	c.JSON(200, r)
}

// TrainPersonList 参加培训人员列表
func TrainPersonList(c *gin.Context) {
	var mos []models.PersonTrain
	var r gin.H
	var id struct {
		ID int64 `json:"id"`
	}
	if c.ShouldBindJSON(&id) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	result := db.Where("train_id = ?", id.ID).Find(&mos)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &mos}
	}
	c.JSON(200, r)
}

// TrainPersonAdd 添加人员参加培训信息
func TrainPersonAdd(c *gin.Context) {
	var r gin.H
	var mos []models.PersonTrain
	if c.ShouldBindJSON(&mos) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	db.Create(&mos)
	r = gin.H{"message": "添加成功！", "code": 20000}
	c.JSON(200, r)
}

// TrainPersonDelete 删除人员参加培训信息
func TrainPersonDelete(c *gin.Context) {
	var r gin.H
	var mos struct {
		PersonnelId string  `json:"personnelId"`
		TrainId     []int64 `json:"trainId"`
	}
	if c.ShouldBindJSON(&mos) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	result := db.Debug().Where("personnel_id = ? and train_id in (?)", mos.PersonnelId, mos.TrainId).Delete(&models.PersonTrain{})
	err := result.Error
	if err != nil {
		log.Error(err)
		r = Errors.ServerError
	} else {
		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
		r = gin.H{"message": message, "code": 20000}
	}
	c.JSON(200, r)
	return
}

//func TrainingDetail(c *gin.Context) {
//	var mos []models.Award
//	var selectStr string
//	var joinStr string
//	getDetail(c, "awards", &mos, &selectStr, &joinStr)
//}
