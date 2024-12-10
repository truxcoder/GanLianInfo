package controllers

import (
	"GanLianInfo/models"
	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
	"slices"
	"time"
)

// AnalysisHighLevelData  获取数据分析高级警长数据
func AnalysisHighLevelData(c *gin.Context) {
	var (
		err error
		r   gin.H
		hl  []HighLevelPolice
	)
	if hl, err = getAnalysisHighLevelData(); err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": hl}
	c.JSON(200, r)
}

func getAnalysisHighLevelData() (hl []HighLevelPolice, err error) {
	var (
		pos []posStruct
	)
	leaderPostMap := make(map[int64][]PostLevelPosition)
	nonLeaderPostMap := make(map[int64][]PostLevelPosition)
	if pos, err = getPostsFromCache(); err != nil {
		if pos, err = getPostsFromDB(); err != nil {
			return
		}
	}
	if leaderPostMap, err = getPostMapFromCache(true); err != nil {
		log.Errorf(err.Error())
		return
	}
	if nonLeaderPostMap, err = getPostMapFromCache(false); err != nil {
		log.Errorf(err.Error())
		return
	}

	for _, p := range pos {
		var (
			h         HighLevelPolice
			g4BaseDay time.Time //四级警长任取基本日期，即一级警长满两年日期
		)
		h.ID = p.ID
		h.Name = p.Name
		h.OrganID = p.OrganID
		h.Birthday = p.Birthday
		h.Gender = p.Gender

		if v, exist := leaderPostMap[p.ID]; exist {
			p.LeaderPosts = v
		}
		if v, exist := nonLeaderPostMap[p.ID]; exist {
			p.NonLeaderPosts = v
		}
		//log.Successf("%s-%s-%d\n", p.Name, p.Gender, getRetireAge(&p))
		p.RetireDay = p.Birthday.AddDate(getRetireAge(&p), 0, 0)
		h.RetireDay = p.RetireDay
		// 如果正科开始日期不为空而副科开始日期为空，则将副科开始日期设为正科开始日期。解决部分人员副科任职经历未录入或直接任正科的情况。
		if !p.ZkStartDay.IsZero() && p.FkStartDay.IsZero() {
			p.FkStartDay = p.ZkStartDay
		}
		//当前为一高情况
		if p.RankName == "一级高级警长" {
			h.G1StartDay = p.RankStartDay
			// 将此人员信息加入切片
			hl = append(hl, h)
			continue
		}
		//当前为二高情况
		if p.RankName == "二级高级警长" {
			h.G2StartDay = p.RankStartDay
			g1BaseDay := h.G2StartDay.AddDate(3, 0, 0)
			h.G1StartDay = getG1StartDay(g1BaseDay, &p)
			// 将此人员信息加入切片
			hl = append(hl, h)
			continue
		}
		//当前为三高情况
		if p.RankName == "三级高级警长" {
			h.G3StartDay = p.RankStartDay
			g2BaseDay := h.G3StartDay.AddDate(2, 0, 0)
			h.G2StartDay = getG2StartDay(g2BaseDay, &p)
			g1BaseDay := h.G2StartDay.AddDate(3, 0, 0)
			h.G1StartDay = getG1StartDay(g1BaseDay, &p)
			// 将此人员信息加入切片
			hl = append(hl, h)
			continue
		}
		//当前为四高情况
		if p.RankName == "四级高级警长" {
			h.G4StartDay = p.RankStartDay
			g3BaseDay := h.G4StartDay.AddDate(2, 0, 0)
			h.G3StartDay = getG3StartDay(g3BaseDay, &p)
			g2BaseDay := h.G3StartDay.AddDate(2, 0, 0)
			h.G2StartDay = getG2StartDay(g2BaseDay, &p)
			g1BaseDay := h.G2StartDay.AddDate(3, 0, 0)
			h.G1StartDay = getG1StartDay(g1BaseDay, &p)
			// 将此人员信息加入切片
			hl = append(hl, h)
			continue
		}
		// 计算出晋升四高的基本日期
		switch p.RankName {
		case "四级警长":
			g4BaseDay = p.RankStartDay.AddDate(8, 0, 0)
		case "三级警长":
			g4BaseDay = p.RankStartDay.AddDate(6, 0, 0)
		case "二级警长":
			g4BaseDay = p.RankStartDay.AddDate(4, 0, 0)
		case "一级警长":
			// 这里需要特殊处理。要计算正科职级的真正起算时间。包括主任科员，一级主任科员和正科级领导职务开始日期一并计算，取最早日期。
			//rankRealStartDay := getRealStartDay(&p)
			//g4BaseDay = rankRealStartDay.AddDate(2, 0, 0)
			g4BaseDay = p.RankStartDay.AddDate(2, 0, 0)
		}
		switch p.RankName {
		case "四级警长", "三级警长", "二级警长", "一级警长":
			h.G4StartDay = getG4StartDay(g4BaseDay, &p)
			g3BaseDay := h.G4StartDay.AddDate(2, 0, 0)
			h.G3StartDay = getG3StartDay(g3BaseDay, &p)
			g2BaseDay := h.G3StartDay.AddDate(2, 0, 0)
			h.G2StartDay = getG2StartDay(g2BaseDay, &p)
			g1BaseDay := h.G2StartDay.AddDate(3, 0, 0)
			h.G1StartDay = getG1StartDay(g1BaseDay, &p)
		}
		// 将此人员信息加入切片
		hl = append(hl, h)
	}
	return
}

