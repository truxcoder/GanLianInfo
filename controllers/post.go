package controllers

import (
	"GanLianInfo/models"
	"fmt"
	log "github.com/truxcoder/truxlog"
	"gorm.io/gorm"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	LevelId []int64 `json:"levelId" gorm:"column:level_id"`
}

type Posts struct {
	models.Post
	PersonnelName string `json:"personnelName"`
	PoliceCode    string `json:"policeCode"`
	PositionName  string `json:"positionName"`
	LevelName     string `json:"levelName"`
}

type PostWithLevel struct {
	models.Post
	LevelName  string `json:"levelName"`
	LevelOrder int    `json:"levelOrder"`
}

type PostDetailStruct struct {
	models.Post
	PositionName string `json:"positionName"`
	LevelName    string `json:"levelName"`
}

func PostList(c *gin.Context) {
	var mos []Posts
	var mo struct {
		PersonnelId int64     `json:"personnelId,string"`
		EndDay      time.Time `json:"endDay"`
		PositionId  []string  `json:"positionId" gorm:"-" query:"posts.position_id" conv:"atoi"`
		LevelId     []string  `json:"levelId" gorm:"-" query:"posts.level_id" conv:"atoi"`
	}
	selectStr := "posts.*,per.name as personnel_name, per.police_code as police_code," +
		"p.name as position_name, l.name as level_name"
	joinStr := "left join personnels as per on posts.personnel_id = per.id " +
		"left join positions as p on posts.position_id = p.id " +
		"left join levels as l on posts.level_id = l.id "

	getList(c, "posts", &mo, &mos, &selectStr, &joinStr)
}

func PostDetail(c *gin.Context) {
	var mos []PostDetailStruct
	selectStr := "posts.*,p.name as position_name, l.name as level_name"
	joinStr := "left join positions as p on posts.position_id = p.id " +
		"left join levels as l on posts.level_id = l.id"
	getDetail(c, "posts", &mos, &selectStr, &joinStr)
}

// PostAdd 牵涉到要同步更新人员现有职务和职级，所以单独处理添加任职信息
func PostAdd(c *gin.Context) {
	var (
		r                      gin.H
		err                    error
		model                  models.Post
		p                      models.Personnel
		position               models.Position
		currentLevel, newLevel models.Level
		count                  int64
	)
	zeroDate := "0001-01-01 00:00:00.000000 +00:00"
	if err = c.ShouldBindJSON(&model); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	// 启用事务，确保所有操作都顺利执行
	// 任职信息的添加逻辑：先判断新信息的结束日期是否为零值。如果否，不做处理。如果是。判断是职务还是职级。
	// 如果是职务，跟人员的current_level作对比。如果新职务高于现职务(即order值更小)，则更新current_level
	// 如果是职级，判断人员是否有未结束任期的职级，如果有，则报错。如果无，则将current_rank修改为新的职级。
	err = db.Transaction(func(tx *gorm.DB) error {
		// 如果提交的是未结束的任职信息
		if model.EndDay.IsZero() {
			if _err := tx.Limit(1).Find(&position, model.PositionId).Error; _err != nil {
				return _err
			}
			if _err := tx.Limit(1).Find(&p, model.PersonnelId).Error; _err != nil {
				return _err
			}
			// 如果提交的是领导职务
			if position.IsLeader == 2 {
				// 先查找现有职务，对比现有职务的新提交的职务的高低
				if _err := tx.Limit(1).Find(&currentLevel, p.CurrentLevel).Error; _err != nil {
					return _err
				}
				if _err := tx.Limit(1).Find(&newLevel, model.LevelId).Error; _err != nil {
					return _err
				}
				if currentLevel.Orders == 0 || currentLevel.Orders > newLevel.Orders {
					if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_level", model.LevelId).Error; _err != nil {
						// 返回任何错误都会回滚事务
						return _err
					}
				}
				// 如果提交的是非领导职务
			} else if position.IsLeader == 1 {
				if _err := tx.Model(&models.Post{}).Where("personnel_id = ? and end_day = ? and position_id in (select id from positions where is_leader = 1)", model.PersonnelId, zeroDate).Count(&count).Error; _err != nil {
					return _err
				}
				if count > 0 {
					return fmt.Errorf("系统检查到该人员有未结束任期的职级信息！")
				} else if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_rank", model.PositionId).Error; _err != nil {
					return _err
				}
			}
		}
		if _err := tx.Create(&model).Error; _err != nil {
			// 返回任何错误都会回滚事务
			return _err
		}
		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		//r = GetError(CodeAdd)
		r = gin.H{"message": err.Error(), "code": 50500}
		log.Error(err)
		c.JSON(200, r)
		return
	}
	r = gin.H{"message": "添加成功！", "code": 20000}
	c.JSON(200, r)
}

