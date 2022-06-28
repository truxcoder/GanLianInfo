package controllers

import (
	"GanLianInfo/models"

	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

type Positions struct {
	models.Position
	LevelName  string `json:"levelName"`
	LevelOrder string `json:"levelOrder"`
}

func PositionList(c *gin.Context) {
	var p []Positions
	var r gin.H
	selectStr := "positions.*,levels.name as level_name, levels.orders as level_order "
	joinStr := "left join levels on positions.level_id = levels.id "
	// 用Order排序要点：Gorm在用Joins时会把关联的表字段重命名为表名__字段名的形式
	result := db.Select(selectStr).Joins(joinStr).Order("is_leader desc,level_order").Find(&p)
	//result := db.Joins("Level", db.Order("orders desc")).Find(&p)
	err := result.Error
	if err != nil {
		r = GetError(CodeServer)
	} else {
		r = gin.H{"code": 20000, "data": &p}
	}
	c.JSON(200, r)
}

func PositionCheck(c *gin.Context) {
	var (
		pos   models.Position
		err   error
		r     gin.H
		count int64
	)
	if err = c.BindJSON(&pos); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	db.Table("positions").Where("name = ? and level_id = ?", pos.Name, pos.LevelId).Count(&count)
	if count > 0 {
		log.Errorf("count:%d\n", count)
		r = GetError(CodeUnique)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "count": &count}
	c.JSON(200, r)
}