// 从数据库获取高警任职信息
func getPostsFromDB() (pos []posStruct, err error) {
	var (
		d models.Department
	)
	selectStr := "current_pos.id,current_pos.name,current_pos.current_rank,current_pos.rank_name,current_pos.current_level,current_pos.level_name,current_pos.rank_start_day," +
		"current_pos.level_start_day,used_pos.fk_start_day,used_pos.zk_start_day,used_pos.fc_start_day,used_pos.zc_start_day," +
		"personnels.birthday, personnels.organ_id, personnels.gender"
	joinStr := "left join used_pos on current_pos.id = used_pos.id left join personnels on personnels.id = current_pos.id"
	db.Model(&models.Department{}).Where("name = ?", "四川省成都强制隔离戒毒所").Limit(1).Find(&d)
	if err = db.Table("current_pos").Select(selectStr).Joins(joinStr).Where("current_pos.id in (select id from personnels where personnels.organ_id = ? and personnels.status = 1)", d.ID).Find(&pos).Error; err != nil {
		log.Errorf(err.Error())
	}
	return
}

// 从数据库中获取任职信息，并以k:v对的形式生成map，分领导职务和非领导职务
func getPostMapFromDB() (leaderPostMap map[int64][]PostLevelPosition, nonLeaderPostMap map[int64][]PostLevelPosition) {
	var (
		leaderPosts, nonLeaderPosts []PostLevelPosition
	)
	leaderPostMap = make(map[int64][]PostLevelPosition)
	nonLeaderPostMap = make(map[int64][]PostLevelPosition)
	selectStr := "posts.id,posts.personnel_id, posts.start_day, posts.end_day,positions.is_leader, levels.name level_name, positions.name position_name, levels.orders"
	joinStr := "left join positions on posts.position_id = positions.id left join levels on posts.level_id = levels.id"
	leaderWhere := "is_leader = 2 and exists (select 1 from personnels where personnels.id = posts.personnel_id and personnels.organ_id = ?)"
	nonLeaderWhere := "is_leader = 1 and exists (select 1 from personnels where personnels.id = posts.personnel_id and personnels.organ_id = ?)"
	db.Table("posts").Select(selectStr).Joins(joinStr).Where(leaderWhere, "3d7e73e3a3034ca1a1da707aa3d54a96").Find(&leaderPosts)
	db.Table("posts").Select(selectStr).Joins(joinStr).Where(nonLeaderWhere, "3d7e73e3a3034ca1a1da707aa3d54a96").Find(&nonLeaderPosts)
	for _, post := range leaderPosts {
		if v, exists := leaderPostMap[post.PersonnelId]; exists {
			leaderPostMap[post.PersonnelId] = append(v, post)
		} else {
			leaderPostMap[post.PersonnelId] = []PostLevelPosition{post}
		}
	}
	for k, v := range leaderPostMap {
		slices.SortFunc(v, func(a, b PostLevelPosition) int {
			return a.StartDay.Compare(b.StartDay)
		})
		leaderPostMap[k] = v
	}
	for _, post := range nonLeaderPosts {
		if v, exists := nonLeaderPostMap[post.PersonnelId]; exists {
			nonLeaderPostMap[post.PersonnelId] = append(v, post)
		} else {
			nonLeaderPostMap[post.PersonnelId] = []PostLevelPosition{post}
		}
	}
	for k, v := range nonLeaderPostMap {
		slices.SortFunc(v, func(a, b PostLevelPosition) int {
			return a.StartDay.Compare(b.StartDay)
		})
		nonLeaderPostMap[k] = v
	}
	return
}

func SyncPostData(c *gin.Context) {
	var (
		errs []error
		err  error
		r    gin.H
	)
	if err = setPostsToCache(); err != nil {
		errs = append(errs, err)
	}
	if err = setPostMapToCache(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		for _, err = range errs {
			log.Error(err.Error())
		}
		r = GetError(CodeServer)
	} else {
		r = gin.H{"code": 20000, "message": "更新成功！"}
	}
	c.JSON(200, r)
}