// PostUpdate 牵涉到要同步更新人员现有职务和职级，所以单独处理修改任职信息
func PostUpdate(c *gin.Context) {
	var (
		r                      gin.H
		err                    error
		model, oldPost         models.Post
		p                      models.Personnel
		position               models.Position
		currentLevel, newLevel models.Level
		existPost              []PostWithLevel
	)
	zeroDate := "0001-01-01 00:00:00.000000 +00:00"
	if err = c.ShouldBindJSON(&model); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Limit(1).Find(&position, model.PositionId)
	db.First(&p, model.PersonnelId)
	db.Limit(1).Find(&oldPost, model.ID)
	// 启用事务，确保所有操作都顺利执行
	// 任职信息的修改逻辑： 先判断新信息的结束日期是否为零值。
	// 结束日期如果不是零值，则判断此条信息的结束日期是否是从零值修改为非零值。如果不是，不做处理。如果是，则判断此条信息是职务还是职级。
	// 如果是职务，判断该人员是否还有未结束任期的职务，有则把current_level改为未结束任期的最高职务，无则把current_level置为null
	// 如果是职级，则将current_rank置为null
	// 结束日期如果是零值，则判断此条信息是职务还是职级。如果是职务，先查找现有职务，对比现有职务的新提交的职务的高低，如果新职务更高，则将
	// current_level修改为新的levelID。否则不做处理。
	// 如果是职级，判断人员是否有未结束任期的职级，如果有，则报错。如果无，则将current_rank修改为新的职级。
	err = db.Transaction(func(tx *gorm.DB) error {
		// 如果提交的是结束的任职信息
		if !model.EndDay.IsZero() {
			//判断任职结束日期是否从未结束修改为结束
			if oldPost.EndDay.IsZero() {
				// 如果提交的是领导职务
				if position.IsLeader == 2 {
					//判断该人员是否还有未结束任期的职务，有则把current_level改为未结束任期的最高职务，无则把current_level置为null
					if _err := tx.Model(&models.Post{}).Select("posts.*, levels.name level_name, levels.orders level_order").Joins("left join levels on levels.id = posts.level_id").Where("personnel_id = ? and end_day = ? and position_id in (select id from positions where is_leader = 2)", model.PersonnelId, zeroDate).Find(&existPost).Error; _err != nil {
						return _err
					}
					if len(existPost) == 1 {
						if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_level", nil).Error; _err != nil {
							return _err
						}
					} else {
						_order := 100
						var _levelId int64
						// 循环判断出级别最高任职信息，取出其id
						for _, v := range existPost {
							if v.LevelOrder < _order && v.ID != model.ID {
								_levelId = v.LevelId
							}
						}
						if _levelId != 0 {
							if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_level", _levelId).Error; _err != nil {
								return _err
							}
						}
					}
				} else {
					// 把人员当前职级置为null
					if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_rank", nil).Error; _err != nil {
						return _err
					}
				}
			}
		} else {

			// 如果提交的是领导职务
			if position.IsLeader == 2 {
				// 先查找现有职务，对比现有职务的新提交的职务的高低
				if _err := tx.Limit(1).Find(&currentLevel, p.CurrentLevel).Error; _err != nil {
					return _err
				}
				if _err := tx.Limit(1).Find(&newLevel, model.LevelId).Error; _err != nil {
					return _err
				}
				if currentLevel.Orders == 0 || currentLevel.Orders > newLevel.Orders {
					if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_level", model.LevelId).Error; _err != nil {
						return _err
					}
				}
				// 如果提交的是非领导职务
			} else if position.IsLeader == 1 {
				if _err := tx.Model(&models.Post{}).Select("posts.*, levels.name level_name, levels.orders level_order").Joins("left join levels on levels.id = posts.level_id").Where("personnel_id = ? and end_day = ? and position_id in (select id from positions where is_leader = 1)", model.PersonnelId, zeroDate).Find(&existPost).Error; _err != nil {
					return _err
				}
				if len(existPost) > 1 || (len(existPost) == 1 && existPost[0].ID != model.ID) {
					return fmt.Errorf("系统检查到该人员有未结束任期的职级信息！")
				} else if _err := tx.Model(&models.Personnel{}).Where("id = ?", model.PersonnelId).Update("current_rank", model.PositionId).Error; _err != nil {
					return _err
				}
			}
		}
		if _err := tx.Model(&model).Updates(&model).Error; _err != nil {
			// 返回任何错误都会回滚事务
			return _err
		}

		// 返回 nil 提交事务
		return nil

	})

	if err != nil {
		r = gin.H{"message": err.Error(), "code": 50500}
		log.Error(err)
		c.JSON(200, r)
		return
	}
	updateZeroFields(&model)
	r = gin.H{"message": "更新成功！", "code": 20000}
	c.JSON(200, r)
}
